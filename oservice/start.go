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

func StartListOsServiceOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsServices oservice.SliceOsService) (string, error) {
	// log
	logger.Debugf("%s: will start following OS service(s): %s", vmName, listOsServices.GetListName())

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// loop over each cli
	for _, osService := range listOsServices {

		// Get the cli to start the service
		cli, err := osService.Start(osFamily)
		if err != nil {
			return "", err
		}
		if cli == "" {
			logger.Debugf("%s: 🅐 skipping starting service %s", vmName, osService.Name)
			continue
		}

		// play it on remote
		// logger.Debugf("%s: starting service file %s", vmName, osService.Path)
		// fmt.Println(cli)
		_, err = run.RunCliSsh(vmName, cli)
		if err != nil {
			return "", fmt.Errorf("%s: failed to play cli  on vm : '%s': %w", vmName, cli, err)
		}
		// get property - before changes
		serviceStatus, err := property.GetProperty(vmName, "serviceStatus", osService.Name)
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		logger.Debugf("%s:%s:%s 🅑 started service & status is %s", vmName, osFamily, osService.Name, serviceStatus)
	}
	return "", nil
}

func createSliceFuncForStartOsService(ctx context.Context, logger logx.Logger, targets []phase.Target, listOsServices oservice.SliceOsService) []syncx.Func {
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
			if _, err := StartListOsServiceOnSingleVm(ctx, logger, vmCopy.Name(), listOsServices); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func StartOsService(listOsServices []oservice.OsService) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "InstallOsService"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForStartOsService(ctx, logger, targets, listOsServices)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
