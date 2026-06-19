package installer

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// updateGolden regenerates the AGENTS.md golden fixtures from current output.
// Run: go test ./internal/installer/ -run TestRenderAgentConfig_GoldenByteIdentity -update
// The golden files were FIRST captured from pre-move output to lock byte-for-byte
// identity across the persona/template relocation (ADR #2115); do not regenerate
// them to paper over an unintended rendering change.
var updateGolden = flag.Bool("update", false, "regenerate AGENTS.md golden fixtures")

// TestRenderAgentConfig_GoldenByteIdentity is the primary guard that relocating
// the persona/template assets out of skill/ into internal/installer/assets/ did
// not alter the rendered AGENTS.md. It renders PersonaPresets[0] (Sky) with
// emoji on and off against the package-private assetsFS and asserts byte-for-byte
// equality with the frozen pre-move goldens.
func TestRenderAgentConfig_GoldenByteIdentity(t *testing.T) {
	cases := []struct {
		name      string
		useEmojis bool
		golden    string
	}{
		{name: "emoji_on", useEmojis: true, golden: "testdata/agents_sky_emoji_on.golden"},
		{name: "emoji_off", useEmojis: false, golden: "testdata/agents_sky_emoji_off.golden"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := renderAgentConfig(PersonaPresets[0], c.useEmojis, "en")
			if err != nil {
				t.Fatalf("renderAgentConfig: %v", err)
			}

			if *updateGolden {
				if err := os.MkdirAll(filepath.Dir(c.golden), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(c.golden, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(c.golden)
			if err != nil {
				t.Fatalf("read golden %q: %v", c.golden, err)
			}
			if string(want) != got {
				t.Errorf("rendered AGENTS.md differs from golden %q — the move was NOT byte-identical.\nlen(golden)=%d len(got)=%d", c.golden, len(want), len(got))
			}
		})
	}
}
