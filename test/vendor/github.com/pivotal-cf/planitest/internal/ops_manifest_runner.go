package internal

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type OpsManifestRunner struct {
	cmdRunner CommandRunner
	FileIO    fileIO
}

func NewOpsManifestRunner(cmdRunner CommandRunner) OpsManifestRunner {
	return OpsManifestRunner{
		cmdRunner: cmdRunner,
		FileIO:    FileIO{},
	}
}

func (o OpsManifestRunner) GetManifest(productProperties, metadataFilePath string) (map[string]interface{}, error) {
	configFile, err := o.FileIO.TempFile("", "")
	configFileJson := fmt.Sprintf("%s.json", configFile.Name())
	os.Rename(configFile.Name(), configFileJson)

	if err != nil {
		return nil, err //not tested
	}

	_, err = configFile.WriteString(productProperties)
	if err != nil {
		return nil, err //not tested
	}

	response, errOutput, err := o.cmdRunner.Run(
		"ops-manifest",
		"--config-file", configFileJson,
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
