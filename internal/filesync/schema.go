package filesync

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SyncDefinition struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Schema []SyncDefinition

func (s *Schema) WriteToFile(path string) error {
	d, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema) Append(definition SyncDefinition) {
	*s = append(*s, definition)
}

func ReadSchema(data []byte) (*Schema, error) {
	res := Schema{}
	err := yaml.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
