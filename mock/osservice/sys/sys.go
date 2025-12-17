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
		// do nothing from now : only log
		logger.Debugf("(%s) %s:%s > ongoing", phaseName, hostName, item.Name)
		// send error if any into the chanel
		errChItem <- fmt.Errorf("%w", nil)
		// logger.Infof("(%s) %s%s > finished", phaseName, hostName, item.Name)
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

	// log
	logger.Infof("Start called with param: %v", serviceSlice)
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
		// do nothing from now : only log
		logger.Debugf("(%s) %s:%s > ongoing", phaseName, hostName, item.Name)
		// send error if any into the chanel
		errChItem <- fmt.Errorf("%w", nil)
		// logger.Infof("(%s) %s%s > finished", phaseName, hostName, item.Name)
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

	// log
	logger.Infof("Enable called with param: %v", serviceSlice)
	// handle success
	return true, nil
}
