package helm

import (
	"context"

	helm "github.com/abtransitionit/gocore/k8s-helm"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/syncx"
)

func installListHelmRepoOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, sliceHelmRepo []helm.HelmRepo) (string, error) {

	logger.Debugf("provision following helm chart repositories : %s on the helm host", vmName, sliceHelmRepo)

	return "", nil
}

func createSliceFuncForInstallHelmRepo(ctx context.Context, logger logx.Logger, targets []phase.Target, sliceHelmRepo []helm.HelmRepo) []syncx.Func {
	var tasks []syncx.Func

	for _, t := range targets {
		if t.Type() != "Vm" {
			continue
		}

		vm, ok := t.(*phase.Vm)
		if !ok {
			logger.Warnf("🅣 Target %s is not a VM, skipping", t.Name())
			continue
		}

		vmCopy := vm // capture for closure
		tasks = append(tasks, func() error {
			if _, err := installListHelmRepoOnSingleVm(ctx, logger, vmCopy.Name(), sliceHelmRepo); err != nil {
				logger.Errorf("🅣 Failed to install Dnfapt repository on VM %s: %v", vmCopy.Name(), err)
				return err
			}

			// logger.Infof("🅣 VM %s package installed successfully", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

// name: InstallCniCilium
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func InstallHelmRepo(sliceHelmRepo []helm.HelmRepo) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "InstallHelmRepo"

		// log
		logger.Infof("🅣 Starting phase: %s", appx)

		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForInstallHelmRepo(ctx, logger, targets, sliceHelmRepo)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}
		// return fmt.Sprintf("🅣 Terminated phase UpdateVmOsApp on %d VM(s)", len(tasks)), nil
		return "", nil
	}

}

// 🟦 all nodes > provision CNI plugin
// provision CNI plugin: Cilium
// all nodes > provision Helm repository
// - helm repo add projectcalico https://docs.tigera.io/calico/charts
// - helm repo update
// 🟦 all nodes > create a K8s namespace
// - kubectl create namespace tigera-operator
// 🟦 all nodes > provision Calico chart
// config helm install calico projectcalico/tigera-operator --namespace tigera-operator
// 🟦 all nodes > check install
// - kubectl get pods -n tigera-operator
// - kubectl get pods -n calico-system
// - kubectl get nodes
// - kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.2/manifests/calico.yaml
// 	→ Calico DaemonSet (runs on each node)
// 	→ CNI configuration
// 	→ CRDs for network policy
// 	→ Calico controller
