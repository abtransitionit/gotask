package ovh

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/ovh"
)

// Description: re-install the same OS image on a Linux host
//
// Notes:
// - the linux host is a remote VM on OVH cloud (aka. VPS)
func InstallVpsImage(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// 1 - define var
	ctx := context.Background()

	// log
	logger.Infof("Processing VPS: %s", hostName)

	// 2 - get the VPS id
	// jsonResponse, err := ovh.VpsReinstallHelper(ctx, logger, id)
	jsonResponse, err := ovh.VpsReinstallHelper(ctx, logger, hostName)
	if err != nil {
		// logger.Errorf("failed to re-install VPS: %v", err)
		return false, fmt.Errorf("❌ re-installing VPS OS image: %w", err)
	}

	// 3 - print the response of the request
	jsonx.PrettyPrintColor(jsonResponse)

	// play CLI
	// out, err := ovh.InstallVpsImage(ctx, hostName, logger)

	// handle system error
	// if err != nil {
	// 	logger.Warnf("%s > system error > updating OS image: %v", hostName, err)
	// }

	// handle success
	// logger.Debugf("%s > rebooting: %s", hostName, out)
	return true, nil
}

// Description: install another OS image on a VPS
func UpdateVpsImage(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {

	// play CLI
	_, err := ovh.UpdateVpsImage()

	// handle system error
	if err != nil {
		logger.Warnf("vps: %s > system error > updating OS image: %v", hostName, err)
	}

	// handle success
	return true, nil
}
