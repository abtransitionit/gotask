package selinux

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/selinux"
)

func Configure(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// get Instance
	i := selinux.GetSelinux()

	// operate
	if _, err := i.Configure(hostName, logger); err != nil {
		return false, fmt.Errorf("%s > configuring selinux > %v", hostName, err)
	}

	// handle success
	return true, nil

}
