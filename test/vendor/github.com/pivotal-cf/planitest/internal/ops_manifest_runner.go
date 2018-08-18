package internal

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type OpsManifestRunner struct {
	cmdRunner CommandRunner
	FileIO    FileIO
}

func NewOpsManifestRunner(cmdRunner CommandRunner, fileIO FileIO) OpsManifestRunner {
	return OpsManifestRunner{
		cmdRunner: cmdRunner,
		FileIO:    fileIO,
	}
}

func (o OpsManifestRunner) GetManifest(productProperties, metadataFilePath string) (map[string]interface{}, error) {
	configFile, err := o.FileIO.TempFile("", "")
	configFileYML := fmt.Sprintf("%s.yml", configFile.Name())
	os.Rename(configFile.Name(), configFileYML)

	if err != nil {
		return nil, err //not tested
	}

	_, err = configFile.WriteString(productProperties)
	if err != nil {
		return nil, err //not tested
	}

	response, errOutput, err := o.cmdRunner.Run(
		"ops-manifest",
		"--config-file", configFileYML,
		"--metadata-path", metadataFilePath,
	)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve manifest: %s: %s", err, errOutput)
	}

	var manifest map[string]interface{}
	err = yaml.Unmarshal([]byte(response), &manifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml", err)
	}

	return manifest, nil
}
