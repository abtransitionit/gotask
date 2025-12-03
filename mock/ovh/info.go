package ovh

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/ovh"
)

// Description: get global informations on OVH API
func ListInfo(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	ovh.ListInfo()
	return true, nil
}
