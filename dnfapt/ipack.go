// File: gotask/dnfapt/upgrade.go
package dnfapt

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/dnfapt"
	"github.com/abtransitionit/golinux/property"
)

func installlistDaPackOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listDaPack dnfapt.SliceDaPack) (string, error) {
	// log
	logger.Debugf("%s: will install following dnfapt package: %s", vmName, listDaPack.GetListName())

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// loop over each cli
	for _, daPkg := range listDaPack {

		// Get the CLI to install the dnfapt package
		cli, err := dnfapt.InstallDaPackage(osFamily, daPkg)
		if err != nil {
			return "", err
		}

		// play the CLI
		logger.Debugf("%s:%s installing package: %s", vmName, osFamily, daPkg.Name)
		_, err = run.RunOnVm(vmName, cli)
		if err != nil {
			return "", err
		}

	}
	return "", nil
}

// Name: createSliceFunc
//
// Description: create the slice of tasks
//
// Parameters:
// - l: logger
// - targets: list of targets
//
// Returns:
//
// - slice of syncx.Func
//
// Notes:
//
// - as many tasks as there are VMs
// - Only VM targets are included; others are skipped with a warning.
func createSliceFuncForInstallPack(ctx context.Context, logger logx.Logger, targets []phase.Target, listDaPack dnfapt.SliceDaPack) []syncx.Func {
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
			if _, err := installlistDaPackOnSingleVm(ctx, logger, vmCopy.Name(), listDaPack); err != nil {
				logger.Errorf("🅣 Failed to install Dnfapt repository on VM %s: %v", vmCopy.Name(), err)
				return err
			}

			// logger.Infof("🅣 VM %s package installed successfully", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

// name: UpgradeVmOs
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func InstallDaPackage(listDaPack dnfapt.SliceDaPack) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: UpdateVmOsApp")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : UpdateVmOsApp")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForInstallPack(ctx, logger, targets, listDaPack)

		// Log number of tasks
		logger.Infof("🅣 Phase UpdateVmOsApp has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase UpdateVmOsApp on %d VM(s)", len(tasks)), nil
	}
}
