package grader

import "fmt"

// ProbeSource is the driven port for loading probe sets. Implementations may be
// filesystem-backed (FixtureSource), embedded, or in-memory for tests.
type ProbeSource interface {
	// Load returns the ProbeSet for the given specialist step.
	Load(specialistID, stepID string) (ProbeSet, error)
}

// ValidateProbeSet checks that every probe in the set declares a known type.
// It returns a wrapped ErrUnknownProbeType naming the first offending probe.
func ValidateProbeSet(set ProbeSet) error {
	for _, probe := range set.Probes {
		if !probe.Type.Valid() {
			return fmt.Errorf("probe %q: %w", probe.Name, ErrUnknownProbeType)
		}
	}
	return nil
}
