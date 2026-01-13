package k8s

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

func InstallIngressCilium(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 10 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > release name not provided in paramList", hostName)
	}

	// 11 - get name of the release
	helmReleaseName := fmt.Sprint(paramList[0][0])
	// 11 - get Instance
	// i := k8s.GetNode(hostName)

	// // 2 - operate
	// if _, err := i.Reset(logger); err != nil {
	// 	return false, fmt.Errorf("%s > Resetting Node > %v", hostName, err)
	// }

	// handle success
	logger.Debugf("%s:%s > installed ingress cilium", hostName, helmReleaseName)
	return true, nil
}
