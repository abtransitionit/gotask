package gopm

import (
	"github.com/abtransitionit/gocore/logx"
)

func InstallGoCli(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("InstallGoCli called")
	// 1 - get parameters
	// 11 - List of go cli to install
	logger.Info("InstallGoCli : get list of go cli to install")

	// handle success
	return true, nil
}
