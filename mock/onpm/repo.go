package onpm

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	lonpm "github.com/abtransitionit/golinux/mock/onpm"
	"gopkg.in/yaml.v3"
)

// Description: add native os package repositories to a Linux host
func AddRepo(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - get parameters
	// 11 - repository:list
	if len(paramList) < 1 || len(paramList[0]) == 0 {
		return false, fmt.Errorf("%s > list of repo not provided in paramList", hostName)
	}
	raw := paramList[0]

	b, err := yaml.Marshal(raw)
	if err != nil {
		return false, err
	}
	repoList, err := filex.GetVarStruct[lonpm.RepoSlice](fmt.Sprint(string(b)))
	if err != nil {
		logger.Errorf("%v", err)
	}

	// 2 - manage error reporting
	nbItem := len(repoList)
	errChItem := make(chan error, nbItem) // define a channel to collect errors from each goroutine

	// 3 - loop over item
	for _, item := range repoList {
		_, grErr := lonpm.AddRepo(hostName, item, logger) // the code to be executed
		if grErr != nil {
			logger.Errorf("(%s) %s/%s > %v", phaseName, hostName, item.Name, grErr) // send goroutines error if any into the chanel
			// send error if any into the chanel
			errChItem <- fmt.Errorf("%w", grErr)
		}
	} // loop
	close(errChItem) // close the channel - signal that no more error will be sent

	// 4 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 5 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ %s > nb repo that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// 6 - handle success
	return true, nil
}
