package installer

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestSchemaInlineDrift is a DETECT-ONLY drift guard (schema-sot Wave 1, P1:
// US-001 + US-002). For each FINAL artifact it compares the TOP-LEVEL field-name
// SET of the inline `payload:` block embedded in the producing step .md against
// the TOP-LEVEL `properties:` field-name SET of its canonical
// schemas/<name>.schema.yaml, and fails ONLY when the divergence differs from
// the documented BASELINE below (a ratchet). This lets the guard land GREEN over
// the existing, known drift while failing loudly the moment NEW divergence is
// introduced.
//
// P1 scope is intentionally SHALLOW: only top-level key SETS are compared. FUTURE
// P2 tightening (recursive/nested key comparison + type-aware checks) is
// deliberately out of scope here and tracked under schema-sot P2/architect.

// schemaArtifact maps a FINAL artifact to its producing step .md, its canonical
// schema file, and the fence label that identifies its inline payload block
// within a (possibly multi-fence) step .md. A static table is used deliberately:
// the set of finals is fixed and small (7), so an explicit table is clearer and
// more robust than re-deriving the mapping from `Produces:` markers at runtime.
type schemaArtifact struct {
	name       string
	mdPath     string // relative to skill/
	schemaFile string // relative to schemas/
	fenceLabel string // substring of the label line preceding the ```yaml fence

	// generated marks an artifact whose inline block is now MACHINE-GENERATED
	// from its canonical schema (schema-sot Wave 2). Generated rows leave the
	// set-parity drift path (TestSchemaInlineDrift) and are instead held to a
	// strict byte-equal contract by TestSchemaInlineGenerated. markerBegin /
	// markerEnd bound that generated region inside the step .md and are empty
	// for non-generated rows.
	generated   bool
	markerBegin string
	markerEnd   string
}

var schemaArtifacts = []schemaArtifact{
	{
		name:        "architectural-decision",
		mdPath:      "asdt-architect/steps/technical-handoff.md",
		schemaFile:  "architectural-decision.schema.yaml",
		fenceLabel:  "architectural-decision schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-architectural-decision -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-architectural-decision -->",
	},
	{
		name:        "system-design",
		mdPath:      "asdt-architect/steps/technical-handoff.md",
		schemaFile:  "system-design.schema.yaml",
		fenceLabel:  "system-design schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-system-design -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-system-design -->",
	},
	{
		name:        "ux-brief",
		mdPath:      "asdt-ux-ui/steps/ux-handoff.md",
		schemaFile:  "ux-brief.schema.yaml",
		fenceLabel:  "ux-brief schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-ux-brief -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-ux-brief -->",
	},
	{
		name:        "component-spec",
		mdPath:      "asdt-ux-ui/steps/ux-handoff.md",
		schemaFile:  "component-spec.schema.yaml",
		fenceLabel:  "component-spec schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-component-spec -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-component-spec -->",
	},
	{
		name:        "security-findings",
		mdPath:      "asdt-security/steps/hardening-checklist.md",
		schemaFile:  "security-findings.schema.yaml",
		fenceLabel:  "security-findings schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-security-findings -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-security-findings -->",
	},
	{
		name:        "hardening-checklist",
		mdPath:      "asdt-security/steps/hardening-checklist.md",
		schemaFile:  "hardening-checklist.schema.yaml",
		fenceLabel:  "hardening-checklist schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-hardening-checklist -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-hardening-checklist -->",
	},
	{
		name:        "test-plan",
		mdPath:      "asdt-qa/steps/quality-report.md",
		schemaFile:  "test-plan.schema.yaml",
		fenceLabel:  "Schema:",
		generated:   true,
		markerBegin: "<!-- ASDT:GENERATED:schema-test-plan -->",
		markerEnd:   "<!-- /ASDT:GENERATED:schema-test-plan -->",
	},
}

// Known pre-existing schema↔inline drift pending reconciliation (schema-sot
// P2/architect). Do NOT add entries to silence NEW drift — reconcile or update
// the canonical source instead.
//
// Extra   = sorted(inlineKeys \ schemaKeys) — keys present inline but not in the schema.
// Missing = sorted(schemaKeys \ inlineKeys) — keys required by the schema but absent inline.
var schemaDriftBaseline = map[string]struct {
	Extra   []string
	Missing []string
}{
	// architectural-decision is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"architectural-decision": {
		Extra:   []string{},
		Missing: []string{},
	},
	// system-design is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"system-design": {
		Extra:   []string{},
		Missing: []string{},
	},
	// ux-brief is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"ux-brief": {
		Extra:   []string{},
		Missing: []string{},
	},
	// component-spec is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"component-spec": {
		Extra:   []string{},
		Missing: []string{},
	},
	// security-findings is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"security-findings": {
		Extra:   []string{},
		Missing: []string{},
	},
	// hardening-checklist is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"hardening-checklist": {
		Extra:   []string{},
		Missing: []string{},
	},
	// test-plan is now GENERATED (schema-sot Wave 2): committed == rendered
	// is asserted byte-for-byte by TestSchemaInlineGenerated, so its set-parity
	// baseline is EMPTY by invariant R-006 (baseline empty IFF generated==true).
	"test-plan": {
		Extra:   []string{},
		Missing: []string{},
	},
}

