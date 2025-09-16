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

func InstallLisOsServiceOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsServices oservice.SliceOsService) (string, error) {
	// log
	logger.Debugf("%s: will install following OS service(s): %s", vmName, listOsServices.GetListName())

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// loop over each cli
	for _, osService := range listOsServices {

		// Get the cli to install the service
		cli, err := osService.Install(osFamily)
		if err != nil {
			return "", err
		}
		if cli == "" {
			continue
		}

		// play it on remote
		// logger.Debugf("%s: creating service file %s", vmName, osService.Path)
		_, err = run.RunCliSsh(vmName, cli)
		if err != nil {
			return "", fmt.Errorf("%s: failed to play cli  on vm : '%s': %w", vmName, cli, err)
		}
		logger.Debugf("%s: created service file %s", vmName, osService.Path)
	}
	return "", nil
}

func createSliceFuncForOsService(ctx context.Context, logger logx.Logger, targets []phase.Target, listOsServices oservice.SliceOsService) []syncx.Func {
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
			if _, err := InstallLisOsServiceOnSingleVm(ctx, logger, vmCopy.Name(), listOsServices); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func InstallOsService(listOsServices []oservice.OsService) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "InstallOsService"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForOsService(ctx, logger, targets, listOsServices)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
