package ovh

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/ovh"
)

// Description: renew the OVH token and store it locally on the host
//
// Notes:
// - the host is a node on which the client resides
func ListImageAvailable(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - define var
	ctx := context.Background()

	// log
	logger.Infof("list availbale OVH VPS OS images (api call)")

	// get the list of allOVH VPS OS images
	vpsList, err := ovh.ImageGetList(ctx, logger)
	if err != nil {
		return false, fmt.Errorf("getting available image:list : %v", err)
	}

	jsonx.PrettyPrintColor(vpsList)

	// handle success
	return true, nil
}
