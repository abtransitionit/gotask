package file

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/file"
)

// Description: add an RC file on a node from hostname
//
// Notes:
//
// - task is executed by a goroutine on hostName
// - task concerns nodeName
func AddRcFile(hostName, nodeName, customRcFileName string, logger logx.Logger) error {
	// 1 - Create the File
	// 11 - get instance
	i := file.GetFile(customRcFileName, "~", "")

	// 12 - operate
	if err := i.ForceCreateRcFile(hostName, nodeName, logger); err != nil {
		return fmt.Errorf("creating rc file %s > %w", i.FullPath, err)
	}
	// log
	logger.Infof("%s:%s > created rc file : %s", hostName, nodeName, i.FullPath)

	// 2 - source this file in the Std RC File
	// 21 - get instance
	// 22 - operate
	// if err := i.AddStringOnce(hostName, nodeName, logger); err != nil {
	// 	return fmt.Errorf("adding string to file %s > %w", filePath, err)
	// }

	// handle success
	return nil
}

func RcAddPath(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// log
	logger.Infof("%s:%s > add envar $PATH to custom rc file ", hostName)
	// handle success
	return true, nil
}
