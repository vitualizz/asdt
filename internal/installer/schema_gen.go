package installer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SchemaKind classifies a canonical schema by the role it plays in the
// schema-as-source-of-truth pipeline. The kind decides whether a schema drives
// a machine-generated inline block in a producing step .md (final_artifact) or
// is documentation-only with no inline mirror (input_platform, wrapper).
type SchemaKind string

const (
	// SchemaKindFinalArtifact is a schema whose example payload is INLINE-RENDERED
	// into a producing specialist step .md between generated-region markers. These
	// are held to a byte-equal contract (committed == rendered).
	SchemaKindFinalArtifact SchemaKind = "final_artifact"
	// SchemaKindInputPlatform is a documentation-only input/platform schema. It has
	// no inline mirror and is never rendered into a step .md.
	SchemaKindInputPlatform SchemaKind = "input_platform"
	// SchemaKindWrapper is the transport envelope schema wrapping specialist
	// payloads. Documentation-only; never inline-rendered.
	SchemaKindWrapper SchemaKind = "wrapper"
)

// SchemaSpec is the canonical registry entry for one schemas/*.schema.yaml file.
// For SchemaKindFinalArtifact all fields are set; for the other kinds only
// Name, SchemaFile, and Kind carry values (the render fields stay empty).
type SchemaSpec struct {
	Name        string // logical artifact name, e.g. "architectural-decision"
	SchemaFile  string // path relative to schemas/, e.g. "developer/dev-spec.schema.yaml"
	Kind        SchemaKind
	MDPath      string // producing step .md relative to skill/ (final_artifact only)
	FenceLabel  string // label line preceding the ```yaml fence (final_artifact only)
	MarkerBegin string // generated-region begin marker (final_artifact only)
	MarkerEnd   string // generated-region end marker (final_artifact only)
}

