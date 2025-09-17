package k8s

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/syncx"
)

func InitSingleWorker(ctx context.Context, logger logx.Logger, vmName string) (string, error) {
	// log
	logger.Debugf("%s: will initialize This VM as a K8s Worker", vmName)

	return "", nil
}

func createSliceFuncForInitListWorker(ctx context.Context, logger logx.Logger, targets []phase.Target) []syncx.Func {
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
			if _, err := InitSingleWorker(ctx, logger, vmCopy.Name()); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func InitWorker(targetsCPlane []phase.Target) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "InitCplane"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targetsCPlane) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForInitListWorker(ctx, logger, targetsCPlane)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
