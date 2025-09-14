package oskernel

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/oskernel"
)

func SetListOsKParamOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsKParam oskernel.SliceOsKParam, kernelFilename string) (string, error) {
	// log
	logger.Debugf("%s: will provision following OS kernel parameter(s) %v", vmName, listOsKParam)

	// load parameters for the current session
	// logger.Debugf("%s: loading kernel parameters(s) : %s for curent session", vmName)
	cli, err := oskernel.LoadOsKParam()
	if err != nil {
		return "", err
	}
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", err
	}
	logger.Debugf("%s: 🅐 loaded kernel parameter(s) : for curent session", vmName)

	// save the content to the file
	// logger.Debugf("%s: persisting kernel module(s) to be applyed after rebbot in file : %s", vmName, filePath)
	// - create content from slice
	stringContent := listOsKParam.GetContent()
	// - define the kernel file path to write this content
	filePath := oskernel.GetKParamFilePath(kernelFilename)
	// - define the cli
	cli = filex.CreateFileFromStringAsSudo(filePath, stringContent)
	// - play the cli
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// success
	logger.Debugf("%s: 🅑 persisted kernel parameter(s) in file : %s", vmName, filePath)
	return "", nil

}

// func SetListOsKParamOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listOsKParam []string) (string, error) {
// 	// log
// 	logger.Debugf("%s: will set following OS kernel parameter(s) : %s", vmName, listOsKParam)

// 	// loop over each cli
// 	for _, osCoreParam := range listOsKParam {

// 		// Get the CLI to activate one OS core module
// 		cli, err := oscore.SetOsCoreParameter(osCoreParam)
// 		if err != nil {
// 			return "", err
// 		}

// 		// // play the CLI
// 		logger.Debugf("%s:%s setiing OS kernel parameter: %s,%s", vmName, osCoreParam, cli)
// 		// _, err = run.RunOnVm(vmName, cli)
// 		// if err != nil {
// 		// 	return "", err
// 		// }

// 	}
// 	return "", nil
// }

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
func createSliceFuncForKParam(ctx context.Context, logger logx.Logger, targets []phase.Target, listOsKParam oskernel.SliceOsKParam, kernelFilename string) []syncx.Func {
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
			if _, err := SetListOsKParamOnSingleVm(ctx, logger, vmCopy.Name(), listOsKParam, kernelFilename); err != nil {
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
func LoadOsKParam(listOsKParam oskernel.SliceOsKParam, kernelFilename string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: LoadOsKParam")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : LoadOsKParam")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForKParam(ctx, logger, targets, listOsKParam, kernelFilename)

		// Log number of tasks
		logger.Infof("🅣 Phase LoadOsKParam has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase LoadOsKParam on %d VM(s)", len(tasks)), nil
	}
}