// schemaSpecs is the ORDERED, canonical registry of every schemas/*.schema.yaml
// file. Order is fixed (a slice, never a map) so generation and coverage checks
// are deterministic. The 7 final_artifact rows carry values copied VERBATIM from
// the drift-guard table so the committed inline blocks stay byte-identical.
var schemaSpecs = []SchemaSpec{
	{
		Name:        "architectural-decision",
		SchemaFile:  "architectural-decision.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-architect/steps/technical-handoff.md",
		FenceLabel:  "architectural-decision schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-architectural-decision -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-architectural-decision -->",
	},
	{
		Name:        "system-design",
		SchemaFile:  "system-design.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-architect/steps/technical-handoff.md",
		FenceLabel:  "system-design schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-system-design -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-system-design -->",
	},
	{
		Name:        "ux-brief",
		SchemaFile:  "ux-brief.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-ux-ui/steps/ux-handoff.md",
		FenceLabel:  "ux-brief schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-ux-brief -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-ux-brief -->",
	},
	{
		Name:        "component-spec",
		SchemaFile:  "component-spec.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-ux-ui/steps/ux-handoff.md",
		FenceLabel:  "component-spec schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-component-spec -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-component-spec -->",
	},
	{
		Name:        "security-findings",
		SchemaFile:  "security-findings.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-security/steps/hardening-checklist.md",
		FenceLabel:  "security-findings schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-security-findings -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-security-findings -->",
	},
	{
		Name:        "hardening-checklist",
		SchemaFile:  "hardening-checklist.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-security/steps/hardening-checklist.md",
		FenceLabel:  "hardening-checklist schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-hardening-checklist -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-hardening-checklist -->",
	},
	{
		Name:        "test-plan",
		SchemaFile:  "test-plan.schema.yaml",
		Kind:        SchemaKindFinalArtifact,
		MDPath:      "asdt-qa/steps/quality-report.md",
		FenceLabel:  "Schema:",
		MarkerBegin: "<!-- ASDT:GENERATED:schema-test-plan -->",
		MarkerEnd:   "<!-- /ASDT:GENERATED:schema-test-plan -->",
	},

	// Wrapper (transport envelope) — documentation-only, never inline-rendered.
	{
		Name:       "envelope",
		SchemaFile: "envelope.schema.yaml",
		Kind:       SchemaKindWrapper,
	},

	// Input / platform schemas — documentation-only, never inline-rendered.
	{
		Name:       "platform",
		SchemaFile: "platform.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "platform-summary",
		SchemaFile: "platform-summary.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "project-context",
		SchemaFile: "project-context.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "requirements-spec",
		SchemaFile: "requirements-spec.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "architect-approaches",
		SchemaFile: "architect-approaches.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "architect-constraints-analysis",
		SchemaFile: "architect-constraints-analysis.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "implementation-plan",
		SchemaFile: "implementation-plan.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "qa-edge-cases",
		SchemaFile: "qa-edge-cases.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "risk-register",
		SchemaFile: "risk-register.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "threat-model",
		SchemaFile: "threat-model.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "quality-report",
		SchemaFile: "quality-report.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "dev-design",
		SchemaFile: "developer/dev-design.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
	{
		Name:       "dev-spec",
		SchemaFile: "developer/dev-spec.schema.yaml",
		Kind:       SchemaKindInputPlatform,
	},
}

// checkMarkersPresent verifies that begin and end each appear EXACTLY once in
// content. replaceMarkerRegion silently no-ops when BOTH markers are absent (so
// marker-free fixtures are tolerated); this check exists so the generator fails
// LOUD when a producing .md is missing (or duplicates) the markers it must
// maintain, naming the offending marker and its count.
func checkMarkersPresent(content []byte, begin, end string) error {
	if n := strings.Count(string(content), begin); n != 1 {
		return fmt.Errorf("begin marker %q appears %d times, want exactly 1", begin, n)
	}
	if n := strings.Count(string(content), end); n != 1 {
		return fmt.Errorf("end marker %q appears %d times, want exactly 1", end, n)
	}
	return nil
}

// renderSchemaRegionForSpec reads a final-artifact spec's canonical schema and
// its producing step .md, renders the inline example region from the schema, and
// splices it between the spec's markers. It returns the updated .md bytes and
// whether they differ from what was on disk. All errors are wrapped with the
// spec name and the relevant file path.
func renderSchemaRegionForSpec(root string, spec SchemaSpec) (updated []byte, changed bool, err error) {
	schemaPath := filepath.Join(root, "schemas", filepath.FromSlash(spec.SchemaFile))
	mdPath := filepath.Join(root, "skill", filepath.FromSlash(spec.MDPath))

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, false, fmt.Errorf("%s: read schema %s: %w", spec.Name, schemaPath, err)
	}
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, false, fmt.Errorf("%s: read producing md %s: %w", spec.Name, mdPath, err)
	}

	if err := checkMarkersPresent(mdBytes, spec.MarkerBegin, spec.MarkerEnd); err != nil {
		return nil, false, fmt.Errorf("%s: markers in %s: %w", spec.Name, mdPath, err)
	}

	body, err := renderSchemaExampleRegion(schemaBytes, spec.FenceLabel)
	if err != nil {
		return nil, false, fmt.Errorf("%s: render example region from %s: %w", spec.Name, schemaPath, err)
	}

	updated, err = replaceMarkerRegion(mdBytes, spec.MarkerBegin, spec.MarkerEnd, body)
	if err != nil {
		return nil, false, fmt.Errorf("%s: splice region into %s: %w", spec.Name, mdPath, err)
	}

	return updated, !bytes.Equal(updated, mdBytes), nil
}

// GenerateSchemaRegions regenerates every final-artifact inline block from its
// canonical schema and writes the updated producing step .md files under root.
// Non-final schemas are skipped. Iteration follows schemaSpecs order for
// determinism, and files are written only when their region actually changed
// (preserving the existing file mode). Errors are named per spec.
func GenerateSchemaRegions(root string) error {
	for _, spec := range schemaSpecs {
		if spec.Kind != SchemaKindFinalArtifact {
			continue
		}

		updated, changed, err := renderSchemaRegionForSpec(root, spec)
		if err != nil {
			return err
		}
		if !changed {
			continue
		}

		mdPath := filepath.Join(root, "skill", filepath.FromSlash(spec.MDPath))
		info, err := os.Stat(mdPath)
		if err != nil {
			return fmt.Errorf("%s: stat %s: %w", spec.Name, mdPath, err)
		}
		if err := os.WriteFile(mdPath, updated, info.Mode()); err != nil {
			return fmt.Errorf("%s: write %s: %w", spec.Name, mdPath, err)
		}
	}
	return nil
}
