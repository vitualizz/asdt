package grader_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/vitualizz/asdt/internal/grader"
	"gopkg.in/yaml.v3"
)

// intp is a helper for the *int bound fields.
func intp(n int) *int { return &n }

// testdataDir resolves the testdata directory relative to this source file so
// tests are independent of the working directory.
func testdataDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "testdata")
}

// TestGrade exercises every probe type across pass/fail/error variants plus the
// three rollup cases.
func TestGrade(t *testing.T) {
	tests := []struct {
		name       string
		probe      grader.Probe
		payload    map[string]any
		wantStatus string
	}{
		// field-present
		{
			name:       "field-present pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldPresent, Target: "decision"},
			payload:    map[string]any{"decision": "do it"},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "field-present fail (absent)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldPresent, Target: "decision"},
			payload:    map[string]any{"other": "x"},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "field-present error (mismatch on path)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldPresent, Target: "a.b"},
			payload:    map[string]any{"a": "scalar"},
			wantStatus: grader.StatusError,
		},

		// field-absent
		{
			name:       "field-absent pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldAbsent, Target: "secret"},
			payload:    map[string]any{"other": "x"},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "field-absent fail (found)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldAbsent, Target: "secret"},
			payload:    map[string]any{"secret": "leaked"},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "field-absent error (mismatch on path)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldAbsent, Target: "a.b"},
			payload:    map[string]any{"a": "scalar"},
			wantStatus: grader.StatusError,
		},

		// negative-present
		{
			name:       "negative-present pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeNegativePresent, Target: "consequences.negative"},
			payload:    map[string]any{"consequences": map[string]any{"negative": []any{"tradeoff"}}},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "negative-present fail (empty list)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeNegativePresent, Target: "consequences.negative"},
			payload:    map[string]any{"consequences": map[string]any{"negative": []any{}}},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "negative-present fail (absent)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeNegativePresent, Target: "consequences.negative"},
			payload:    map[string]any{"consequences": map[string]any{}},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "negative-present error (not a list)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeNegativePresent, Target: "consequences"},
			payload:    map[string]any{"consequences": "scalar"},
			wantStatus: grader.StatusError,
		},

		// enum-member
		{
			name:       "enum-member pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeEnumMember, Target: "status", ExpectedValues: []string{"accepted", "rejected"}},
			payload:    map[string]any{"status": "accepted"},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "enum-member fail (not in set)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeEnumMember, Target: "status", ExpectedValues: []string{"accepted", "rejected"}},
			payload:    map[string]any{"status": "draft"},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "enum-member error (absent)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeEnumMember, Target: "status", ExpectedValues: []string{"accepted"}},
			payload:    map[string]any{"other": "x"},
			wantStatus: grader.StatusError,
		},

		// min-count
		{
			name:       "min-count pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMinCount, Target: "items", Min: intp(2)},
			payload:    map[string]any{"items": []any{"a", "b", "c"}},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "min-count fail (too few)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMinCount, Target: "items", Min: intp(3)},
			payload:    map[string]any{"items": []any{"a"}},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "min-count error (missing bound)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMinCount, Target: "items"},
			payload:    map[string]any{"items": []any{"a"}},
			wantStatus: grader.StatusError,
		},

		// max-count
		{
			name:       "max-count pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMaxCount, Target: "items", Max: intp(3)},
			payload:    map[string]any{"items": []any{"a", "b"}},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "max-count fail (too many)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMaxCount, Target: "items", Max: intp(1)},
			payload:    map[string]any{"items": []any{"a", "b"}},
			wantStatus: grader.StatusFail,
		},
		{
			name:       "max-count error (not a list)",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeMaxCount, Target: "items", Max: intp(1)},
			payload:    map[string]any{"items": "scalar"},
			wantStatus: grader.StatusError,
		},

		// fk-integrity
		{
			name: "fk-integrity pass",
			probe: grader.Probe{
				Name: "p", Type: grader.ProbeFKIntegrity,
				Target: "edges", FKTarget: "nodes", FKField: "id",
			},
			payload: map[string]any{
				"edges": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}},
				"nodes": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}},
			},
			wantStatus: grader.StatusPass,
		},
		{
			name: "fk-integrity fail (dangling)",
			probe: grader.Probe{
				Name: "p", Type: grader.ProbeFKIntegrity,
				Target: "edges", FKTarget: "nodes", FKField: "id",
			},
			payload: map[string]any{
				"edges": []any{map[string]any{"id": "a"}, map[string]any{"id": "z"}},
				"nodes": []any{map[string]any{"id": "a"}},
			},
			wantStatus: grader.StatusFail,
		},
		{
			name: "fk-integrity error (target absent)",
			probe: grader.Probe{
				Name: "p", Type: grader.ProbeFKIntegrity,
				Target: "edges", FKTarget: "nodes", FKField: "id",
			},
			payload: map[string]any{
				"nodes": []any{map[string]any{"id": "a"}},
			},
			wantStatus: grader.StatusError,
		},

		// rollup cases
		{
			name:       "rollup all pass",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeFieldPresent, Target: "decision"},
			payload:    map[string]any{"decision": "x"},
			wantStatus: grader.StatusPass,
		},
		{
			name:       "rollup unknown type errors",
			probe:      grader.Probe{Name: "p", Type: grader.ProbeType("bogus"), Target: "x"},
			payload:    map[string]any{"x": "y"},
			wantStatus: grader.StatusError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			set := grader.ProbeSet{Probes: []grader.Probe{tc.probe}}
			got := grader.Grade(set, tc.payload)
			if got.Status != tc.wantStatus {
				t.Errorf("Status: got %q, want %q (diagnostics: %v)", got.Status, tc.wantStatus, got.Diagnostics)
			}
		})
	}
}

