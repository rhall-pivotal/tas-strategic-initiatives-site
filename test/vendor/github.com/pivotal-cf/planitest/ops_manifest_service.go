package planitest

import (
	"fmt"
)

type OpsManifestService struct {
	cmdRunner CommandRunner
}

func NewOpsManifestService() (*OpsManifestService, error) {
	return NewOpsManifestServiceWithRunner(NewExecutor())
}

func NewOpsManifestServiceWithRunner(cmdRunner CommandRunner) (*OpsManifestService, error) {
	return &OpsManifestService{cmdRunner: cmdRunner}, nil
}

func (o OpsManifestService) RenderManifest(config ProductConfig) (Manifest, error) {
	response, errOutput, err := o.cmdRunner.Run(
		"ops-manifest",
		"-m", config.MetadataFile,
		"-c", config.ConfigFile,
	)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve bosh manifest: %s: %s", err, errOutput)
	}
	return NewManifest(string(response), o.cmdRunner), nil
}
