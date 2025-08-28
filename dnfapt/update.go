// File: gotask/dnfapt/upgrade.go
package dnfapt

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/dnfapt"
	"github.com/abtransitionit/golinux/property"
)

// Name: upgradeSingleVm
//
// Description: the single task: update the OS with standard/required/missing dnfapt packages
//
// Parameters:
//
// - vmName: name of the VM
//
// Returns:
// - nil if the VM is reachable,
// - an error if the VM is not configured, not reachable or if there was an SSH failure.
//
// Notes:
// - pure logic : no logging
func updateSingleVmOsApp(logger logx.Logger, vmName string, requiredPackages []string) (string, error) {

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	logger.Debugf("will install required or missing packages: %v\n", requiredPackages)
	logger.Debugf("%s:%s\n", vmName, osFamily)

	for _, pkgName := range requiredPackages {
		install := false
		// logic for installtion
		switch pkgName {
		case "uidmap":
			if osFamily == "debian" {
				install = true
			}
		case "gnupg":
			install = true
		}

		// logic for log
		if install {
			err := dnfapt.InstallPackage(logger, osFamily, pkgName)
			if err != nil {
				return "", err
			}
			// run installation here
		} else {
			logger.Debugf("Skipping package installation for %s:%s:%s", vmName, osFamily, pkgName)
		}
	}
	return "", nil
}

// NAme: createSliceFunc
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
func createSliceFuncForUpdate(logger logx.Logger, targets []phase.Target, requiredPackages []string) []syncx.Func {
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
			if _, err := updateSingleVmOsApp(logger, vmCopy.Name(), requiredPackages); err != nil {
				logger.Errorf("🅣 Failed to install VM %s: %v", vmCopy.Name(), err)
				return err
			}

			logger.Infof("🅣 VM %s package installed successfully", vmCopy.Name())
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
func UpdateVmOsApp(listRequiredPackage []string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: UpdateVmOsApp")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : UpdateVmOsApp")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForUpdate(logger, targets, listRequiredPackage)

		// Log number of tasks
		logger.Infof("🅣 Phase UpdateVmOsApp has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase UpdateVmOsApp on %d VM(s)", len(tasks)), nil
	}
}
