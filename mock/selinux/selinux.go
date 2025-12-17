package selinux

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/selinux"
)

func ConfigureSelinux(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	_, err := selinux.ConfigureSelinux(hostName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("host: %s > system error > upgrading OS: %v", hostName, err)
	}

	// log
	// logger.Infof("ConfigureSelinux called ")

	// handle success
	return true, nil

}
