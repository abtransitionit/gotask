package ovh

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/ovh"
)

// Description: renew the OVH token and store it locally on the host
//
// Notes:
// - the host is a node on which the client resides
func RenewToken(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - define var
	logger.Infof("Renewing the OVH token")

	// 2 - renew the token
	_, err := ovh.RefreshToken(context.Background(), logger)
	if err != nil {
		logger.Errorf("%v", err)
	}

	// handle success
	return true, nil
}