// schemaDir resolves the repository's schemas/ directory relative to this
// package directory (internal/installer), where `go test` runs. schemas/ lives
// at the repo root and is NOT embedded, so it is read via this relative path —
// mirroring the skillDir(t) idiom in workflow_agents_test.go.
func schemaDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join("..", "..", "schemas")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("cannot locate schemas/ relative to internal/installer: %v", err)
	}
	return dir
}

// fenceRe matches a fenced ```yaml ... ``` block, capturing its body.
var fenceRe = regexp.MustCompile("(?s)```yaml\\s*\\n(.*?)\\n```")

// inlinePayloadKeys extracts the set of TOP-LEVEL keys under `payload:` from the
// ```yaml fence in mdPath whose immediately-preceding non-blank line contains
// fenceLabel. Returns the keys and a found flag for discovery-integrity checks.
func inlinePayloadKeys(t *testing.T, mdPath, fenceLabel string) ([]string, bool) {
	t.Helper()
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("read %s: %v", mdPath, err)
	}
	content := string(data)

	locs := fenceRe.FindAllStringSubmatchIndex(content, -1)
	for _, loc := range locs {
		// loc[0] = fence start, loc[2:4] = captured body.
		preceding := content[:loc[0]]
		// label line = last non-blank line before the fence.
		lines := strings.Split(preceding, "\n")
		var label string
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				label = lines[i]
				break
			}
		}
		if !strings.Contains(label, fenceLabel) {
			continue
		}
		body := content[loc[2]:loc[3]]
		var doc map[string]yaml.Node
		if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
			t.Fatalf("parse inline payload fence (%s) in %s: %v", fenceLabel, mdPath, err)
		}
		payload, ok := doc["payload"]
		if !ok {
			return nil, false
		}
		return mappingKeys(payload), true
	}
	return nil, false
}

// schemaPropertyKeys extracts the set of TOP-LEVEL keys under `properties:` from
// a JSON-Schema-style YAML schema file.
func schemaPropertyKeys(t *testing.T, schemaPath string) ([]string, bool) {
	t.Helper()
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, false
	}
	var doc map[string]yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse schema %s: %v", schemaPath, err)
	}
	props, ok := doc["properties"]
	if !ok {
		return nil, false
	}
	return mappingKeys(props), true
}

// mappingKeys returns the keys of a YAML mapping node. Values are ignored, so
// comments and placeholder values are tolerated.
func mappingKeys(node yaml.Node) []string {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keys = append(keys, node.Content[i].Value)
	}
	return keys
}

