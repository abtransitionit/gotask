package oservice

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/oservice"
	"github.com/abtransitionit/golinux/property"
)

func InstallSingleOsServiceOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, osService oservice.OsService) (string, error) {

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	// get service canonical name
	osServiceCName, err := oservice.OsServiceReference.GetCName(osService)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// logic for installtion
	install := false
	switch osServiceCName {
	case "apparmor.service":
		if osFamily == "debian" {
			install = true
		}
	}

	// if nothing to install
	if !install {
		logger.Debugf("Skipping service installation for %s:%s:%s", vmName, osFamily, osServiceCName)
		return "", nil
	}

	// create file from string
	logger.Debugf("%s: create file %s from content", vmName, osService.Path)
	// get cli
	cli := filex.CreateFileFromStringAsSudo(osService.Path, osService.Content)
	// play it on remote
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	return "", nil
}

func InstallLisOsServiceOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsServices oservice.SliceOsService) (string, error) {
	// log
	logger.Debugf("%s: will install following OS service(s): %s", vmName, listOsServices.GetListName())

	// loop over each cli
	for _, osService := range listOsServices {

		// install the cli
		_, err := InstallSingleOsServiceOnSingleVm(ctx, logger, vmName, osService)
		if err != nil {
			return "", err
		}
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
