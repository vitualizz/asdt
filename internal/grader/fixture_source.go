package grader

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FixtureSource loads probe sets from YAML files laid out on disk as
// {Root}/{specialistID}/{stepID}.yaml.
type FixtureSource struct {
	Root string
}

// NewFixtureSource returns a FixtureSource rooted at the given directory.
func NewFixtureSource(root string) *FixtureSource {
	return &FixtureSource{Root: root}
}

// Load reads, decodes, and validates the probe set for the given step.
func (s *FixtureSource) Load(specialistID, stepID string) (ProbeSet, error) {
	p := filepath.Join(s.Root, specialistID, stepID+".yaml")

	data, err := os.ReadFile(p)
	if err != nil {
		return ProbeSet{}, fmt.Errorf("probeset read %s/%s: %w", specialistID, stepID, err)
	}

	var set ProbeSet
	if err := yaml.Unmarshal(data, &set); err != nil {
		return ProbeSet{}, fmt.Errorf("probeset unmarshal %s/%s: %w", specialistID, stepID, err)
	}

	if err := ValidateProbeSet(set); err != nil {
		return ProbeSet{}, fmt.Errorf("probeset validate %s/%s: %w", specialistID, stepID, err)
	}

	return set, nil
}

// Compile-time assertion that FixtureSource satisfies ProbeSource.
var _ ProbeSource = (*FixtureSource)(nil)
