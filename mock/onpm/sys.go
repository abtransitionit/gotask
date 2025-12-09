package onpm

import (
	"github.com/abtransitionit/gocore/logx"
	lonpm "github.com/abtransitionit/golinux/mock/onpm"
)

// Description: updates the linux OS of a host
//
// Notes:
// - update mean: add the missing/reuired standard native OS package repositories and packages
func UpdateOs(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// log
	// play CLI
	_, err := lonpm.UpdateOs(hostName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("host: %s > system error > updating OS > %v", hostName, err)
	}

	// handle success
	return true, nil
}

// Description: upgrade the linux OS of a host
//
// Notes:
// - upgrade mean: set the OS native repositories and packages to the latest version
func UpgradeOs(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	_, err := lonpm.UpgradeOs(hostName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("host: %s > system error > upgrading OS: %v", hostName, err)
	}

	// handle success
	return true, nil
}

// Description: check if a reboot is needed for a host
//
// TODO:
// - this task should not be part of the gcore/golinux:omp package but rather gcore/golinux:node package
func NeedReboot(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	out, err := lonpm.NeedReboot(hostName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("host: %s > system error > getting reboot status: %v", hostName, err)
	}

	// handle success
	logger.Debugf("host: %s > need reboot: %s", hostName, out)
	return true, nil
}
