package util

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/property"
)

func CreateCustomRcFileOnSigleVm(ctx context.Context, logger logx.Logger, vmName string, fileName string) (string, error) {

	// get property
	stdRcFilePath, err := property.GetProperty(vmName, "rcfilepath")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// create custom RC file
	cli := filex.TouchFile(fmt.Sprintf(`$HOME/%s`, fileName))
	customRcFilePath, err := run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// add line to std rc file
	cli = filex.EnsureLineInFile(stdRcFilePath, fmt.Sprintf(`. %s`, customRcFilePath))
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// success
	logger.Debugf("%s: std rc file %s source custom rc file %s", vmName, stdRcFilePath, customRcFilePath)
	return "", nil
}

func createSliceFuncRcFile(ctx context.Context, logger logx.Logger, targets []phase.Target, fileName string) []syncx.Func {
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
			if _, err := CreateCustomRcFileOnSigleVm(ctx, logger, vmCopy.Name(), fileName); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func CreateCustomRcFile(fileName string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "GetPath"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncRcFile(ctx, logger, targets, fileName)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
