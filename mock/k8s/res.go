package k8s

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/k8scli/kubectl"
)

func CreateSecret(phaseName, kubectlHost string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// check
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("local:%s > list of secrets not properly provided in paramList", kubectlHost)
	}
	// 11 - slice to managed
	slice, err := filex.GetVarStructFromYaml[kubectl.SliceResource](paramList[0])
	if err != nil {
		return false, fmt.Errorf("local:%s > getting slice from paramList: %w", kubectlHost, err)
	}

	// 2 - manage error reporting
	nbItem := len(slice)
	logger.Debugf("local:%s > nb item: %d > slice %+v", kubectlHost, nbItem, slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range slice {
		// 31 - get instance and operate
		i := kubectl.Resource{Type: kubectl.ResSecret, Name: item.Name, Ns: item.Ns, UserName: item.UserName}
		if _, err := i.Create("local", kubectlHost, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("creating secret %s > %w", item.Name, err)
		}
	} // loop

	// 4 - manage error
	close(errChItem) // close the channel - signal that no more error will be sent
	// 41 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 42 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ local:%s > nb item installation that failed: %d", kubectlHost, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}
func ApplyManifest(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 10 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > list manifet or helm client not properly provided in paramList", hostName)
	}
	// 11 - name of helm client host
	kubectlClientNodeName := fmt.Sprint(paramList[1][0])

	// 12 - slice to managed
	slice, err := filex.GetVarStructFromYaml[kubectl.SliceResource](paramList[0])
	if err != nil {
		return false, fmt.Errorf("%s > getting slice from paramList: %w", hostName, err)
	}

	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range slice {
		// 31 - get instance and operate
		i := kubectl.Resource{Type: kubectl.ResManifest, Name: item.Name}
		if _, err := i.Apply("local", kubectlClientNodeName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("applying manifest %s > %w", item.Name, err)
		}
	} // loop

	// 4 - manage error
	close(errChItem) // close the channel - signal that no more error will be sent
	// 41 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 42 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ %s > nb item installation that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}

// func AddNs(phaseName, kubectlHost string, paramList [][]any, logger logx.Logger) (bool, error) {
// 	// 1 - get parameters
// 	// check
// 	if len(paramList) < 1 || len(paramList[0]) == 0 {
// 		return false, fmt.Errorf("local:%s > list of GO clis or Destination folder not properly provided in paramList", kubectlHost)
// 	}
// 	// 11 - List of namespaces to create
// 	slice, err := filex.GetVarStructFromYaml[kubectl.SliceResource](paramList[0])
// 	if err != nil {
// 		return false, fmt.Errorf("local:%s > getting cliName from paramList: %w", kubectlHost, err)
// 	}

// 	// 2 - manage error reporting
// 	nbItem := len(slice)
// 	errChItem := make(chan error, nbItem) // define a channel to collect errors
// 	logger.Debugf("slice %v", slice)

// 	// 3 - loop over item
// 	for _, item := range slice {
// 		// 31 - get instance and operate
// 		logger.Debugf("item is %+v", item)
// 		i := kubectl.Resource{Type: kubectl.ResNS, Name: item.Name}
// 		if _, err := i.Create("local", kubectlHost, logger); err != nil {
// 			// send error if any into the chanel
// 			errChItem <- fmt.Errorf("installing GO cli %s: %w", item.Name, err)
// 		}
// 	} // loop

// 	// 4 - manage error
// 	close(errChItem) // close the channel - signal that no more error will be sent
// 	// 41 - collect errors
// 	var errList []error
// 	for e := range errChItem {
// 		errList = append(errList, e)
// 	}

// 	// 42 - handle errors
// 	nbGroutineFailed := len(errList)
// 	errCombined := errors.Join(errList...)
// 	if nbGroutineFailed > 0 {
// 		logger.Errorf("❌ local:%s > nb cli installation that failed: %d", kubectlHost, nbGroutineFailed)
// 		return false, errCombined
// 	}

// 	// handle success
// 	return true, nil
// }
