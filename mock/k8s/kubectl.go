package k8s

import (
	"github.com/abtransitionit/gocore/logx"
)

func ConfigureKubectl(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("ConfigureKubectl called")
	// 1 - get parameters
	// 11 - control plane node
	logger.Info("ConfigureKubectl : get cplane node")
	// 2 - do the job
	logger.Info("ConfigureKubectl : get the initial config from the cplane ([]byte)")
	logger.Info("ConfigureKubectl : install it on the kubectl node (~/.kube/config)")

	// handle success
	return true, nil
}
