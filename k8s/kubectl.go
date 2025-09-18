package k8s

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/golinux/k8s"
)

func ConfigureKubectlOnCplane(cPlane phase.Target) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		vmName := cPlane.Name()
		// log
		logger.Debugf("%s: will configure Kubectl on this K8s control plane", vmName)

		// get cli to configure kubectl
		cli := k8s.ConfigureKubectlOnCPlane()

		// play the cli on the control plane
		output, err := run.RunCliSsh(vmName, cli)
		if err != nil {
			fmt.Println(output)
			return "", fmt.Errorf("%s: failed to configure Kubectl on this K8s control plane", vmName)
		}
		// fmt.Println(output)
		logger.Debugf("%s: configured Kubectl on this K8s control plane", vmName)

		// fmt.Println(cli)
		return "", nil
	}
}
