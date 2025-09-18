package k8s

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/k8s"
)

// Name: ResetSingleNode
//
// Description: Reset a Node (control plane or worker) to a standard VM
func ResetSingleNode(ctx context.Context, logger logx.Logger, vmName string) (string, error) {
	// log
	logger.Debugf("%s: will reset this Node (control plane or worker) to a standard VM", vmName)

	// get CLI to initialize the control plane
	cli, err := k8s.ResetNode()
	if err != nil {
		return "", err
	}

	// play the cli on the control plane
	output, err := run.RunCliSsh(vmName, cli)
	if err != nil {
		fmt.Println(output)
		return "", fmt.Errorf("%s: failed to reset this control plane to a standard VM", vmName)
	}
	// fmt.Println(output)
	logger.Debugf("%s: reset this control plane to a standard VM", vmName)
	return "", nil
}

func createSliceFuncForResetNode(ctx context.Context, logger logx.Logger, targets []phase.Target) []syncx.Func {
	var tasks []syncx.Func

	for _, target := range targets {
		if target.Type() != "Vm" {
			continue
		}

		vm, ok := target.(*phase.Vm)
		if !ok {
			logger.Warnf("🅣 Target %s is not a VM, skipping", target.Name())
			continue
		}

		vmCopy := vm // capture for closure
		// define the job of the task and add it to the slice
		tasks = append(tasks, func() error {
			if _, err := ResetSingleNode(ctx, logger, vmCopy.Name()); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func ResetNode(targetsNode []phase.Target) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "ResetNode"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targetsNode) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForResetNode(ctx, logger, targetsNode)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
