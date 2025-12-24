package selinux

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/property"
	"github.com/abtransitionit/golinux/selinux"
)

func configureSelinuxOnSingleVm(ctx context.Context, logger logx.Logger, vmName string) (string, error) {

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// get property - before changes
	selinuxStatus, err := property.GetProperty(vmName, "selinuxStatus")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	// get property - before changes
	selinuxMode, err := property.GetProperty(vmName, "selinuxMode")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// skip if Os:family not rhel or not fedora
	if osFamily != "rhel" && osFamily != "fedora" {
		// logger.Debugf("%s:%s 🅐 Skipping selinux configuration", vmName, osFamily)
		return "", nil
	}

	// here: osFamily is in ["rhel", "fedora"]
	logger.Debugf("%s:%s    🅐 Selinux status before is : %s / %s", vmName, osFamily, selinuxStatus, selinuxMode)
	// logger.Debugf("%s:%s 🅐 configuring Selinux at startup and runtime", vmName, osFamily)
	cli := selinux.ConfigureSelinuxAtRuntime()
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	cli = selinux.ConfigureSelinuxAtStartup()
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("%s:%s  failed to play cli %s : %w", vmName, osFamily, cli, err)
	}

	// get property - after changes
	selinuxMode, err = property.GetProperty(vmName, "selinuxMode")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// get property - after changes
	selinuxStatus, err = property.GetProperty(vmName, "selinuxStatus")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// success
	logger.Debugf("%s:%s 🅑 Selinux status after is : %s / %s", vmName, osFamily, selinuxStatus, selinuxMode)
	// logger.Debugf("%s:%s 🅑 configured Selinux", vmName, osFamily)
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
func createSliceFuncForSelinux(ctx context.Context, logger logx.Logger, targets []phase.Target) []syncx.Func {
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
			if _, err := configureSelinuxOnSingleVm(ctx, logger, vmCopy.Name()); err != nil {
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
func ConfigureSelinux() phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		phaseName := "ConfigureSelinux"
		logger.Infof("🅣 Starting phase: %s", phaseName)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to: %s", phaseName)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForSelinux(ctx, logger, targets)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", phaseName, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase %s on %d VM(s)", phaseName, len(tasks)), nil
	}
}
