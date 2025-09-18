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

func AddSingleWorker(ctx context.Context, logger logx.Logger, vmName string, cPlane phase.Target) (string, error) {
	// define var
	cPlaneVm := cPlane.Name()
	workerVm := vmName
	// log
	logger.Debugf("%s: will Add This VM as a K8s Worker to the K8s cluster using the join cli from the CPlane %s", workerVm, cPlaneVm)

	// get the cli to get the join CLI from the control plane
	cli := k8s.GetJoinCli()
	// print(cli)

	// play the cli on the control plane to get the join CLI
	joinCli, err := run.RunCliSsh(cPlaneVm, cli)
	if err != nil {
		fmt.Println(joinCli)
		return "", fmt.Errorf("%s: failed to get the join CLI from the control plane %s", workerVm, cPlaneVm)
	}
	fmt.Println(joinCli)

	// here we got the join cli - build the cli to add the worker
	cli, err = k8s.AddWorker(joinCli)
	// cli, err = k8s.AddWorkerWithReset(joinCli)
	if err != nil {
		return "", err
	}
	// play the join cli on the worker
	output, err := run.RunCliSsh(workerVm, cli)
	if err != nil {
		fmt.Println(output)
		return "", fmt.Errorf("%s: failed to add the worker to the cluster", workerVm)
	}
	// success
	fmt.Println(output)
	logger.Infof("%s: added the worker to the cluster.", workerVm)
	return "", nil
}

func createSliceFuncForAddListWorker(ctx context.Context, logger logx.Logger, targets []phase.Target, cPlane phase.Target) []syncx.Func {
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
			if _, err := AddSingleWorker(ctx, logger, vmCopy.Name(), cPlane); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func AddWorker(cPlane phase.Target, targetsWorker []phase.Target) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "AddWorker"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targetsWorker) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForAddListWorker(ctx, logger, targetsWorker, cPlane)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
