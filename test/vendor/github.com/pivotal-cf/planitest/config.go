package planitest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yamlConverter "github.com/ghodss/yaml"
)

// Things I don't like
// Arguments:
// - configFile is a path to a file, requires IO to read
// - additionalProperties is JSON, decoded into a map
// additionalProperties is really only additionalProductProperties, and gets merged with configFile's productProperties section
func MergeAdditionalProductProperties(configFile string, additionalProperties map[string]interface{}) ([]byte, error) {
	yamlInput, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to configure product: %s (config file %s)", err, configFile)
	}
	jsonInput, err := yamlConverter.YAMLToJSON(yamlInput)
	if err != nil {
		return nil, err
	}

	var inputConfig ConfigJSON
	err = json.Unmarshal(jsonInput, &inputConfig)

	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %s", err)
	}

	if inputConfig.NetworkProperties == nil {
		return nil, fmt.Errorf("network properties must be provided in the config file")
	}

	if inputConfig.ProductProperties == nil {
		return nil, fmt.Errorf("product properties must be provided in the config file")
	}

	propertiesJSON, err := json.Marshal(inputConfig.ProductProperties)
	if err != nil {
		return nil, err // un-tested
	}

	var minimalProperties *map[string]interface{}
	err = json.Unmarshal(propertiesJSON, &minimalProperties)
	if err != nil {
		return nil, fmt.Errorf("could not parse product properties: %s", err)
	}

	combinedProperties := mergeProperties(*minimalProperties, additionalProperties)

	propertiesJSON, err = json.Marshal(combinedProperties)
	if err != nil {
		return nil, err // un-tested
	}

	err = json.Unmarshal(propertiesJSON, &inputConfig.ProductProperties)
	if err != nil {
		return nil, err // un-tested
	}

	outputJSON, err := json.Marshal(inputConfig)
	if err != nil {
		return nil, err
	}

	return outputJSON, nil
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
