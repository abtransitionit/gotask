package git

import (
	"errors"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	lgit "github.com/abtransitionit/golinux/mock/git"
)

// Description: git merge branch dev to main and git push to github (for a set of git repositories)
//
// Notes:
//
// - the host contains the git repository
func MergeDevToMain(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	// 11 - repo:list
	repoList := []string{}
	for _, v := range paramList[0] {
		repoList = append(repoList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - repo:folder
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("host: %s > repo folder not provided in paramList", hostName)
	}
	repoFolder := fmt.Sprint(paramList[1][0])

	// 2 - manage goroutines concurrency
	nbItem := len(repoList)
	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item (node)
	for _, repoName := range repoList {
		wgHost.Add(1) // Increment the WaitGroup:counter for each item
		logger.Infof("↪ (%s) %s/%s > running", phaseName, hostName, repoName)
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer func() {
				logger.Infof("↩ (%s) %s/%s > finished", phaseName, hostName, oneItem)
				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
			}()
			// defer wgHost.Done()                                                     // Decrement the WaitGroup counter - when the goroutine complete
			_, grErr := lgit.MergeDevToMain(hostName, repoFolder, repoName, logger) // the code to be executed by the goroutine
			if grErr != nil {
				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
				// send goroutines error if any into the chanel
				errChItem <- fmt.Errorf("%w", grErr)
			}

		}(repoName) // pass the node to the goroutine
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
