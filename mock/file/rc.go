package file

import (
	"errors"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/file"
)

// Description: for each node, add to the user's custom RC file the envvar:PATH built from a root folder path
func RcAddPath(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	// 11 - check them
	if len(paramList) < 3 || len(paramList[0]) == 0 || len(paramList[1]) == 0 || len(paramList[2]) == 0 {
		return false, fmt.Errorf("%s > Node(s) or folder path or custom rc file not provided in paramList", hostName)
	}
	// 12 - node:List
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}
	// 13 - folder from which to build a path
	folderRootPath := fmt.Sprint(paramList[1][0])
	// 14 - custom rc file
	rcCustom := fmt.Sprint(paramList[2][0])

	// 2 - manage goroutines concurrency
	nbItem := len(nodeList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChNode := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item (node)
	for _, node := range nodeList {
		wgHost.Add(1)             // Increment the WaitGroup:counter for each item
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer func() {
				logger.Infof("↩ (%s) %s:%s > finished", phaseName, hostName, oneItem)
				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
			}()
			logger.Infof("↪ (%s) %s:%s > ongoing", phaseName, hostName, oneItem)
			grErr := file.RcAddPath(hostName, node, folderRootPath, rcCustom, logger) // the code to be executed by the goroutine
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
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

func AddPathToRcFile(hostName, nodeName, folderRootPath, customRcFileName string, logger logx.Logger) error {
	// operate
	err := file.RcAddPath(hostName, nodeName, folderRootPath, customRcFileName, logger)
	if err != nil {
		return fmt.Errorf("%s:%s > adding path to rc file > %w", hostName, nodeName, err)
	}

	// handle success
	return nil

}
