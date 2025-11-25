package file

import (
	"errors"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	lfile "github.com/abtransitionit/golinux/mock/file"
)

// Description: sudo copy a file to a set of nodes from hostname
func CopyFileWithSudo(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	// 11 - node:List
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - fileProperty
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("host: %s > fileProperty not provided in paramList", hostName)
	}
	fileProperty, err := lfile.GetVarStruct[lfile.FileProperty](fmt.Sprint(paramList[1][0]))
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage goroutines concurrency
	nbNode := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChNode := make(chan error, nbNode) // define a channel to collect errors from each goroutine

	// 3 - loop over item (node)
	for _, node := range nodeList {
		wgHost.Add(1)             // Increment the WaitGroup:counter for each node
		go func(oneNode string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer func() {
				logger.Infof("↩ (%s) %s/%s > finished", phaseName, hostName, oneNode)
				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
			}()
			logger.Infof("↪ (%s) %s/%s > ongoing", phaseName, hostName, oneNode)
			_, grErr := lfile.CopyFileWithSudo(hostName, node, fileProperty, logger) // the code to be executed by the goroutine
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
