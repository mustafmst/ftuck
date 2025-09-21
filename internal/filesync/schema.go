package filesync

import (
	"bytes"
	"io"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type SyncDefinition struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dest"`
}

type Schema []SyncDefinition

func (s *Schema) WriteToFile(path string) error {
	slog.Debug("writing sync schema to file", "path", path, "entries", len(*s))
	
	d, err := yaml.Marshal(s)
	if err != nil {
		slog.Error("failed to marshal sync schema", "path", path, "error", err)
		return err
	}

	err = os.WriteFile(path, d, 0644)
	if err != nil {
		slog.Error("failed to write sync file", "path", path, "error", err)
		return err
	}

	slog.Debug("sync schema written successfully", "path", path)
	return nil
}

func (s *Schema) Append(definition SyncDefinition) {
	*s = append(*s, definition)
}

func (s *Schema) ForEach(f func(SyncDefinition) error) error {
	for _, sd := range *s {
		err := f(sd)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadOrCreate(path string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()
	io.Copy(buf, f)
	return buf.Bytes(), nil
}

func ReadSchema(data []byte) (*Schema, error) {
	res := Schema{}
	err := yaml.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
