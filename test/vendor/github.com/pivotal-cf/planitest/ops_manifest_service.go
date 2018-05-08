package planitest

import (
	"fmt"
)

type OpsManifestService struct {
	config    OpsManifestConfig
	cmdRunner CommandRunner
}

type OpsManifestConfig struct {
	Name         string
	Version      string
	ConfigFile   string
	MetadataFile string
}

func NewOpsManifestService(config OpsManifestConfig) (*OpsManifestService, error) {
	return NewOpsManifestServiceWithRunner(config, NewExecutor())
}

func NewOpsManifestServiceWithRunner(config OpsManifestConfig, cmdRunner CommandRunner) (*OpsManifestService, error) {
	return &OpsManifestService{config: config, cmdRunner: cmdRunner}, nil
}

func (o *OpsManifestService) Configure(additionalProperties map[string]interface{}) error {
	return nil
}

func (o OpsManifestService) RenderManifest(additionalProperties map[string]interface{}) (Manifest, error) {
	response, errOutput, err := o.cmdRunner.Run(
		"ops-manifest",
		"-m", o.config.MetadataFile,
		"-c", o.config.ConfigFile,
	)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve bosh manifest: %s: %s", err, errOutput)
	}
	return NewManifest(string(response)), nil
}
