package node

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: check if a node is SSH configured on a target (for a set of nodes).
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a target is a node from which the ssh command is executed
func CheckSshConf(targetName string, nodeList []string, logger logx.Logger) (bool, error) {

	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of node that are not SSH configured

	// loop over each node
	for _, node := range nodeList {

		// play CLI - check if SSH is configured for the couple target/node
		oko, err := lnode.IsSshConfigured(targetName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("target: %s > Node %s: > system error > checking SSH config: %v", targetName, node, err)
			continue
		}
		// failedNodes = append(failedNodes, node)

		// collect results
		results[node] = oko
		if !oko {
			failedNodes = append(failedNodes, node) // logical error: SSH simply not configured
			logger.Debugf("taget: %s > Node %s: > is not SSH configured", targetName, node)
		}

	}

	// If any node failed, return a single error message
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("target: %s > Node(s) that are not SSH configured: %v", targetName, failedNodes)
	}

	// handle success
	return true, nil
}

// Description: check if a node is SSH reachable from a target (for a set of nodes).
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a target is a node from which the ssh command is executed
func CheckSshAccess(targetName string, nodeList []string, logger logx.Logger) (bool, error) {

	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of nodes that are not SSH reachable

	// loop over item (node)
	for _, node := range nodeList {

		// play CLI for each item - check if node is SSH reachable for the couple target/node
		oko, err := lnode.IsSshReachable(targetName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("target: %s > node %s: > system error > checking SSH access: %v", targetName, node, err)
			continue
		}

		// manage and collect logic errors
		results[node] = oko
		if !oko {
			failedNodes = append(failedNodes, node)
			// log
			logger.Infof("target: %s > node %s: > is not SSH reachable", targetName, node)
		}
	}

	// errors summary
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("target: %s > Node(s) that are not SSH reachable: %v", targetName, failedNodes)
	}

	// handle success
	return true, nil
}
