// File: gotask/dnfapt/upgrade.go
package dnfapt

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
	"github.com/abtransitionit/gocore/run"
	"github.com/abtransitionit/gocore/syncx"
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
func updateSingleVmOsApp(vmName string) error {
	var cmds []string

	// get property
	osFamily, err := property.GetProperty(vmName, "osfamily")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	switch osFamily {
	case "debian":
		cmds = []string{
			"DEBIAN_FRONTEND=noninteractive sudo apt-get -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' update -qq -y",
			"DEBIAN_FRONTEND=noninteractive sudo apt-get -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' upgrade -qq -y",
			"DEBIAN_FRONTEND=noninteractive sudo apt-get -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' clean -qq",
		}
	case "rhel", "fedora":
		cmds = []string{
			"sudo dnf update -q -y",
			"sudo dnf upgrade -q -y",
			"sudo dnf clean all",
		}
	default:
		return fmt.Errorf("unsupported Linux OS Family: %s", osFamily)
	}

	// Join commands with && to run them sequentially
	cli := strings.Join(cmds, " && ")

	// play the CLI
	_, err = run.RunOnVm(vmName, cli)
	if err != nil {
		// fmt.Printf("Command output from %s:\n%s\n", vmName, output)
		return fmt.Errorf("failed to upgrade OS on VM %s: %w", vmName, err)
	}
	// fmt.Printf("Successfully upgraded OS on %s. Output:\n%s\n", vmName, output)
	return nil
}

// func upgradeSingleVmOs(vm *phase.Vm) error {
// 	reachable, err := executor.IsVmSshReachable(vm.Name())
// 	if err != nil {
// 		return errorx.NewWithNoStack("🅣 SSH check failed for VM %s", vm.Name())
// 	}
// 	if !reachable {
// 		return errorx.NewWithNoStack("🅣 VM %s is not reachable via SSH", vm.Name())
// 	}
// 	return nil
// }

// NAme: createSliceFunc
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
func createSliceFuncForUpdate(l logx.Logger, targets []phase.Target) []syncx.Func {
	var tasks []syncx.Func

	for _, t := range targets {
		if t.Type() != "Vm" {
			continue
		}

		vm, ok := t.(*phase.Vm)
		if !ok {
			l.Warnf("🅣 Target %s is not a VM, skipping", t.Name())
			continue
		}

		vmCopy := vm // capture for closure
		tasks = append(tasks, func() error {
			if err := upgradeSingleVmOs(vmCopy.Name()); err != nil {
				l.Errorf("🅣 Failed to upgrade VM %s: %v", vmCopy.Name(), err)
				return err
			}

			l.Infof("🅣 VM %s upgraded successfully", vmCopy.Name())
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
func UpdateVmOsApp(listRequiredPackage []string) phase.PhaseFunc {
	return func(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

		logger.Info("🅣 Starting phase: UpgradeVmOs")
		// check paramaters
		if len(targets) == 0 {
			logger.Warn("🅣 No targets provided to : UpgradeVmOs")
			return "", nil
		}

		// Build slice of functions
		tasks := createSliceFuncForUpgrade(logger, targets)

		// Log number of tasks
		logger.Infof("🅣 Phase UpgradeVmOs has %d concurent tasks", len(tasks))

		// Run tasks in the slice concurrently
		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			return "", errs[0] // return first error encountered
		}

		return fmt.Sprintf("🅣 Terminated phase UpgradeVmOs on %d VM(s)", len(tasks)), nil
	}
}
