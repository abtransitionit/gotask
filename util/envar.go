package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/envar"
	"github.com/abtransitionit/golinux/filex"
)

func SetSingleEnvarOnSigleVm(ctx context.Context, logger logx.Logger, vmName string, customRcFileName string, envVVar envar.EnvVar) (string, error) {

	// persist this envar key value into the user's custom RC file
	userCustomRcFile := fmt.Sprintf(`$HOME/%s`, strings.TrimSpace(customRcFileName))
	line := fmt.Sprintf(`export %s=%s`, envVVar.Name, envVVar.Value)
	cli := filex.EnsureLineInFile(userCustomRcFile, line)
	_, err := run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// success
	logger.Debugf("%s: persist envar %s into file %s. check it", vmName, envVVar.Name, userCustomRcFile)
	return "", nil
}
func SetListEnvarOnSigleVm(ctx context.Context, logger logx.Logger, vmName string, customRcFileName string, listEnvVar envar.SliceEnvVar) (string, error) {
	// log
	logger.Debugf("%s: will persist following env var(s): %s", vmName, listEnvVar)

	// loop over each cli
	for _, envar := range listEnvVar {

		// persist the envar
		_, err := SetSingleEnvarOnSigleVm(ctx, logger, vmName, customRcFileName, envar)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func createSliceFuncEnvar(ctx context.Context, logger logx.Logger, targets []phase.Target, customRcFileName string, listEnvVar envar.SliceEnvVar) []syncx.Func {
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
			if _, err := SetListEnvarOnSigleVm(ctx, logger, vmCopy.Name(), customRcFileName, listEnvVar); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

func SetEnvar(customRcFileName string, listEnvVar envar.SliceEnvVar) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "GetPath"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncEnvar(ctx, logger, targets, customRcFileName, listEnvVar)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
