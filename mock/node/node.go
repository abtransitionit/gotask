package node

import (
	"errors"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: check if a set of node is SSH configured on a host.
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func CheckSshConf(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// define var
	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of node that are not SSH configured

	// loop over each node
	for _, node := range nodeList {

		// play CLI - check if SSH is configured for the couple host/node
		ok, err := lnode.IsSshConfigured(hostName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > Node %s: > system error > checking SSH config: %v", hostName, node, err)
			continue
		}
		// failedNodes = append(failedNodes, node)

		// collect results
		results[node] = ok
		if !ok {
			failedNodes = append(failedNodes, node) // logical error: SSH simply not configured
			logger.Debugf("taget: %s > Node %s: > is not SSH configured", hostName, node)
		}

	}

	// If any node failed, return a single error message
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("host: %s > Node(s) that are not SSH configured: %v", hostName, failedNodes)
	}

	// handle success
	return true, nil
}

// Description: check if a set of node is SSH reachable from a host
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func CheckSshAccess(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// define var
	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of nodes that are not SSH reachable

	// loop over item (node)
	for _, node := range nodeList {

		// play CLI for each item - check if node is SSH reachable for the couple host/node
		ok, err := lnode.IsSshReachable(hostName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > node %s: > system error > checking SSH access: %v", hostName, node, err)
			continue
		}

		// manage and collect logic errors
		results[node] = ok
		if !ok {
			failedNodes = append(failedNodes, node)
			// log
			logger.Infof("host: %s > node %s: > is not SSH reachable", hostName, node)
		}
	}

	// errors summary
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("host: %s > Node(s) that are not SSH reachable: %v", hostName, failedNodes)
	}

	// handle success
	return true, nil
}

// Description: Wait a set of node to be SSH reachable (within a delayMax) from a host.
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func WaitIsSshOnline(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
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
	nbNode := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each node : wait for all (concurent) goroutines (one per node) to complete
	errChNode := make(chan error, nbNode) // define a channel to collect errors from goroutines

	// 3 - loop over item (node)
	for _, node := range nodeList {
		wgHost.Add(1) // Increment the WaitGroup:counter for each node
		logger.Infof("↪ (%s) %s/%s > running", phaseName, hostName, node)
		go func(oneNode string) { // create as goroutine (that will run concurrently) as node  AND pass it as an argument
			defer wgHost.Done()                                                // Decrement the WaitGroup counter - when the goroutine complete
			_, grErr := lnode.IsSshOnline(hostName, oneNode, delayMax, logger) // the goroutin execute in fact this code
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneNode, grErr) // send goroutines error if any into the chanel
				// send goroutines error if any into the chanel
				errChNode <- fmt.Errorf("%w", grErr)
			}

		}(node) // pass the node to the goroutine
	} // node loop

	wgHost.Wait()    // Wait for all goroutines to complete - done with the help of the WaitGroup:counter
	close(errChNode) // close the channel - signal that no more error will be sent

	// 4 - collect errors
	var errList []error
	for e := range errChNode {
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

// Description: reboot a set of node if needed
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
// func NeedReboot(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
// 	// 1 - extract parameters
// 	// 11 - node:List
// 	nodeList := []string{}
// 	for _, v := range paramList[0] {
// 		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
// 	}

// 	// 2 - manage goroutines concurrency
// 	nbNode := len(nodeList)
// 	var wgHost sync.WaitGroup             // define a WaitGroup instance for each node : wait for all (concurent) goroutines (one per node) to complete
// 	errChNode := make(chan error, nbNode) // define a channel to collect errors from goroutines

// 	// 3 - loop over item (node)
// 	for _, node := range nodeList {
// 		wgHost.Add(1) // Increment the WaitGroup:counter for each node
// 		logger.Infof("↪ (goroutine) %s/%s > running", hostName, node)
// 		go func(oneNode string) { // create as goroutine (that will run concurrently) as node  AND pass it as an argument
// 			defer wgHost.Done()                                                // Decrement the WaitGroup counter - when the goroutine complete
// 			_, grErr := lnode.IsSshOnline(hostName, oneNode, delayMax, logger) // the goroutin execute in fact this code
// 			if grErr != nil {
// 				logger.Errorf("(goroutine) %s/%s > %v", hostName, oneNode, grErr) // send goroutines error if any into the chanel
// 				// send goroutines error if any into the chanel
// 				errChNode <- fmt.Errorf("%w", grErr)
// 			}

// 		}(node) // pass the node to the goroutine
// 	} // node loop

// 	wgHost.Wait()    // Wait for all goroutines to complete - done with the help of the WaitGroup:counter
// 	close(errChNode) // close the channel - signal that no more error will be sent

// 	// 4 - collect errors
// 	var errList []error
// 	for e := range errChNode {
// 		errList = append(errList, e)
// 	}

// 	// 5 - handle errors
// 	nbGroutineFailed := len(errList)
// 	errCombined := errors.Join(errList...)
// 	if nbGroutineFailed > 0 {
// 		logger.Errorf("❌ host: %s > nb node that failed: %d", hostName, nbGroutineFailed)
// 		return false, errCombined
// 	}

// 	// 6 - handle success
// 	return true, nil
// }

// ------- pattern for target/host that has node in param -------
// for _, item := range itemList {
//     result, err := check(item)
//     if err != nil {
//         logger.Warnf("item %s: check error: %v", item, err)
//         result = false
//     }
//     if !result {
//         failedItems = append(failedItems, item)
//     }
// }
// if len(failedItems) > 0 {
//     return false, fmt.Errorf("some items failed: %v", failedItems)
// }
// return true, nil

// ------- pattern for target/host that has node in param -------
// for each node in nodeList:
//     run code concurrently (goroutine) : check shh access for example
// wait for all goroutines to comple --> waitgroup
// collect errors / failures --> chanel
// return aggregated result
