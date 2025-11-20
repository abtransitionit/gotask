package git

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	lgit "github.com/abtransitionit/golinux/mock/git"
)

// Description: git merge branch dev to main and git push to github (for a set of git repositories)
//
// Notes:
// - the host contains the git repository
func MergeDevToMain(hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - extract parameters
	// 11 - repo:list
	repoList := []string{}
	for _, v := range paramList[0] {
		repoList = append(repoList, fmt.Sprint(v)) // converts any -> string
	}
	// 12 - repo:folder
	if len(paramList) < 2 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("host: %s > repo folder not provided in paramList", hostName)
	}
	repoFolder := fmt.Sprint(paramList[1][0])

	// define var
	var failed []string
	results := make(map[string]bool)
	// const repoFolder = "/Users/max/wkspc/git" // TODO : externalize it to config file

	// loopt over item (git repo)
	for _, repoName := range repoList {

		// play CLI for each item - merge dev to main and push
		ok, err := lgit.MergeDevToMain(hostName, repoFolder, repoName, logger)

		// handle system error
		if err != nil {
			logger.Warnf("host: %s > repo %s > system error during git ops: %v", hostName, repoName, err)
			continue
		}

		// manage and collect logic errors
		results[repoName] = ok
		if !ok {
			failed = append(failed, repoName)
			logger.Debugf("host: %s > repo %s > git op failed", hostName, repoName)
		} else {
			logger.Debugf("host: %s > repo %s > update with success", hostName, repoName)
		}
	}

	// errors summary
	if len(failed) > 0 {
		return false, fmt.Errorf("host: %s > repo(s) failed: %v", hostName, failed)
	}

	// handle success
	return true, nil
}