// TestGrade_RollupPrecedence verifies a multi-probe set rolls up with
// ERROR > FAIL > PASS precedence and populates the name lists.
func TestGrade_RollupPrecedence(t *testing.T) {
	set := grader.ProbeSet{Probes: []grader.Probe{
		{Name: "ok", Type: grader.ProbeFieldPresent, Target: "decision"},
		{Name: "bad", Type: grader.ProbeFieldPresent, Target: "missing"},
		{Name: "broken", Type: grader.ProbeType("bogus"), Target: "x"},
	}}
	payload := map[string]any{"decision": "x"}
	got := grader.Grade(set, payload)

	if got.Status != grader.StatusError {
		t.Fatalf("Status: got %q, want ERROR", got.Status)
	}
	if len(got.ProbesPassed) != 1 || got.ProbesPassed[0] != "ok" {
		t.Errorf("ProbesPassed: got %v, want [ok]", got.ProbesPassed)
	}
	if len(got.ProbesFailed) != 2 {
		t.Errorf("ProbesFailed: got %v, want 2 entries", got.ProbesFailed)
	}
}

func TestResolve_ViaEvaluators(t *testing.T) {
	// resolve is unexported; exercise its semantics through the field-present
	// evaluator, which maps each pathState to a distinct Verdict status.
	tests := []struct {
		name       string
		target     string
		payload    map[string]any
		wantStatus string
	}{
		{"found top", "a", map[string]any{"a": 1}, grader.StatusPass},
		{"found nested", "a.b.c", map[string]any{"a": map[string]any{"b": map[string]any{"c": 1}}}, grader.StatusPass},
		{"absent top", "z", map[string]any{"a": 1}, grader.StatusFail},
		{"absent nested", "a.z", map[string]any{"a": map[string]any{"b": 1}}, grader.StatusFail},
		{"mismatch on path", "a.b", map[string]any{"a": "scalar"}, grader.StatusError},
		{"nil payload", "a", nil, grader.StatusFail},
		{"empty payload", "a", map[string]any{}, grader.StatusFail},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			set := grader.ProbeSet{Probes: []grader.Probe{
				{Name: "r", Type: grader.ProbeFieldPresent, Target: tc.target},
			}}
			got := grader.Grade(set, tc.payload)
			if got.Status != tc.wantStatus {
				t.Errorf("Status: got %q, want %q", got.Status, tc.wantStatus)
			}
		})
	}
}

func TestParseProbeType(t *testing.T) {
	valid := []grader.ProbeType{
		grader.ProbeFieldPresent, grader.ProbeFieldAbsent, grader.ProbeNegativePresent,
		grader.ProbeEnumMember, grader.ProbeFKIntegrity, grader.ProbeMinCount, grader.ProbeMaxCount,
	}
	for _, v := range valid {
		got, err := grader.ParseProbeType(string(v))
		if err != nil {
			t.Errorf("ParseProbeType(%q): unexpected error %v", v, err)
		}
		if got != v {
			t.Errorf("ParseProbeType(%q): got %q", v, got)
		}
	}

	_, err := grader.ParseProbeType("bogus")
	if !errors.Is(err, grader.ErrUnknownProbeType) {
		t.Errorf("ParseProbeType(bogus): want ErrUnknownProbeType, got %v", err)
	}
}

func TestValidateProbeSet(t *testing.T) {
	valid := grader.ProbeSet{Probes: []grader.Probe{
		{Name: "a", Type: grader.ProbeFieldPresent, Target: "x"},
	}}
	if err := grader.ValidateProbeSet(valid); err != nil {
		t.Errorf("valid set: unexpected error %v", err)
	}

	bogus := grader.ProbeSet{Probes: []grader.Probe{
		{Name: "a", Type: grader.ProbeType("nope"), Target: "x"},
	}}
	if err := grader.ValidateProbeSet(bogus); !errors.Is(err, grader.ErrUnknownProbeType) {
		t.Errorf("bogus set: want ErrUnknownProbeType, got %v", err)
	}
}

func TestFixtureSource_Load_Golden(t *testing.T) {
	root := filepath.Join(testdataDir(t), "probes")
	src := grader.NewFixtureSource(root)

	set, err := src.Load("architect", "decision-record")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Load the golden payload envelope and grade against it.
	payloadPath := filepath.Join(testdataDir(t), "payloads", "architect-decision-record.golden.yaml")
	data, err := os.ReadFile(payloadPath)
	if err != nil {
		t.Fatalf("read golden payload: %v", err)
	}
	var envelope map[string]any
	if err := yaml.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("unmarshal golden payload: %v", err)
	}

	got := grader.Grade(set, envelope)
	if got.Status != grader.StatusPass {
		t.Errorf("Status: got %q, want PASS (diagnostics: %v)", got.Status, got.Diagnostics)
	}
	if len(got.ProbesFailed) != 0 {
		t.Errorf("ProbesFailed: got %v, want empty", got.ProbesFailed)
	}
}

func TestFixtureSource_Load_NotFound(t *testing.T) {
	root := filepath.Join(testdataDir(t), "probes")
	src := grader.NewFixtureSource(root)

	_, err := src.Load("architect", "does-not-exist")
	if err == nil {
		t.Fatal("expected error loading missing probe set, got nil")
	}
}
