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
func installListStdPkgOnSingleVmOs(logger logx.Logger, vmName string, requiredPackages []string) (string, error) {

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	// get property
	osDistro, err := property.GetProperty(vmName, "osdistro")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// loop over each package
	logger.Debugf("will install required or missing packages: %v\n", requiredPackages)
	for _, pkgName := range requiredPackages {

		// get the CLI
		cli, err := dnfapt.InstallStdDaPackage(osFamily, osDistro, pkgName)
		if err != nil {
			return "", err
		}
		// if nothing is to be installed
		if cli == "" {
			logger.Debugf("%s:%s:%s Skipping dnfapt package installation of : %s", vmName, osFamily, osDistro, pkgName)
			continue
		}

		// play the CLI
		logger.Debugf("%s:%s:%s Installing dnfapt package : %s", vmName, osFamily, osDistro, pkgName)
		_, err = run.RunOnVm(vmName, cli)
		if err != nil {
			return "", fmt.Errorf("failed to install package %s OS on VM %s: %w", pkgName, vmName, err)
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
func createSliceFuncForStdInstall(logger logx.Logger, targets []phase.Target, requiredPackages []string) []syncx.Func {
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
			if _, err := installListStdPkgOnSingleVmOs(logger, vmCopy.Name(), requiredPackages); err != nil {
				logger.Errorf("🅣 Failed to install VM %s: %v", vmCopy.Name(), err)
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
func UpdateVmOsApp(listRequiredPackage []string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "UpdateVmOsApp"

		// log
		logger.Infof("🅣 Starting phase: %s", appx)

		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForStdInstall(logger, targets, listRequiredPackage)

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
