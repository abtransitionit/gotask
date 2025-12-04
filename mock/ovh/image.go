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
	// define var
	// ctx := context.Background()

	// 1 - get parameters

	// log
	logger.Infof("list availbale OVH VPS OS images (api call)")

	// // get the list of all OVH VPS OS images
	// vpsList, err := ovh.VpsImageGetList(ctx, "o5d", logger)
	// if err != nil {
	// 	return false, fmt.Errorf("getting available image:list : %v", err)
	// }

	// jsonx.PrettyPrintColor(vpsList)

	// handle success
	return true, nil
}

func GetImageInfo(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - define var
	ctx := context.Background()

	// 2 -
	info, err := ovh.GetVpsImageId2(ctx, "vps-a7a8f7f6.vps.ovh.net", logger)
	if err != nil {
		return false, fmt.Errorf("getting info : %v", err)
	}

	jsonx.PrettyPrintColor(info)

	// handle success
	return true, nil
}
