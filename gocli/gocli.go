// file gotask/gocli/icli.go
package gocli

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/gocli"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/gocore/url"
	"github.com/abtransitionit/golinux/property"
)

// Name: Install
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
func InstallOnSingleVm(logger logx.Logger, vmName string, listGoClis []gocli.GoCli) (string, error) {

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

	// log
	logger.Debugf("%s: will install following GO CLI(s): %s", vmName, listGoClis)

	// loop over each cli
	for _, goCli := range listGoClis {

		// get the URL of the CLI to install that is also VM specific
		urlResolved, err := gocli.ResolveURL(logger, goCli, osType, osArch, uname)
		if err != nil {
			return "", err
		}

		// download file pointed by URL - on local host
		localPath, err := url.Download(goCli.Name, urlResolved)
		if err != nil {
			return "", err
		}

		// detect the type of the downloaded file
		fileType, err := filex.DetectBinaryType(localPath) // eg. tar.gz, zip, binary
		if err != nil {
			fmt.Println("Detection failed:", err)
			return "", err
		}

		// move file when possible
		switch fileType {
		case "zip":
			logger.Debugf("🌐 Cli: % not yet managed", goCli.Name)
		case "tgz":
			logger.Debugf("🌐 Cli: %s:type:tgz:%s - need more works", goCli.Name, localPath)
		case "exe":
			logger.Debugf("🌐 Cli: %s:type:Exe:%s : now sudo copy %s to folder /usr/local/bin with name xxx", goCli.Name, localPath)
			// lookup OsName for a goCli that end up with an Exe type

		default:
			return "", fmt.Errorf("Unsupported file type %s", fileType)
		}
	}

	// logger.Infof("Cli: %s Url: %s", cli.Name, url)
	return "", nil

	// // loop over cli
	// for _, goCli := range listGoClis {
	// 	logger.Debugf("%s: installing GO cli: %s ", vmName, goCli.Name)
	// 	_, err := gocli.Install(logger, goCli, osType, osArch, uname)
	// 	if err != nil {
	// 		logger.Debugf("%s: error installing GO cli: %s ", vmName, goCli.Name)
	// 		return "", err
	// 	}
	// 	logger.Debugf("%s: successfully installed GO cli: %s ", vmName, goCli.Name)
	// }
	// return "", nil
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
func createSliceFuncForInstall(logger logx.Logger, targets []phase.Target, listGoClis []gocli.GoCli) []syncx.Func {
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
			if _, err := InstallOnSingleVm(logger, vmCopy.Name(), listGoClis); err != nil {
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
		tasks := createSliceFuncForInstall(logger, targets, listGoClis)

		// Log number of tasks
		logger.Infof("🅣 Phase %s has %d concurent tasks", appx, len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		// return fmt.Sprintf("🅣 Terminated phase InstallOnVm on %d VM(s)", len(tasks)), nil
		return "", nil
	}
}
