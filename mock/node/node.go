package node

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: check if a set of nodes are SSH configured.
// CheckSshConf checks SSH configuration for a list of nodes.
func CheckSshConf(nodes []string, logger logx.Logger) (bool, error) {
	results := make(map[string]bool)

	for _, node := range nodes {
		ok, err := lnode.IsSshConfigured(node, logger)
		if err != nil {
			return false, fmt.Errorf("checking SSH config for node %q: %w", node, err)
		}
		results[node] = ok
		logger.Infof("Node %q: SSH configured = %v", node, ok)
	}

	// check if any node failed
	for node, ok := range results {
		if !ok {
			return false, fmt.Errorf("SSH not configured for node %q", node)
		}
	}

	// success
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
