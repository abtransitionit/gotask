package selinux

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/selinux"
)

func Configure(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// get Instance
	selinux := selinux.GetSelinux()

	// operate
	if _, err := selinux.Configure(hostName, logger); err != nil {
		// handle error
		return false, fmt.Errorf("%s > configuring selinux > %v", hostName, err)
	}

	// handle success
	return true, nil

}