// diff returns sorted(a \ b).
func diff(a, b []string) []string {
	set := make(map[string]bool, len(b))
	for _, k := range b {
		set[k] = true
	}
	var out []string
	for _, k := range a {
		if !set[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// extraDiff reports which keys in got are not in want — used to name NEW drift.
func extraDiff(got, want []string) []string {
	set := make(map[string]bool, len(want))
	for _, k := range want {
		set[k] = true
	}
	var out []string
	for _, k := range got {
		if !set[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func TestSchemaInlineDrift(t *testing.T) {
	skills := skillDir(t)
	schemas := schemaDir(t)

	for _, art := range schemaArtifacts {
		art := art
		t.Run(art.name, func(t *testing.T) {
			if art.generated {
				t.Skip("generated rows covered by TestSchemaInlineGenerated")
				return
			}
			mdPath := filepath.Join(skills, filepath.FromSlash(art.mdPath))
			schemaPath := filepath.Join(schemas, art.schemaFile)

			schemaKeys, ok := schemaPropertyKeys(t, schemaPath)
			if !ok {
				t.Fatalf("%s: schema file %s missing or has no top-level properties", art.name, schemaPath)
			}

			inlineKeys, ok := inlinePayloadKeys(t, mdPath, art.fenceLabel)
			if !ok {
				t.Fatalf("%s: could not resolve inline payload fence %q in %s", art.name, art.fenceLabel, mdPath)
			}

			extra := diff(inlineKeys, schemaKeys)
			missing := diff(schemaKeys, inlineKeys)

			base, baselined := schemaDriftBaseline[art.name]
			if !baselined {
				t.Fatalf("%s: no baseline entry — every final artifact must have a documented baseline", art.name)
			}

			newExtra := extraDiff(extra, base.Extra)
			newMissing := extraDiff(missing, base.Missing)
			goneExtra := extraDiff(base.Extra, extra)
			goneMissing := extraDiff(base.Missing, missing)

			if len(newExtra) > 0 || len(newMissing) > 0 || len(goneExtra) > 0 || len(goneMissing) > 0 {
				t.Errorf("schema↔inline drift changed for %q (md: %s)\n"+
					"  NEW extra keys (inline, not in schema): %v\n"+
					"  NEW missing keys (schema, not inline):  %v\n"+
					"  baseline extra no longer present:       %v\n"+
					"  baseline missing no longer present:     %v\n"+
					"  -> reconcile the canonical source or update the baseline intentionally",
					art.name, art.mdPath, newExtra, newMissing, goneExtra, goneMissing)
			}
		})
	}
}

// TestSchemaInlineGenerated is the schema-sot Wave 2 byte-equal contract for
// GENERATED inline blocks (R-002, R-005). For each generated artifact it renders
// the canonical schema via renderSchemaExampleRegion and asserts the bytes
// committed between the artifact's generated-region markers are IDENTICAL. It
// also asserts both markers are present (R-005: absence must fail LOUD, never
// silently skip) and that the renderer is a fixed point under round-trip
// (R-002).
func TestSchemaInlineGenerated(t *testing.T) {
	skills := skillDir(t)
	schemas := schemaDir(t)

	any := false
	for _, art := range schemaArtifacts {
		if !art.generated {
			continue
		}
		any = true
		art := art
		t.Run(art.name, func(t *testing.T) {
			mdPath := filepath.Join(skills, filepath.FromSlash(art.mdPath))
			schemaPath := filepath.Join(schemas, art.schemaFile)

			schemaBytes, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Fatalf("%s: read schema %s: %v", art.name, schemaPath, err)
			}
			contentBytes, err := os.ReadFile(mdPath)
			if err != nil {
				t.Fatalf("%s: read md %s: %v", art.name, mdPath, err)
			}

			want, err := renderSchemaExampleRegion(schemaBytes, art.fenceLabel)
			if err != nil {
				t.Fatalf("%s: render schema example region: %v", art.name, err)
			}

			// R-005: BOTH markers must be present; absence fails loud here,
			// it does NOT silently skip the byte-equal contract.
			got, err := extractMarkerRegion(contentBytes, art.markerBegin, art.markerEnd)
			if err != nil {
				t.Fatalf("%s: generated-region markers missing/malformed in %s: %v\n"+
					"  expected begin %q and end %q to bound the generated block",
					art.name, art.mdPath, err, art.markerBegin, art.markerEnd)
			}

			if got != want {
				t.Errorf("generated inline block is STALE for %q (md: %s)\n"+
					"  the committed block diverges from the schema-rendered block.\n"+
					"  Re-run the schema-sot generator (or sync the .md) so committed == rendered.\n"+
					"--- COMMITTED (between markers) ---\n%s\n--- RENDERED (from %s) ---\n%s\n--- END ---",
					art.name, art.mdPath, got, art.schemaFile, want)
			}

			// R-002: golden round-trip — the rendered payload must be a fixed
			// point. Parse the rendered YAML payload and re-render from a
			// schema-shaped re-encode is not meaningful here; instead assert the
			// renderer is idempotent: rendering the same schema twice yields
			// identical bytes, and the rendered payload parses as valid YAML.
			want2, err := renderSchemaExampleRegion(schemaBytes, art.fenceLabel)
			if err != nil {
				t.Fatalf("%s: re-render for round-trip: %v", art.name, err)
			}
			if want != want2 {
				t.Errorf("%s: renderer is NOT deterministic — two renders differ", art.name)
			}
			payload := renderedPayloadBody(t, want)
			var parsed map[string]yaml.Node
			if err := yaml.Unmarshal([]byte(payload), &parsed); err != nil {
				t.Errorf("%s: rendered payload is not valid YAML: %v\n%s", art.name, err, payload)
			}
		})
	}

	if !any {
		t.Skip("no generated artifacts to verify")
	}
}

// renderedPayloadBody extracts the YAML body inside the ```yaml fence of a
// rendered marker-region string, so the round-trip test can re-parse it.
func renderedPayloadBody(t *testing.T, region string) string {
	t.Helper()
	m := fenceRe.FindStringSubmatch(region)
	if m == nil {
		t.Fatalf("rendered region has no ```yaml fence:\n%s", region)
	}
	return m[1]
}

// TestSchemaInlineCoverage is the R-006 coverage meta-test: every schema
// artifact must be covered by EXACTLY ONE path. A non-generated artifact lives
// under the set-parity drift path (TestSchemaInlineDrift) and MUST have a
// baseline entry; a generated artifact lives under the byte-equal path
// (TestSchemaInlineGenerated) and MUST have an EMPTY baseline entry. The two
// conditions are XOR, and every artifact must have a baseline key.
func TestSchemaInlineCoverage(t *testing.T) {
	for _, art := range schemaArtifacts {
		art := art
		t.Run(art.name, func(t *testing.T) {
			base, ok := schemaDriftBaseline[art.name]
			if !ok {
				t.Fatalf("%s: no baseline key — every artifact must have a baseline entry", art.name)
			}

			baselineEmpty := len(base.Extra) == 0 && len(base.Missing) == 0

			// XOR: generated IFF baseline empty IFF markers present.
			if art.generated {
				if !baselineEmpty {
					t.Errorf("%s: generated artifact MUST have an empty baseline (got Extra=%v Missing=%v)",
						art.name, base.Extra, base.Missing)
				}
				if art.markerBegin == "" || art.markerEnd == "" {
					t.Errorf("%s: generated artifact MUST declare both markers", art.name)
				}
			} else {
				if baselineEmpty {
					// A non-generated artifact with an empty baseline would be
					// silently uncovered by either path — disallowed by R-006.
					t.Errorf("%s: non-generated artifact has an EMPTY baseline — it would be covered by NO path; either generate it or document its drift baseline",
						art.name)
				}
				if art.markerBegin != "" || art.markerEnd != "" {
					t.Errorf("%s: non-generated artifact must NOT declare markers", art.name)
				}
			}
		})
	}
}
