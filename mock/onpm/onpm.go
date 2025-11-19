package onpm

import (
	"github.com/abtransitionit/gocore/logx"
	lonpm "github.com/abtransitionit/golinux/mock/onpm"
)

// Description: upgrade the linux OS of a set of nodes
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a target is a node from which the ssh command is executed
// - upgrade mean: set the OS native repositories and packages to the latest version
func UpgradeOs(targetName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	_, err := lonpm.UpgradeOs(targetName, logger)

	// handle system error
	if err != nil {
		logger.Warnf("target: %s > system error > checking SSH config: %v", targetName, err)
	}

	// handle success
	return true, nil
}
