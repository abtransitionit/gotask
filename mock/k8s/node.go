package k8s

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/k8s"
)

func ResetNode(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get Instance
	i := k8s.GetNode(hostName)

	// 2 - operate
	if err := i.Reset(logger); err != nil {
		return false, fmt.Errorf("%s > Resetting Node > %v", hostName, err)
	}

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
	// 1 - get parameters
	// 11 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > cluster config not provided in paramList", hostName)
	}
	// 12 - get cluster config
	clusterParam, err := filex.GetVarStructFromYamlString[k8s.ClusterParam](fmt.Sprint(paramList[0][0]))
	if err != nil {
		return false, fmt.Errorf("getting cluster config from paramList: %v", err)
	}
	// 13 - get container runtime config AND put it directly in the clusterParam
	clusterParam.CrSocketName = fmt.Sprint(paramList[1][0])

	// 2 - get Instance
	i := k8s.GetCplane(hostName)

	// 3 - operate
	if err := i.Init(clusterParam, logger); err != nil {
		return false, fmt.Errorf("%s > initializing control plane > %v", hostName, err)
	}

	// handle success
	return true, nil

}

func AddWorker(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get Instance
	i := k8s.GetWorker(hostName)

	// 2 - operate
	if err := i.Add(logger); err != nil {
		// handle error
		return false, fmt.Errorf("%s > configuring selinux > %v", hostName, err)
	}
	// handle success
	return true, nil

}
