package sys

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/osservice"
)

func Start(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of service
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > serviceList not provided in paramList", hostName)
	}
	serviceSlice, err := filex.GetVarStructFromYaml[osservice.ServiceSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage error reporting
	nbItem := len(serviceSlice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range serviceSlice {
		// get instance
		osService := osservice.GetService(item.Name)
		// operate
		if _, err := osService.Start(hostName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("enabling os service %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb service starting that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}

func Enable(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of service
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > serviceList not provided in paramList", hostName)
	}
	serviceSlice, err := filex.GetVarStructFromYaml[osservice.ServiceSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage error reporting
	nbItem := len(serviceSlice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range serviceSlice {
		// get instance
		osService := osservice.GetService(item.Name)
		// operate
		if _, err := osService.Enable(hostName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("enabling os service %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb service enabling that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}

func Install(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 11 - list of service
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > serviceList not provided in paramList", hostName)
	}
	serviceSlice, err := filex.GetVarStructFromYaml[osservice.ServiceSlice](paramList[0])
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage error reporting
	nbItem := len(serviceSlice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range serviceSlice {
		// get instance
		osService := osservice.GetService(item.Name)
		// operate
		if _, err := osService.Install(hostName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("installing os service %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb service installing that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}
