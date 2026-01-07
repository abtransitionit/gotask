package k8s

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/k8s"
)

func ConfigureKubectl(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// check
	if len(paramList) < 2 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > control plane and install nodes not properly provided in paramList", hostName)
	}
	// 11 - get control plane node
	cplaneNode := fmt.Sprint(paramList[0][0])
	// 12 - get install node
	installNode := fmt.Sprint(paramList[1][0])

	// 2 - get Instance
	i := k8s.GetKubectl(cplaneNode, installNode)

	// // 3 - operate
	if err := i.Configure(hostName, logger); err != nil {
		return false, fmt.Errorf("%s:%s > configuring kubectl > %w", hostName, i.InstallNodeName, err)
	}

	// handle success
	return true, nil
}
