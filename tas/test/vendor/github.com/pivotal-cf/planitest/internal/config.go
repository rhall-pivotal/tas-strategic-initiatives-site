package internal

import (
	"fmt"
	"io"
	"io/ioutil"

	"bytes"

	"github.com/pivotal-cf/om/config"
	"gopkg.in/yaml.v2"
)

// Things I don't like
// Arguments:
// - additionalProperties is JSON, decoded into a map
// additionalProperties is really only additionalProductProperties, and gets merged with configFile's productProperties section
// Starts with an io.Reader to a YAML structured config file, returns a JSON structured string/[]byte containing the added product properties

//MergeAdditionalProductProperties
func MergeAdditionalProductProperties(configFile io.Reader, additionalProperties map[string]interface{}) (io.Reader, error) {
	yamlInput, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var inputConfig config.ProductConfiguration
	err = yaml.Unmarshal(yamlInput, &inputConfig)

	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %s", err)
	}

	if inputConfig.NetworkProperties == nil {
		return nil, fmt.Errorf("network properties must be provided in the config file")
	}

	if inputConfig.ProductProperties == nil {
		return nil, fmt.Errorf("product properties must be provided in the config file")
	}

	inputConfig.ProductProperties = mergeProperties(inputConfig.ProductProperties, additionalProperties)

	modifiedConfigFile := bytes.NewBufferString("")
	err = yaml.NewEncoder(modifiedConfigFile).Encode(&inputConfig)
	if err != nil {
		return nil, err
	}

	return modifiedConfigFile, nil
}

func mergeProperties(minimalProperties, additionalProperties map[string]interface{}) map[string]interface{} {
	combinedProperties := make(map[string]interface{}, len(minimalProperties)+len(additionalProperties))
	for k, v := range minimalProperties {
		combinedProperties[k] = v
	}

	for k, v := range additionalProperties {
		combinedProperties[k] = map[string]interface{}{
			"value": v,
		}
	}
	return combinedProperties
}
