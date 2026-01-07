package gopm

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/gopm"
)

func AddPkgGo(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// check
	if len(paramList) < 1 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s > list of GO clis or Destination folder not properly provided in paramList", hostName)
	}
	// 11 - List of go cli to install
	slice, err := filex.GetVarStructFromYaml[gopm.CliSlice](paramList[0])
	if err != nil {
		return false, fmt.Errorf("%s > getting cliName from paramList: %w", hostName, err)
	}
	// 12 - binary folder
	folderPath := fmt.Sprint(paramList[1][0])
	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range slice {
		// 31 - get instance
		i := gopm.GetCli(item)
		// 32 - operate
		if err := i.Install(hostName, folderPath, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("installing GO cli %s: %w", item.Name, err)
		}
	} // loop

	// 4 - manage error
	close(errChItem) // close the channel - signal that no more error will be sent
	// 41 - collect errors
	var errList []error
	for e := range errChItem {
		errList = append(errList, e)
	}

	// 42 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ %s > nb cli installation that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	return true, nil
}
