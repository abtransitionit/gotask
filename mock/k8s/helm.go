package k8s

import (
	"github.com/abtransitionit/gocore/logx"
)

// Description: Add helm repo to a Helm client
//
// Notes:
//
// - it just adds the repo to the Helm client, it does not install any chart
func AddRepoHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddRepoHelm : adding repo: TODO")
	// 1 - get parameters
	// 11 - List helm repo

	// handle success
	return true, nil
}

// Description: Add helm charts into a K8s cluster from a Helm client
//
// Notes:
//
// - the helm repo of the chart must be presents in the Helm client
func AddChartHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddChartHelm : adding chart: TODO")
	// 1 - get parameters
	// 11 - List helm chart

	// handle success
	return true, nil
}

// Description: Add helm releases into a K8s cluster from a Helm client
//
// Notes:
//
// - the helm repo of the chart must be presents in the Helm client
func AddReleaseHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddReleaseHelm : adding release: TODO")
	// 1 - get parameters
	// 11 - List helm release

	// handle success
	return true, nil
}
