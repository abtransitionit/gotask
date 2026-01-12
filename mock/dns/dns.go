package dns

import (
	"github.com/abtransitionit/gocore/logx"
	ldns "github.com/abtransitionit/golinux/mock/dns"
)

// Description: create missing files
//
// Notes:
//
// - this fix is only necessary by certain tools
//   - mandatory to run K8S chart cilium:1.18.5
func FixDns(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// log
	// play CLI
	if err := ldns.FixDns(hostName, logger); err != nil {
		logger.Warnf("%s > fixing dns > %w", hostName, err)
	}

	// handle success
	return true, nil
}

// func FixDns(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

// 	// 1 - get parameters
// 	// 11 - node:List
// 	nodeList := []string{}
// 	for _, v := range paramList[0] {
// 		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
// 	}
// 	// 12 - fileProperty
// 	if len(paramList) < 2 || len(paramList[1]) == 0 {
// 		return false, fmt.Errorf("host: %s > fileProperty not provided in paramList", hostName)
// 	}
// 	fileProperty, err := cfile.GetVarStructFromYamlString[lfile.FileProperty](fmt.Sprint(paramList[1][0]))
// 	if err != nil {
// 		logger.Errorf("%v", err)
// 	}

// 	// 2 - manage goroutines concurrency
// 	nbItem := len(nodeList)
// 	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
// 	errChNode := make(chan error, nbItem) // define a channel to collect errors from each goroutine

// 	// 3 - loop over item (node)
// 	for _, node := range nodeList {
// 		wgHost.Add(1)             // Increment the WaitGroup:counter for each item
// 		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
// 			defer func() {
// 				logger.Infof("↩ (%s) %s/%s > finished", phaseName, hostName, oneItem)
// 				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
// 			}()
// 			logger.Infof("↪ (%s) %s/%s > ongoing", phaseName, hostName, oneItem)
// 			_, grErr := dns.CopyFileWithSudo(hostName, node, fileProperty, logger) // the code to be executed by the goroutine
// 			if grErr != nil {
// 				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
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

// // Description: create a RC custom file on a set of nodes from a hostname
// func CreateRcFile(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

// 	// 1 - get parameters
// 	// 11 - node:List
// 	nodeList := []string{}
// 	for _, v := range paramList[0] {
// 		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
// 	}
// 	// 12 - file path
// 	if len(paramList) < 2 || len(paramList[1]) == 0 {
// 		return false, fmt.Errorf("host: %s > RC 	file path not provided in paramList", hostName)
// 	}
// 	fileName := fmt.Sprint(paramList[1][0])

// 	// 2 - manage goroutines concurrency
// 	nbItem := len(nodeList)
// 	var wgHost sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
// 	errChNode := make(chan error, nbItem) // define a channel to collect errors from each goroutine

// 	// 3 - loop over item (node)
// 	for _, node := range nodeList {
// 		wgHost.Add(1)             // Increment the WaitGroup:counter for each item
// 		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
// 			defer func() {
// 				logger.Infof("↩ (%s) %s/%s > finished", phaseName, hostName, oneItem)
// 				wgHost.Done() // Decrement the WaitGroup counter - when the goroutine complete
// 			}()
// 			logger.Infof("↪ (%s) %s/%s > ongoing", phaseName, hostName, oneItem)
// 			grErr := lfile.ForceCreateRcFile(hostName, node, fileName, logger) // the code to be executed by the goroutine
// 			if grErr != nil {
// 				logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, oneItem, grErr) // send goroutines error if any into the chanel
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
