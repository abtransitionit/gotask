// File gotask/phase/show.go
package vm

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
)

// Name: checkSingleVmIsSshReachable
//
// Description: the single task: checks if a single VM is SSH reachable
//
// Returns:
// - nil if the VM is reachable,
// - an error if the VM is not configured, not reachable or if there was an SSH failure.
//
// Notes:
// - pure logic : no logging
func checkSingleVmIsSshReachable(vm *phase.Vm, logger logx.Logger) error {
	reachable, err := run.IsVmSshReachable(vm.Name())
	if err != nil {
		return errorx.NewWithNoStack("🅣 SSH check failed for VM %s", vm.Name())
	}
	if !reachable {
		return errorx.NewWithNoStack("🅣 VM %s is not reachable via SSH", vm.Name())
	}
	logger.Infof("%s: VM is ssh reachable", vm.NameStr)
	return nil
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
func createSliceFunc(targets []phase.Target, logger logx.Logger) []syncx.Func {

	var tasks []syncx.Func // the slice

	// loop over targets
	for _, t := range targets {
		// skip non-VM targets
		if t.Type() != "Vm" {
			continue
		}
		// get the vm struct
		vm, ok := t.(*phase.Vm)
		if !ok {
			logger.Warnf("🅣 Target %s is not a VM, skipping", t.Name())
			continue
		}

		vmCopy := vm // capture for closure
		// define and add each task to the slice
		tasks = append(tasks, func() error {
			if err := checkSingleVmIsSshReachable(vmCopy, logger); err != nil {
				logger.Errorf("%s", err)
				return err
			}
			// logger.Debugf("VM %s passed SSH check", vmCopy.Name()) // log success
			return nil
		})
	}

	return tasks
}

// name: checkVmSshAccess
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func CheckVmSshAccess(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
	// check paramaters
	if len(targets) == 0 {
		logger.Warn("🅣 No targets provided to checkVmSshAccess")
		return "", nil
	}

	// Build slice of functions
	tasks := createSliceFunc(targets, logger)

	// Log number of tasks
	logger.Infof("🅣 Phase has %d concurent tasks", len(tasks))

	// Run tasks in the slice concurrently
	if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
		return "", errs[0] // return first error encountered
	}

	return fmt.Sprintf("🅣 SSH successfully verified on %d VM(s)", len(tasks)), nil
}
