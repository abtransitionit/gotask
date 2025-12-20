package k8s

import (
	"github.com/abtransitionit/gocore/logx"
)

// Description: Add helm repo to a Helm client
func ConfigureCilium(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("ConfigureCilium called")
	// 1 - get parameters
	// 11 - List helm repo
	logger.Info("ConfigureCilium : get list repo")

	// handle success
	return true, nil
}
