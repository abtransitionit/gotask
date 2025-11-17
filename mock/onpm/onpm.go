package onpm

import (
	"fmt"

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

	// 1 - extract parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// define var
	// results := make(map[string]bool) // collector
	var failedNodes []string // slice of node that are not SSH configured

	// loop over each node
	for _, node := range nodeList {

		// play CLI - check if SSH is configured for the couple target/node
		_, err := lonpm.UpgradeOs(logger)

		// handle system error
		if err != nil {
			logger.Warnf("target: %s > Node %s: > system error > checking SSH config: %v", targetName, node, err)
			continue
		}
		// failedNodes = append(failedNodes, node)

		// // collect results
		// results[node] = oko
		// if !oko {
		// 	failedNodes = append(failedNodes, node) // logical error: SSH simply not configured
		// 	logger.Debugf("taget: %s > Node %s: > is not SSH configured", targetName, node)
		// }

	}

	// If any node failed, return a single error message
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("target: %s > Node(s) that are not SSH configured: %v", targetName, failedNodes)
	}

	// handle success
	return true, nil
}
