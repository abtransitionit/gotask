package onpm

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	lonpm "github.com/abtransitionit/golinux/mock/onpm"
	"gopkg.in/yaml.v3"
)

// Description: add native os packages to a Linux host
func AddPkg(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// logger.Debugf("paramList: %v", paramList)
	// 1 - get parameters
	// 11 - package:list
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > list of package not provided in paramList", hostName)
	}
	raw := paramList[0]

	b, err := yaml.Marshal(raw)
	if err != nil {
		return false, err
	}
	pkgList, err := filex.GetVarStructFromYamlString[lonpm.PkgSlice](fmt.Sprint(string(b)))
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage error reporting
	nbItem := len(pkgList)
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item
	for _, item := range pkgList {
		_, grErr := lonpm.AddPkg(hostName, item.Name, logger) // the code to be executed
		if grErr != nil {
			logger.Errorf("%s:%s (%s) > %v", hostName, item.Name, phaseName, grErr) // send goroutines error if any into the chanel
			// send error if any into the chanel
			errChItem <- fmt.Errorf("%w", grErr)
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
		logger.Errorf("❌ %s > nb pkg that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// 6 - handle success
	return true, nil
}
