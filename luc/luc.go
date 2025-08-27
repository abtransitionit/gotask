// File gotask/phase/show.go
package luc

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/cli"
)

// Name: checkSingleVmIsSshReachable
//
// Description: the single task: deploys LUC on a single VM
//
// Returns:
// - nil if the VM is reachable,
// - an error if the VM is not configured, not reachable or if there was an SSH failure.
//
// Prerequisites:
//
// - VM must be reachable via SSH
//
// Notes:
// - pure logic : no logging
func DeployLucOnSingleVm(logger logx.Logger, vm *phase.Vm) error {
	// define vars
	localArtifactPath := "/tmp/goluc-linux"
	dstPath := "/usr/local/bin/goluc"
	hostFilePath := fmt.Sprintf("%s:%s", vm.Name(), dstPath)

	// deploy the artifact : ie. scp fie to remote
	deployOk, err := cli.DeployGoArtifactAsSudo(logger, localArtifactPath, hostFilePath)
	if err != nil {
		return errorx.NewWithNoStack("🅣 %s: %v ", err, vm.Name())
	}
	if !deployOk {
		return errorx.NewWithNoStack("🅣 %s: %v ", err, vm.Name())
	}
	return nil
}

// NAme: createSliceFunc
//
// Description: create the slice of tasks
//
// Parameters:
// - logger: logger
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
func createSliceFunc(logger logx.Logger, targets []phase.Target) []syncx.Func {

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
			if err := DeployLucOnSingleVm(logger, vmCopy); err != nil {
				logger.Errorf("%s", err)
				return err
			}
			return nil
		})
	}

	return tasks
}

// name: DeployLuc
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func DeployLuc(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

	// check paramaters
	if len(targets) == 0 {
		logger.Warn("🅣 No targets provided to checkVmSshAccess")
		return "", nil
	}

	// Build slice of functions
	tasks := createSliceFunc(logger, targets)

	// Log number of tasks
	logger.Infof("🅣 Phase has %d concurent tasks", len(tasks))

	// Run tasks in the slice concurrently
	if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
		return "", errs[0] // return first error encountered
	}

	return fmt.Sprintf("🅣 SSH successfully verified on %d VM(s)", len(tasks)), nil
}
