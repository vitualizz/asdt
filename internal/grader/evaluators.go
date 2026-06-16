package grader

import (
	"fmt"
	"strings"
)

// Evaluator evaluates a single probe against a payload and returns its result.
// Evaluators are pure and must never panic, even on malformed payloads.
type Evaluator func(p Probe, payload map[string]any) ProbeResult

// evaluators is the dispatch map from probe type to its evaluator. Adding a new
// probe type requires adding an entry here plus the ProbeType constant.
var evaluators = map[ProbeType]Evaluator{
	ProbeFieldPresent:    evalFieldPresent,
	ProbeFieldAbsent:     evalFieldAbsent,
	ProbeNegativePresent: evalNegativePresent,
	ProbeEnumMember:      evalEnumMember,
	ProbeFKIntegrity:     evalFKIntegrity,
	ProbeMinCount:        evalMinCount,
	ProbeMaxCount:        evalMaxCount,
}

func pass(p Probe) ProbeResult {
	return ProbeResult{ProbeName: p.Name, Type: p.Type, Passed: true}
}

func fail(p Probe, detail string) ProbeResult {
	return ProbeResult{ProbeName: p.Name, Type: p.Type, Detail: detail}
}

func errored(p Probe, detail string) ProbeResult {
	return ProbeResult{ProbeName: p.Name, Type: p.Type, Errored: true, Detail: detail}
}

func evalFieldPresent(p Probe, payload map[string]any) ProbeResult {
	r := resolve(payload, p.Target)
	switch r.State {
	case stateFound:
		return pass(p)
	case stateAbsent:
		return fail(p, fmt.Sprintf("expected field %q to be present", p.Target))
	default: // stateMismatchOnPath
		return errored(p, fmt.Sprintf("path %q traverses a non-map value", p.Target))
	}
}

func evalFieldAbsent(p Probe, payload map[string]any) ProbeResult {
	r := resolve(payload, p.Target)
	switch r.State {
	case stateAbsent:
		return pass(p)
	case stateFound:
		return fail(p, fmt.Sprintf("expected field %q to be absent", p.Target))
	default: // stateMismatchOnPath
		return errored(p, fmt.Sprintf("path %q traverses a non-map value", p.Target))
	}
}

func evalNegativePresent(p Probe, payload map[string]any) ProbeResult {
	r := resolve(payload, p.Target)
	switch r.State {
	case stateAbsent:
		return fail(p, fmt.Sprintf("expected non-empty list at %q, but it is absent", p.Target))
	case stateMismatchOnPath:
		return errored(p, fmt.Sprintf("path %q traverses a non-map value", p.Target))
	}
	list, ok := r.Val.([]any)
	if !ok {
		return errored(p, fmt.Sprintf("expected %q to be a list, got %T", p.Target, r.Val))
	}
	if len(list) > 0 {
		return pass(p)
	}
	return fail(p, fmt.Sprintf("expected %q to be non-empty", p.Target))
}

func evalEnumMember(p Probe, payload map[string]any) ProbeResult {
	r := resolve(payload, p.Target)
	switch r.State {
	case stateAbsent:
		return errored(p, fmt.Sprintf("field %q is absent; cannot check enum membership", p.Target))
	case stateMismatchOnPath:
		return errored(p, fmt.Sprintf("path %q traverses a non-map value", p.Target))
	}
	got := fmt.Sprint(r.Val)
	for _, want := range p.ExpectedValues {
		if got == want {
			return pass(p)
		}
	}
	return fail(p, fmt.Sprintf("value %q not in allowed set [%s]", got, strings.Join(p.ExpectedValues, ", ")))
}

func evalMinCount(p Probe, payload map[string]any) ProbeResult {
	if p.Min == nil {
		return errored(p, "missing bound: min is required for min-count")
	}
	list, res, ok := resolveList(p, payload)
	if !ok {
		return res
	}
	if len(list) >= *p.Min {
		return pass(p)
	}
	return fail(p, fmt.Sprintf("expected at least %d items at %q, got %d", *p.Min, p.Target, len(list)))
}

func evalMaxCount(p Probe, payload map[string]any) ProbeResult {
	if p.Max == nil {
		return errored(p, "missing bound: max is required for max-count")
	}
	list, res, ok := resolveList(p, payload)
	if !ok {
		return res
	}
	if len(list) <= *p.Max {
		return pass(p)
	}
	return fail(p, fmt.Sprintf("expected at most %d items at %q, got %d", *p.Max, p.Target, len(list)))
}

// resolveList resolves p.Target to a []any. On any failure it returns ok=false
// and a fully-formed errored ProbeResult to return verbatim.
func resolveList(p Probe, payload map[string]any) (list []any, res ProbeResult, ok bool) {
	r := resolve(payload, p.Target)
	switch r.State {
	case stateAbsent:
		return nil, errored(p, fmt.Sprintf("field %q is absent", p.Target)), false
	case stateMismatchOnPath:
		return nil, errored(p, fmt.Sprintf("path %q traverses a non-map value", p.Target)), false
	}
	l, isList := r.Val.([]any)
	if !isList {
		return nil, errored(p, fmt.Sprintf("expected %q to be a list, got %T", p.Target, r.Val)), false
	}
	return l, ProbeResult{}, true
}

func evalFKIntegrity(p Probe, payload map[string]any) ProbeResult {
	// Resolve the referencing collection.
	srcItems, srcOK := resolveMapList(payload, p.Target)
	if !srcOK {
		return errored(p, fmt.Sprintf("target %q must be a list of objects", p.Target))
	}
	// Resolve the referenced collection.
	refItems, refOK := resolveMapList(payload, p.FKTarget)
	if !refOK {
		return errored(p, fmt.Sprintf("fk_target %q must be a list of objects", p.FKTarget))
	}

	// Build the set of valid key values from the referenced collection.
	valid := make(map[string]struct{}, len(refItems))
	for _, item := range refItems {
		if v, present := item[p.FKField]; present {
			valid[fmt.Sprint(v)] = struct{}{}
		}
	}

	// Every referencing item's FK value must exist in the valid set.
	var dangling []string
	for _, item := range srcItems {
		v, present := item[p.FKField]
		if !present {
			dangling = append(dangling, "<missing>")
			continue
		}
		key := fmt.Sprint(v)
		if _, found := valid[key]; !found {
			dangling = append(dangling, key)
		}
	}

	if len(dangling) == 0 {
		return pass(p)
	}
	return fail(p, fmt.Sprintf("dangling %s references in %q: [%s]",
		p.FKField, p.Target, strings.Join(dangling, ", ")))
}

// resolveMapList resolves dotPath to a []any whose elements are all
// map[string]any. Returns ok=false on absence, path mismatch, or wrong shape.
func resolveMapList(payload map[string]any, dotPath string) ([]map[string]any, bool) {
	r := resolve(payload, dotPath)
	if r.State != stateFound {
		return nil, false
	}
	list, ok := r.Val.([]any)
	if !ok {
		return nil, false
	}
	out := make([]map[string]any, 0, len(list))
	for _, el := range list {
		m, isMap := el.(map[string]any)
		if !isMap {
			return nil, false
		}
		out = append(out, m)
	}
	return out, true
}
