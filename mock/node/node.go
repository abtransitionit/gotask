package node

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: check if a set of targets are SSH configured.
//
// Notes:
// - a target is a remote VM, the localhost or a container
func CheckSshConf(targetList []string, logger logx.Logger) (bool, error) {

	results := make(map[string]bool) // collector
	var failedTargets []string       // slice of taget for which SSH is not configured

	// loop over each target
	for _, target := range targetList {
		// check if SSH is configured
		oko, err := lnode.IsSshConfigured(target, logger)
		// handle system error
		if err != nil {
			logger.Warnf("Target %s: > system error > checking SSH config: %v", target, err)
			failedTargets = append(failedTargets, target)
			continue
		}
		results[target] = oko
		logger.Infof("Target %s: > SSH configured > %v", target, oko)

		if !oko {
			failedTargets = append(failedTargets, target)
		}
	}

	// If any node failed, return a single error message
	if len(failedTargets) > 0 {
		return false, fmt.Errorf("SSH not configured for nodes: %v", failedTargets)
	}

	return true, nil
}

// Description: check if a set of nodes are SSH reachable.
func CheckSshAccess(nodes []string, logger logx.Logger) (bool, error) {
	results := make(map[string]bool)

	for _, node := range nodes {
		ok, err := lnode.IsSshReachable(node, logger)
		if err != nil {
			return false, fmt.Errorf("node %q: error checking SSH reachability: %w", node, err)
		}
		results[node] = ok
		logger.Infof("Node %q: SSH reachable = %v", node, ok)
	}
	// success
	return true, nil
}
