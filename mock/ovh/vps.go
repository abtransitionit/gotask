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
	jsonResponse, err := ovh.VpsReinstall(ctx, logger, hostName)
	if err != nil {
		return false, fmt.Errorf("%s > re-installing VPS OS image >  %w", hostName, err)
	}

	// 3 - print the response of the request
	jsonx.PrettyPrintColor(jsonResponse)

	// handle success
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
