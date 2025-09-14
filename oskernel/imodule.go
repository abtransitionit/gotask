package oskernel

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/oskernel"
)

func loadListOsKModuleOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsKModule []string, kernelFilename string) (string, error) {

	logger.Debugf("%s: provision following OS kernel module : %s", vmName, listOsKModule)

	// load module for the current session
	// logger.Debugf("%s: loading kernel module(s) : %s for curent session", vmName, listOsKModule)
	cli, err := oskernel.LoadOsKModule(listOsKModule)
	if err != nil {
		return "", err
	}
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", err
	}
	logger.Debugf("%s: 🅐 loaded kernel module(s) : %s for curent session", vmName, listOsKModule)

	// save the content to the file
	// logger.Debugf("%s: persisting kernel module(s) to be applyed after rebbot in file : %s", vmName, filePath)
	// - create content from slice
	stringContent := list.GetStringWithSepFromSlice(listOsKModule, "\n")
	// - define the kernel file path to write this content
	filePath := oskernel.GetKModuleFilePath(kernelFilename)
	// - define the cli
	cli = filex.CreateFileFromStringAsSudo(filePath, stringContent)
	// - play the cli
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// success
	logger.Debugf("%s: 🅑 persisted kernel module(s) in file : %s", vmName, filePath)
	return "", nil
}

// Name: createSliceFunc
//
// Description: create the slice of tasks
//
// Parameters:
// - l: logger
// - targets: list of targets
//
// Returns:
//
// - slice of syncx.Func
//
// Notes:
//
// - as many tasks as there are VMs
// - Only VM targets are included; others are skipped with a warning.
func createSliceFuncForKModule(ctx context.Context, logger logx.Logger, targets []phase.Target, listOsKModule []string, kernelFilename string) []syncx.Func {
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
		tasks = append(tasks, func() error {
			if _, err := loadListOsKModuleOnSingleVm(ctx, logger, vmCopy.Name(), listOsKModule, kernelFilename); err != nil {
				logger.Errorf("🅣 Failed to install Dnfapt repository on VM %s: %v", vmCopy.Name(), err)
				return err
			}

			// logger.Infof("🅣 VM %s package installed successfully", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

// name: UpgradeVmOs
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func LoadOsKModule(listOsKModule []string, kernelFilename string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: LoadOsKModule")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : LoadOsKModule")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForKModule(ctx, logger, targets, listOsKModule, kernelFilename)

		// Log number of tasks
		logger.Infof("🅣 Phase LoadOsKModule has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase LoadOsKModule on %d VM(s)", len(tasks)), nil
	}
}
