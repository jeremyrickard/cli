package v3action

import (
	"code.cloudfoundry.org/cli/util/manifest"
)

// ApplyApplicationManifest reads in the manifest from the path and provides it
// to the cloud controller.
func (actor Actor) ApplyApplicationManifest(pathToManifest string, spaceGUID string) (Warnings, error) {
	var allWarnings Warnings

	manifestApps, err := manifest.ReadAndMergeManifests(pathToManifest)
	if err != nil {
		return nil, err
	}

	appName := manifestApps[0].Name
	app, getAppWarnings, err := actor.GetApplicationByNameAndSpace(appName, spaceGUID)

	allWarnings = append(allWarnings, getAppWarnings...)
	if err != nil {
		return allWarnings, err
	}

	jobURL, applyManifestWarnings, err := actor.CloudControllerClient.CreateApplicationActionsApplyManifestByApplication(manifestApps[0], app.GUID)
	allWarnings = append(allWarnings, applyManifestWarnings...)
	if err != nil {
		return allWarnings, err
	}

	pollWarnings, err := actor.CloudControllerClient.PollJob(jobURL)
	allWarnings = append(allWarnings, pollWarnings...)
	return allWarnings, err
}
