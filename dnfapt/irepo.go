// File: gotask/dnfapt/upgrade.go
package dnfapt

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
	"github.com/abtransitionit/golinux/dnfapt"
	"github.com/abtransitionit/golinux/filex"
	"github.com/abtransitionit/golinux/property"
)

// Name: upgradeSingleVm
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
func installSingleDaRepoOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, daRepo dnfapt.DaRepo) (string, error) {
	// logger.Debugf("will install dnfapt package repository: %v\n", daRepo.Name)
	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	// get property
	osDistro, err := property.GetProperty(vmName, "osdistro")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	logger.Debugf("%s:%s:%s Installing dnfapt package repository: %s", vmName, osFamily, osDistro, daRepo.Name)

	// get the repo file path
	repoFilePath, err := dnfapt.GetRepoFilePath(osFamily, daRepo)
	if err != nil {
		return "", err
	}
	// get the URLs resolve from template
	repoFileContent, err := dnfapt.GetRepoFileContent(osFamily, daRepo)
	if err != nil {
		return "", err
	}

	// save the repo file
	logger.Debugf("🅰️ %s:%s:%s creating file repoFilePath: %s", vmName, osFamily, osDistro, repoFilePath)
	cli := filex.CreateFileFromStringAsSudo(repoFilePath, repoFileContent)
	_, err = run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}

	// see the repo file
	cli = fmt.Sprintf("ls -ial %s && echo && cat %s", repoFilePath, repoFilePath)
	output, err := run.RunCliSsh(vmName, cli)
	if err != nil {
		return "", fmt.Errorf("failed to play cli %s on vm '%s': %w", cli, vmName, err)
	}
	fmt.Println(output)
	// log

	// // install the dnfapt package repository
	// _, err = dnfapt.InstallDaRepository(ctx, logger, osFamily, daRepo)
	// if err != nil {
	// 	return "", err
	// }
	// success
	logger.Debugf("%s:%s:%s Installed dnfapt package repository: %s", vmName, osFamily, osDistro, daRepo.Name)
	return "", nil
}
func installListDaRepoOnSingleVm(ctx context.Context, logger logx.Logger, vmName string, listDaRepo dnfapt.SliceDaRepo) (string, error) {
	// log
	logger.Debugf("%s: will install following dnfapt package repositories: %s", vmName, listDaRepo.GetListName())

	// loop over each cli
	for _, daRepoName := range listDaRepo {

		// install the dnfapt package repository
		_, err := installSingleDaRepoOnSingleVm(ctx, logger, vmName, daRepoName)
		if err != nil {
			return "", err
		}
	}
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
func createSliceFuncForInstallRepo(ctx context.Context, logger logx.Logger, targets []phase.Target, listDaRepo dnfapt.SliceDaRepo) []syncx.Func {
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
			if _, err := installListDaRepoOnSingleVm(ctx, logger, vmCopy.Name(), listDaRepo); err != nil {
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
func InstallDaRepository(listDaRepo dnfapt.SliceDaRepo) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: UpdateVmOsApp")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : UpdateVmOsApp")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForInstallRepo(ctx, logger, targets, listDaRepo)

		// Log number of tasks
		logger.Infof("🅣 Phase UpdateVmOsApp has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase UpdateVmOsApp on %d VM(s)", len(tasks)), nil
	}
}
