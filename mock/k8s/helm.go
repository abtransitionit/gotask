package k8s

import (
	"errors"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
	"github.com/abtransitionit/golinux/mock/k8scli/helm"
)

// Description: Add helm repo to a Helm client
//
// Notes:
//
// - it just adds the repo to the Helm client cfg files, it does not install any chart
func AddRepoHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 10 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s:helm > list repos or helm client not properly provided in paramList", hostName)
	}
	// 11 - name of helm client host
	helmClientNodeName := fmt.Sprint(paramList[1][0])

	// 12 - List of helm repos
	slice, err := filex.GetVarStructFromYaml[helm.RepoSlice](paramList[0])
	if err != nil {
		return false, fmt.Errorf("%s:helm > getting list from paramList: %w", hostName, err)
	}

	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range slice {
		// 31 - get instance and operate
		// i := helm.GetRepo(item.Name, "")
		// if err := i.Add(hostName, helmClientNodeName, logger); err != nil {
		i := helm.Resource{Type: helm.ResRepo, Name: item.Name}
		if _, err := i.Add("local", helmClientNodeName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("adding Helm repo %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb item that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	// logger.Infof("%s:%s > added repo %v on helm client", hostName, helmClientNodeName, slice)
	return true, nil
}

// Description: Add helm charts into a K8s cluster from a Helm client
//
// Notes:
//
// - the helm repo of the chart must be presents in the Helm client
func AddChartHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	logger.Info("AddChartHelm : adding chart: TODO")
	return true, nil
}

// func InstallReleaseHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
// 	// 1 - get parameters
// 	// 10 - check
// 	if len(paramList) < 1 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
// 		return false, fmt.Errorf("%s:helm > list repos or helm client not properly provided in paramList", hostName)
// 	}
// 	// log
// 	logger.Info("AddReleaseHelm : adding release: TODO")
// 	// 1 - get parameters
// 	// 11 - List helm release

// 	// handle success
// 	return true, nil
// }

// Description: Add a helm release from a Helm Chart into a K8s cluster's namespaces
//
// Notes:
//
// - the helm repo of the chart must be presents in the Helm client
func InstallReleaseHelm(phaseName, hostName string, paramList [][]any, logger logx.Logger) (bool, error) {
	// 1 - get parameters
	// 10 - check
	if len(paramList) < 1 || len(paramList[0]) == 0 || len(paramList[1]) == 0 {
		return false, fmt.Errorf("%s:helm > list releases or helm client not properly provided in paramList", hostName)
	}
	// 11 - name of helm client host
	helmClientNodeName := fmt.Sprint(paramList[1][0])

	// 12 - List of helm release
	slice, err := filex.GetVarStructFromYaml[helm.ReleaseSlice](paramList[0])
	if err != nil {
		return false, fmt.Errorf("%s:helm > getting list from paramList: %w", hostName, err)
	}

	// 2 - manage error reporting
	nbItem := len(slice)
	errChItem := make(chan error, nbItem) // define a channel to collect errors

	// 3 - loop over item
	for _, item := range slice {
		// 31 - get instance
		logger.Debugf("item.Param >  %v", item.Param)
		i := helm.GetRelease(item.Name, item.Chart.QName, item.Chart.Version, item.Namespace, item.Param)
		// 32 - operate
		if err := i.Install(hostName, helmClientNodeName, logger); err != nil {
			// send error if any into the chanel
			errChItem <- fmt.Errorf("installing Helm release %s: %w", item.Name, err)
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
		logger.Errorf("❌ %s > nb item that failed: %d", hostName, nbGroutineFailed)
		return false, errCombined
	}

	// handle success
	// logger.Infof("%s:%s > added repo %v on helm client", hostName, helmClientNodeName, slice)
	return true, nil
}
