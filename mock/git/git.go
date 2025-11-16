package git

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/git"
)

// Description: git merge branch dev to main and push to github (for a set of git repositories)
func MergeDevToMain(targetName string, repoList []string, logger logx.Logger) (bool, error) {

	var failed []string
	results := make(map[string]bool)

	for _, repo := range repoList {

		ok, err := git.MergeDevToMain(targetName, repo, logger)

		// system error → log and continue
		if err != nil {
			logger.Warnf("target: %s > repo %s > system error during git ops: %v", targetName, repo, err)
			continue
		}

		// collect
		results[repo] = ok
		if !ok {
			failed = append(failed, repo)
			logger.Debugf("target: %s > repo %s > git op failed", targetName, repo)
		} else {
			logger.Debugf("target: %s > repo %s > update with success", targetName, repo)
		}
	}

	if len(failed) > 0 {
		return false, fmt.Errorf("target: %s > repo(s) failed: %v", targetName, failed)
	}

	return true, nil
}
