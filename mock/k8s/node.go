package k8s

import (
	"github.com/abtransitionit/gocore/logx"
)

func ResetNode(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("ResetNode called")
	// handle success
	return true, nil
}

// func ResetCP(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
// 	// log
// 	logger.Info("ResetCP called")
// 	// handle success
// 	return true, nil
// }

func InitCplane(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("InitCP called")
	// 1 - get parameters
	// 11 - cluster config
	logger.Info("InitCP : get cluster config")
	// 12 - Container runtilme config
	logger.Info("InitCP : get container runtime config")

	// handle success
	return true, nil
}

// func ResetWorker(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
// 	// log
// 	logger.Info("ResetWorker called")
// 	// handle success
// 	return true, nil
// }

func AddWorker(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Info("AddWorker called")
	// handle success
	return true, nil
}
