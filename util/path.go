package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/property"
	"github.com/abtransitionit/golinux/stringx"
)

func SetPathTreeOnSigleVm(ctx context.Context, logger logx.Logger, vmName string, basePath string, customRcFileName string) (string, error) {

	// get property
	pathTree, err := property.GetProperty(vmName, "pathtree", basePath)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	// currentUser, err := property.GetProperty(vmName, "osuser")
	// if err != nil {
	// 	return "", fmt.Errorf("%v", err)
	// }
	envarPATH, err := property.GetProperty(vmName, "envar", "PATH")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// concatenate path
	cli := stringx.EnsureFusionStringUniq(envarPATH, pathTree, ":")
	fusionPATH, err := run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// persist this PATH to the user's custom RC file
	userCustomRcFile := fmt.Sprintf(`$HOME/%s`, strings.TrimSpace(customRcFileName))
	line := fmt.Sprintf(`export PATH=%s`, fusionPATH)
	cli = filex.EnsureLineInFile(userCustomRcFile, line)
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// // create tmp file from path
	// cli = filex.CreateTmpFileFromString(vmName, pathTree)
	// filename, err := run.RunCliSsh(vmName, cli)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	// }

	// success
	logger.Debugf("%s: persist PATH into file %s. check it", vmName, userCustomRcFile)
	return "", nil
}

func createSliceFuncPathTree(ctx context.Context, logger logx.Logger, targets []phase.Target, basePath string, customRcFileName string) []syncx.Func {
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
			if _, err := SetPathTreeOnSigleVm(ctx, logger, vmCopy.Name(), basePath, customRcFileName); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func SetPath(basePath string, customRcFileName string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "GetPath"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncPathTree(ctx, logger, targets, basePath, customRcFileName)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
