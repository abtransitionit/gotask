package oskernel

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/oskernel"
)

// description: Load a list of kernel modules
func LoadModule(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// check
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > Module list not provided in paramList", hostName)
	}
	// 11 - list of KModule
	slice, err := filex.GetVarStructFromYaml[oskernel.ModuleSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}
	// check
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > file.kernel not provided in paramList", hostName)
	}
	// 12 - kernel config file
	kernelFileName := fmt.Sprint(paramList[1][0])

	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - get instance
	i := oskernel.GetModuleSet(slice, kernelFileName)

	// 4 - operate
	if _, err := i.Load(hostName, logger); err != nil {
		// send error if any into the chanel
		errChItem <- fmt.Errorf("adding kernel module %s: %w", slice[0].Name, err)
	}

	// 5 - manage error
	close(errChItem) // close the channel - signal that no more error will be sent
	// 51 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 52 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ %s > nb module loading that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}

// description: Load a list of kernel module parameters
func LoadParam(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of kParameter
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > Parameter list not provided in paramList", hostName)
	}
	slice, err := filex.GetVarStructFromYaml[oskernel.ParameterSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}
	// 12 - kernel config file for module
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > file.kernel not provided in paramList", hostName)
	}
	kernelFileName := fmt.Sprint(paramList[1][0])

	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - get instance
	i := oskernel.GetParameterSet(slice, kernelFileName)

	// 4 - operate
	if err := i.Load(hostName, logger); err != nil {
		// send error if any into the chanel
		errChItem <- fmt.Errorf("adding kernel parameter: %w", err)
	}

	// 5 - manage error
	close(errChItem) // close the channel - signal that no more error will be sent
	// 51 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 52 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ %s > nb parameter loading that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}
