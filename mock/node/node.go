package node

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lnode "github.com/abtransitionit/golinux/mock/node"
)

// Description: check if a set of node is SSH configured on a host.
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func CheckSshConf(hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// define var
	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of node that are not SSH configured

	// loop over each node
	for _, node := range nodeList {

		// play CLI - check if SSH is configured for the couple host/node
		ok, err := lnode.IsSshConfigured(hostName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > Node %s: > system error > checking SSH config: %v", hostName, node, err)
			continue
		}
		// failedNodes = append(failedNodes, node)

		// collect results
		results[node] = ok
		if !ok {
			failedNodes = append(failedNodes, node) // logical error: SSH simply not configured
			logger.Debugf("taget: %s > Node %s: > is not SSH configured", hostName, node)
		}

	}

	// If any node failed, return a single error message
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("host: %s > Node(s) that are not SSH configured: %v", hostName, failedNodes)
	}

	// handle success
	return true, nil
}

// Description: check if a set of node is SSH reachable from a host
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func CheckSshAccess(hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}

	// define var
	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of nodes that are not SSH reachable

	// loop over item (node)
	for _, node := range nodeList {

		// play CLI for each item - check if node is SSH reachable for the couple host/node
		ok, err := lnode.IsSshReachable(hostName, node, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > node %s: > system error > checking SSH access: %v", hostName, node, err)
			continue
		}

		// manage and collect logic errors
		results[node] = ok
		if !ok {
			failedNodes = append(failedNodes, node)
			// log
			logger.Infof("host: %s > node %s: > is not SSH reachable", hostName, node)
		}
	}

	// errors summary
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("host: %s > Node(s) that are not SSH reachable: %v", hostName, failedNodes)
	}

	// handle success
	return true, nil
}

// Description: Wait a set of node to be SSH reachable (within a delayMax) from a host.
//
// Notes:
// - a node is a remote VM, the localhost, a container or a remote container
// - a host is a node from which the ssh command is executed
func WaitIsSshOnline(hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	// 11 - node:List
	nodeList := []string{}
	for _, v := range paramList[0] {
		nodeList = append(nodeList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - repo:folder
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("host: %s > delay not provided in paramList", hostName)
	}
	delayMax := fmt.Sprint(paramList[1][0])

	// define var
	results := make(map[string]bool) // collector
	var failedNodes []string         // slice of nodes that are not SSH reachable (within a delayMax)

	// loop over item (node)
	for _, node := range nodeList {

		// play CLI for each item - check if node is SSH reachable for the couple host/node
		ok, err := lnode.IsSshOnline(hostName, node, delayMax, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > node %s: > system error > checking SSH access: %v", hostName, node, err)
			continue
		}

		// manage and collect logic errors
		results[node] = ok
		if !ok {
			failedNodes = append(failedNodes, node)
			// log
			logger.Infof("host: %s > node %s: > is not SSH reachable within the delay of %s", hostName, node, delayMax)
		}
	}

	// errors summary
	if len(failedNodes) > 0 {
		return false, fmt.Errorf("host: %s > Node(s) that are not SSH reachable: %v", hostName, failedNodes)
	}

	// handle success
	return true, nil
}
