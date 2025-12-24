package oskernel

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/oskernel"
)

func AddKModule(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of KModule
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > Module list not provided in paramList", hostName)
	}
	moduleSlice, err := filex.GetVarStructFromYaml[oskernel.ModuleSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}
	// 12 - kernel config file for module
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > file.kernel not provided in paramList", hostName)
	}
	kernelFileName := fmt.Sprint(paramList[1][0])

	// 2 - manage error reporting
	nbItem := len(moduleSlice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range moduleSlice {
		// get instance
		kModule := oskernel.GetModule(item.Name, kernelFileName)
		// operate
		if _, err := kModule.Add(hostName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("adding kernel module %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb module that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// log
	// logger.Infof("AddKModule called with param: %v and %s", moduleSlice, kernelFileName)
	// handle success
	return true, nil
}

func AddKParam(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of kParameter
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > Parameter list not provided in paramList", hostName)
	}
	parameterSlice, err := filex.GetVarStructFromYaml[oskernel.ParameterSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}
	// 12 - kernel config file for module
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > file.kernel not provided in paramList", hostName)
	}
	kernelFileName := fmt.Sprint(paramList[1][0])

	// 2 - manage error reporting
	nbItem := len(parameterSlice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range parameterSlice {
		// get instance
		kParameter := oskernel.GetParameter(item.Name, item.Value, kernelFileName)
		// operate
		if _, err := kParameter.Add(hostName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("adding kernel parameter %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb parameter that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// log
	// logger.Infof("AddKParam called with param: %v and %s", parameterSlice, kernelFileName)
	// handle success
	return true, nil
}
