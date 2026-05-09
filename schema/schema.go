package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VarDefinition describes a single environment variable contract.
type VarDefinition struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
	Pattern     string `yaml:"pattern"`
}

// Schema represents the top-level envlock schema file.
type Schema struct {
	Version string                    `yaml:"version"`
	Vars    map[string]VarDefinition  `yaml:"vars"`
}

// Load reads and parses an envlock schema file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file %q: %w", path, err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema file %q: %w", path, err)
	}

	if s.Version == "" {
		s.Version = "1"
	}

	if s.Vars == nil {
		s.Vars = make(map[string]VarDefinition)
	}

	return &s, nil
}
