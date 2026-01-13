package k8s

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	helm2 "github.com/abtransitionit/golinux/mock/k8scli/helm"
)

func InstallIngressCilium(phaseName, helmHost string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 10 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("local:%s > release name or namespace not properly provided in paramList", helmHost)
	}
	// 11 - get release to upgrade
	release, err := filex.GetVarStructFromYaml[helm2.Release](paramList[0][0])
	if err != nil {
		return false, fmt.Errorf("local:%s > getting list from paramList: %w", helmHost, err)
	}
	// 2 - get Instance
	i := helm2.GetRelease(release.Name, strings.TrimSpace(release.Chart.QName), "", release.Namespace, nil)

	// 3 - operate
	if err := i.InstallIngressCilium("local", helmHost, logger); err != nil {
		return false, fmt.Errorf("local:%s > installing ingress cilium to existing release > %w", helmHost, err)
	}

	// handle success
	return true, nil
}
