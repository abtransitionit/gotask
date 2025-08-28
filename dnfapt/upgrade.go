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

// Name: upgradeSingleVmOs
//
// Description: the single task: upgrade an OS
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
func upgradeSingleVmOs(vmName string) (string, error) {

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// get the cli
	cli, err := dnfapt.UpgradeOs(osFamily)
	if err != nil {
		return "", err
	}
	// play the CLI
	_, err = run.RunOnVm(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to upgrade OS on VM %s: %w", vmName, err)
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
func createSliceFuncForUpgrade(l logx.Logger, targets []phase.Target) []syncx.Func {
	var tasks []syncx.Func

	for _, t := range targets {
		if t.Type() != "Vm" {
			continue
		}

		vm, ok := t.(*phase.Vm)
		if !ok {
			l.Warnf("🅣 Target %s is not a VM, skipping", t.Name())
			continue
		}

		vmCopy := vm // capture for closure
		tasks = append(tasks, func() error {
			if _, err := upgradeSingleVmOs(vmCopy.Name()); err != nil {
				l.Errorf("🅣 Failed to upgrade VM %s: %v", vmCopy.Name(), err)
				return err
			}

			l.Infof("🅣 VM %s upgraded successfully", vmCopy.Name())
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
func UpgradeVmOs(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

	logger.Info("🅣 Starting phase: UpgradeVmOs")
	// check paramaters
	if len(targets) == 0 {
		logger.Warn("🅣 No targets provided to : UpgradeVmOs")
		return "", nil
	}

	// Build slice of functions
	tasks := createSliceFuncForUpgrade(logger, targets)

	// Log number of tasks
	logger.Infof("🅣 Phase UpgradeVmOs has %d concurent tasks", len(tasks))

	// Run tasks in the slice concurrently
	if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
		return "", errs[0] // return first error encountered
	}

	return fmt.Sprintf("🅣 Terminated phase UpgradeVmOs on %d VM(s)", len(tasks)), nil
}
