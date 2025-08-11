package dnfapt

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/dnfapt"
)

// Name: CreateRepo
// Description: Create a Dnfapt repository
func CreateRepo(repoName string) {
	logx.Info("[%s] Installalling package", repoName)

	if err := dnfapt.Install(repoName); err != nil {
		logx.Error("creation of repo failed: %s", err)
		return
	}

	logx.Info("[%s] successfully created dnfapt repository", repoName)
}

// // Inside a function in your gotask library

// hosts := []string{"o1u", "o2r"} // Provided by the user

// for _, host := range hosts {
// 	// Step 1: Check if the VM is configured before trying to connect
// 	if configured, err := executor.IsSSHConfigured(host); err != nil {
// 		logx.Error("Failed to check SSH config for %s: %s", host, err)
// 		continue
// 	} else if !configured {
// 		logx.Error("VM '%s' is not configured in your SSH config. Skipping.", host)
// 		continue
// 	}

// 	// Step 2: Now that we know it's configured, it's safe to run the command
// 	output, err := executor.RunSSH(host, "systemctl status sshd")
// 	if err != nil {
// 		logx.Error("Error running command on %s: %s", host, err)
// 		continue
// 	}
// 	logx.Info("Output from %s: %s", host, output)
// }
