// Package grader defines a pure, panic-free probe grader. A ProbeSet is a
// declarative list of probes that are evaluated against an artifact payload
// (a map[string]any decoded from YAML) to produce a Verdict.
// This package imports only stdlib and gopkg.in/yaml.v3.
package grader

import (
	"errors"
	"fmt"
)

// ProbeType is the kind of check a probe performs. Values are kebab-case so
// they read naturally in probe-set YAML.
type ProbeType string

// The supported probe types. Each maps to exactly one evaluator in the
// dispatch map (see evaluators.go).
const (
	ProbeFieldPresent    ProbeType = "field-present"
	ProbeFieldAbsent     ProbeType = "field-absent"
	ProbeNegativePresent ProbeType = "negative-present"
	ProbeEnumMember      ProbeType = "enum-member"
	ProbeFKIntegrity     ProbeType = "fk-integrity"
	ProbeMinCount        ProbeType = "min-count"
	ProbeMaxCount        ProbeType = "max-count"
)

// ErrUnknownProbeType is returned (wrapped) when a probe type string does not
// match any known ProbeType.
var ErrUnknownProbeType = errors.New("unknown probe type")

// Valid reports whether t is one of the known probe types.
func (t ProbeType) Valid() bool {
	switch t {
	case ProbeFieldPresent, ProbeFieldAbsent, ProbeNegativePresent,
		ProbeEnumMember, ProbeFKIntegrity, ProbeMinCount, ProbeMaxCount:
		return true
	default:
		return false
	}
}

// ParseProbeType converts a string to a ProbeType, returning a wrapped
// ErrUnknownProbeType when the value is not recognized.
func ParseProbeType(s string) (ProbeType, error) {
	t := ProbeType(s)
	if !t.Valid() {
		return "", fmt.Errorf("%q: %w", s, ErrUnknownProbeType)
	}
	return t, nil
}

// Probe is a single declarative check within a ProbeSet.
type Probe struct {
	Name           string    `yaml:"name"`
	Type           ProbeType `yaml:"type"`
	Target         string    `yaml:"target"`
	ExpectedValues []string  `yaml:"expected_values,omitempty"`
	Min            *int      `yaml:"min,omitempty"`
	Max            *int      `yaml:"max,omitempty"`
	FKTarget       string    `yaml:"fk_target,omitempty"`
	FKField        string    `yaml:"fk_field,omitempty"`
}

// ProbeSet is a versioned collection of probes scoped to a specialist step.
type ProbeSet struct {
	SchemaVersion string  `yaml:"schema_version"`
	SpecialistID  string  `yaml:"specialist_id"`
	StepID        string  `yaml:"step_id"`
	MVPGate       string  `yaml:"mvp_gate,omitempty"`
	Probes        []Probe `yaml:"probes"`
}
