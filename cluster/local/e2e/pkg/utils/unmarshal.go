package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// UnmarshalAnyYaml unmarshals a yaml file into a struct
func UnmarshalAnyYaml[T any](path string, out *T) (*T, error) {
	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	byteValue, _ := io.ReadAll(jsonFile)
	err = jsonFile.Close()
	if err != nil {
		return nil, err
	}
	j, err := yaml.YAMLToJSON(byteValue)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(j, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
