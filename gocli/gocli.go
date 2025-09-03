// file gotask/gocli/icli.go
package gocli

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/gocli"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/property"
)

// Name: InstallSingleGoCliOnSingleVm
//
// Description: the single task: update the OS with standard/required/missing dnfapt packages
//
// Parameters:
//
// - vmName: name of the VM
//
// Returns:
// - nil if the VM is reachable,
// - an error if the VM is not configured, not reachable or if there was an SSH failure.
//
// Notes:
// - pure logic : no logging
func InstallSingleGoCliOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, goCli gocli.GoCli) (string, error) {
	// get vm property
	osType, err := property.GetProperty(vmName, "ostype")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	osArch, err := property.GetProperty(vmName, "osarch")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	uname, err := property.GetProperty(vmName, "uname")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// get the URL of the CLI to install that is also VM specific
	urlResolved, err := gocli.ResolveURL(logger, goCli, osType, osArch, uname)
	if err != nil {
		return "", err
	}

	// download file pointed by URL
	cmd := fmt.Sprintf("goluc do download %s -p %s", urlResolved, goCli.Name)
	filePath, err := run.RunCliSsh(vmName, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to play cli on vm: '%s': '%s' : %w", vmName, cmd, err)
	}

	// detect the type of the downloaded file
	cmd = fmt.Sprintf("goluc do detect %s ", filePath)
	fileType, err := run.RunCliSsh(vmName, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to play cli on vm: '%s': '%s' : %w", vmName, cmd, err)
	}

	// copy artifact to destination
	switch fileType {
	case "zip":
		logger.Debugf("🌐 ZIp:%s not yet managed", filePath)
		return "", fmt.Errorf("filetype: Zip not yet managed")

	case "tgz":
		// get cli
		cli := filex.CpTgzFile(filePath, goCli.Name)
		if err != nil {
			return "", fmt.Errorf("failed to get code from library : %w", err)
		}

		// play it on remote
		_, err := run.RunCliSsh(vmName, cli)
		if err != nil {
			return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
		}
		logger.Debugf("🌐🅣 Tgz:%s untarred to /usr/local/bin/%s ", goCli.Name, goCli.Name)
	case "exe":

		// get cli
		dstFile := "/usr/local/bin/" + goCli.Name
		cli, err := filex.CpAsSudo(ctx, logger, filePath, dstFile)
		if err != nil {
			return "", fmt.Errorf("failed to get code from library : %w", err)
		}

		// play it on remote
		_, err = run.RunCliSsh(vmName, cli)
		if err != nil {
			return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
		}

		//success
		logger.Debugf("🌐🅔 Exe:%s:%s copied to ", goCli.Name, filePath, dstFile)
	default:
		return "", fmt.Errorf("unsupported file type %s", fileType)
	}

	return "", nil
}

// Name: InstallLisGoCliOnSingleVm
//
// Description: the single task: update the OS with standard/required/missing dnfapt packages
//
// Parameters:
//
// - vmName: name of the VM
//
// Returns:
// - nil if the VM is reachable,
// - an error if the VM is not configured, not reachable or if there was an SSH failure.
//
// Notes:
// - pure logic : no logging
func InstallLisGoCliOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listGoClis []gocli.GoCli) (string, error) {
	// log
	logger.Debugf("%s: will install following GO CLI(s): %s", vmName, listGoClis)

	// loop over each cli
	for _, goCli := range listGoClis {

		// install the cli
		_, err := InstallSingleGoCliOnSingleVm(ctx, logger, vmName, goCli)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// Name: createSliceFuncForInstall
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
func createSliceFuncForInstall(ctx context.Context, logger logx.Logger, targets []phase.Target, listGoClis []gocli.GoCli) []syncx.Func {
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
			if _, err := InstallLisGoCliOnSingleVm(ctx, logger, vmCopy.Name(), listGoClis); err != nil {
				logger.Errorf("🅣 Failed to execute task on VM %s: %v", vmCopy.Name(), err)
				return err
			}
			logger.Infof("🅣 task on VM %s succeded", vmCopy.Name())
			return nil
		})
	}

	return tasks
}

// name: InstallOnVm
//
// description: the overall task.
//
// Notes:
// - Each target must implement the Target interface.
func InstallOnVm(listGoClis []gocli.GoCli) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {
		appx := "InstallGoCliOnVm"
		logger.Infof("🅣 Starting phase: %s", appx)
		// check paramaters
		if len(targets) == 0 {
			logger.Warnf("🅣 No targets provided to phase: %s", appx)
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForInstall(ctx, logger, targets, listGoClis)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return "", nil
	}
}
