package oservice

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/oservice"
	"github.com/abtransitionit/golinux/property"
)

func EnableLingerOnSingleVm(ctx context.Context, logger logx.Logger, vmName string) (string, error) {

	// get property
	currentUser, err := property.GetProperty(vmName, "osuser")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	logger.Debugf("Enable linger for current user: %s", currentUser)
	// get cli
	cli := oservice.EnableLinger()

	// play it on remote
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	return "", nil
}

func createSliceFuncForLinger(ctx context.Context, logger logx.Logger, targets []phase.Target) []syncx.Func {
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
		// define the job of the task and add it to the slice
		tasks = append(tasks, func() error {
			if _, err := EnableLingerOnSingleVm(ctx, logger, vmCopy.Name()); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func EnableLinger(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
	appx := "EnableLinger"
	logger.Infof("🅣 Starting phase: %s", appx)
	// check paramaters
	if len(targets) == 0 {
		logger.Warnf("🅣 No targets provided to phase: %s", appx)
		return "", nil
	}

	// Build slice of functions
	tasks := createSliceFuncForLinger(ctx, logger, targets)

	// Log number of tasks
	logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

	// Run tasks in the slice concurrently
	if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
		return "", errs[0] // return first error encountered
	}

	return "", nil
}
