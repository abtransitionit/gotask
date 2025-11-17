package git

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lgit "github.com/abtransitionit/golinux/mock/git"
)

// Description: git merge branch dev to main and push to github (for a set of git repositories)
// func MergeDevToMain(targetName string, repoList []string, logger logx.Logger) (bool, error) {
func MergeDevToMain(targetName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	// 11 - repo:list
	repoList := []string{}
	for _, v := range paramList[0] {
		repoList = append(repoList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - repo:folder
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("target: %s > repo folder not provided in paramList", targetName)
	}
	repoFolder := fmt.Sprint(paramList[1][0])

	// define var
	var failed []string
	results := make(map[string]bool)
	// const repoFolder = "/Users/max/wkspc/git" // TODO : externalize it to config file

	// loopt over item (git repo)
	for _, repoName := range repoList {

		// play CLI for each item - merge dev to main and push
		ok, err := lgit.MergeDevToMain(targetName, repoFolder, repoName, logger)

		// handle system error
		if err != nil {
			logger.Warnf("target: %s > repo %s > system error during git ops: %v", targetName, repoName, err)
			continue
		}

		// manage and collect logic errors
		results[repoName] = ok
		if !ok {
			failed = append(failed, repoName)
			logger.Debugf("target: %s > repo %s > git op failed", targetName, repoName)
		} else {
			logger.Debugf("target: %s > repo %s > update with success", targetName, repoName)
		}
	}

	// errors summary
	if len(failed) > 0 {
		return false, fmt.Errorf("target: %s > repo(s) failed: %v", targetName, failed)
	}

	// handle success
	return true, nil
}
