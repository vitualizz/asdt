package grader

import "strings"

// Verdict status constants. A Verdict is always exactly one of these.
const (
	StatusPass  = "PASS"
	StatusFail  = "FAIL"
	StatusError = "ERROR"
)

// pathState classifies the outcome of resolving a dot path within a payload.
type pathState int

const (
	// stateFound means the final key was reached and its value is in Val.
	stateFound pathState = iota
	// stateAbsent means a key along the path was missing (but every hop that
	// did exist was a navigable map).
	stateAbsent
	// stateMismatchOnPath means a non-final hop resolved to a non-map value,
	// so the path could not be walked. This is a structural error, distinct
	// from a simple absence.
	stateMismatchOnPath
)

// resolved is the outcome of resolve.
type resolved struct {
	Val   any
	State pathState
}

// resolve walks a dot-separated path through a YAML-decoded payload.
//
// Semantics (never panics, handles nil payload):
//   - Each non-final hop must be a map[string]any; if it is any other type the
//     result is stateMismatchOnPath.
//   - A missing key at any hop yields stateAbsent.
//   - Reaching the final key yields stateFound with its value (which may be nil).
func resolve(payload map[string]any, dotPath string) resolved {
	parts := strings.Split(dotPath, ".")
	var cur any = payload
	for i, key := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			// We needed to descend but the current value is not a map.
			return resolved{State: stateMismatchOnPath}
		}
		val, present := m[key]
		if !present {
			return resolved{State: stateAbsent}
		}
		if i == len(parts)-1 {
			return resolved{Val: val, State: stateFound}
		}
		cur = val
	}
	// Only reachable when dotPath is empty (Split returns one empty element);
	// treat as absent rather than panicking.
	return resolved{State: stateAbsent}
}

// ProbeResult is the outcome of evaluating a single probe.
type ProbeResult struct {
	ProbeName string
	Type      ProbeType
	Passed    bool
	Errored   bool
	Detail    string
}

// Verdict is the aggregate outcome of grading a ProbeSet against a payload.
type Verdict struct {
	Status       string        `yaml:"status"`
	ProbesPassed []string      `yaml:"probes_passed"`
	ProbesFailed []string      `yaml:"probes_failed"`
	Diagnostics  []string      `yaml:"diagnostics"`
	Results      []ProbeResult `yaml:"results"`
}

// Grade evaluates every probe in set against payload and rolls the per-probe
// results up into a single Verdict. It is pure, total, and never panics.
//
// Rollup precedence: any errored probe ⇒ ERROR; otherwise any failed probe ⇒
// FAIL; otherwise PASS.
func Grade(set ProbeSet, payload map[string]any) Verdict {
	v := Verdict{Status: StatusPass}

	anyErrored := false
	anyFailed := false

	for _, probe := range set.Probes {
		eval, ok := evaluators[probe.Type]
		var res ProbeResult
		if !ok {
			res = ProbeResult{
				ProbeName: probe.Name,
				Type:      probe.Type,
				Errored:   true,
				Detail:    "unknown probe type",
			}
		} else {
			res = eval(probe, payload)
		}

		v.Results = append(v.Results, res)

		switch {
		case res.Errored:
			anyErrored = true
			v.ProbesFailed = append(v.ProbesFailed, res.ProbeName)
			v.Diagnostics = append(v.Diagnostics, diag(res))
		case !res.Passed:
			anyFailed = true
			v.ProbesFailed = append(v.ProbesFailed, res.ProbeName)
			v.Diagnostics = append(v.Diagnostics, diag(res))
		default:
			v.ProbesPassed = append(v.ProbesPassed, res.ProbeName)
		}
	}

	switch {
	case anyErrored:
		v.Status = StatusError
	case anyFailed:
		v.Status = StatusFail
	default:
		v.Status = StatusPass
	}

	return v
}

// diag formats a single non-passing probe result for the Diagnostics list.
func diag(res ProbeResult) string {
	if res.Detail == "" {
		return res.ProbeName
	}
	return res.ProbeName + ": " + res.Detail
}
