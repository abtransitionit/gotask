package k8s

import (
	"github.com/abtransitionit/gocore/logx"
)

// Description: Add helm repo to a Helm client
func AddRepoHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddRepoHelm called")
	// 1 - get parameters
	// 11 - List helm repo
	logger.Info("AddRepoHelm : get list repo")

	// handle success
	return true, nil
}

// Description: Add helm charts into to a K8s cluster from a Helm client
func AddChartHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddChartHelm called")
	// 1 - get parameters
	// 11 - List helm chart
	logger.Info("AddChartHelm : get list chart")

	// handle success
	return true, nil
}
