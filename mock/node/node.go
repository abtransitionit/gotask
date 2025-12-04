package node

import (
	"errors"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: checks if a set of node is SSH configured on a host.
func CheckSshConf(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// 2 - manage goroutines concurrency
	nbItem := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item
	for _, item := range nodeList {
		wgHost.Add(1)             // Increment the WaitGroup:counter for this item
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer func() {
				logger.Infof("↩ (%s) %s/%s > finished", phaseName, hostName, oneItem)
				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
			}()
			logger.Infof("↪ (%s) %s/%s > ongoing", phaseName, hostName, oneItem)
			_, grErr := lnode.IsSshConfigured(hostName, item, logger) // the code to be executed by the goroutine
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
				// send goroutines error if any into the chanel
				errChItem <- fmt.Errorf("%w", grErr)
			}

		}(item) // pass the node to the goroutine
	} // node loop

	wgHost.Wait()    // Wait for all goroutines to complete - done with the help of the WaitGroup:counter
	close(errChItem) // close the channel - signal that no more error will be sent

	// 4 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 5 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ host: %s > nb node that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// 6 - handle success
	return true, nil
}

// Description: checks if a set of node is SSH reachable from a host
func CheckSshAccess(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// 2 - manage goroutines concurrency
	nbItem := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item (node)
	for _, node := range nodeList {
		wgHost.Add(1) // Increment the WaitGroup:counter for this item
		logger.Infof("↪ (%s) %s/%s > running", phaseName, hostName, node)
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer wgHost.Done()                                      // Decrement the WaitGroup counter - when the goroutine complete
			_, grErr := lnode.IsSshReachable(hostName, node, logger) // the code to be executed by the goroutine
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
				// send goroutines error if any into the chanel
				errChItem <- fmt.Errorf("%w", grErr)
			}

		}(node) // pass the node to the goroutine
	} // node loop

	wgHost.Wait()    // Wait for all goroutines to complete - done with the help of the WaitGroup:counter
	close(errChItem) // close the channel - signal that no more error will be sent

	// 4 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 5 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ host: %s > nb node that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// 6 - handle success
	return true, nil
}

// Description: Waits a set of node to be SSH reachable (within a delayMax) from a host.
func WaitIsSshOnline(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	// 11 - node:List
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - repo:folder
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("host: %s > delay not provided in paramList", hostName)
	}
	delayMax := fmt.Sprint(paramList[1][0])

	// 2 - manage goroutines concurrency
	nbItem := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item (node)
	for _, node := range nodeList {
		wgHost.Add(1) // Increment the WaitGroup:counter for each item
		logger.Infof("↪ (%s) %s/%s > running", phaseName, hostName, node)
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer wgHost.Done()                                                // Decrement the WaitGroup counter - when the goroutine complete
			_, grErr := lnode.IsSshOnline(hostName, oneItem, delayMax, logger) // the code to be executed by the goroutine
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
				// send goroutines error if any into the chanel
				errChItem <- fmt.Errorf("%w", grErr)
			}

		}(node) // pass the node to the goroutine
	} // node loop

	wgHost.Wait()    // Wait for all goroutines to complete - done with the help of the WaitGroup:counter
	close(errChItem) // close the channel - signal that no more error will be sent

	// 4 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 5 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ host: %s > nb node that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// 6 - handle success
	return true, nil
}

// Description: reboots a set of node if needed
func RebootIfNeeded(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	_, err := lnode.RebootIfNeeded(hostName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("host: %s > system error > getting reboot status: %v", hostName, err)
	}

	// handle success
	// logger.Debugf("host: %s > need reboot: %s", hostName, out)
	return true, nil
}
