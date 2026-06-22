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

// TestRenderAgentConfig_GoldenByteIdentity is the primary guard that the
// persona/template assets render to a stable AGENTS.md. It renders every
// PersonaPreset with emoji on and off against the package-private assetsFS and
// asserts byte-for-byte equality with the frozen goldens. The 5-preset ×
// {emoji_on,off} matrix byte-guards all persona demotions/dedups (US-108);
// goldens prove byte-stability, not content correctness — a human eyeballs the
// -update diff before the suite is trusted.
func TestRenderAgentConfig_GoldenByteIdentity(t *testing.T) {
	cases := []struct {
		name      string
		preset    PersonaPreset
		useEmojis bool
		golden    string
	}{
		{name: "sky_emoji_on", preset: PersonaPresets[0], useEmojis: true, golden: "testdata/agents_sky_emoji_on.golden"},
		{name: "sky_emoji_off", preset: PersonaPresets[0], useEmojis: false, golden: "testdata/agents_sky_emoji_off.golden"},
		{name: "toffy_emoji_on", preset: PersonaPresets[1], useEmojis: true, golden: "testdata/agents_toffy_emoji_on.golden"},
		{name: "toffy_emoji_off", preset: PersonaPresets[1], useEmojis: false, golden: "testdata/agents_toffy_emoji_off.golden"},
		{name: "atreus_emoji_on", preset: PersonaPresets[2], useEmojis: true, golden: "testdata/agents_atreus_emoji_on.golden"},
		{name: "atreus_emoji_off", preset: PersonaPresets[2], useEmojis: false, golden: "testdata/agents_atreus_emoji_off.golden"},
		{name: "babi_emoji_on", preset: PersonaPresets[3], useEmojis: true, golden: "testdata/agents_babi_emoji_on.golden"},
		{name: "babi_emoji_off", preset: PersonaPresets[3], useEmojis: false, golden: "testdata/agents_babi_emoji_off.golden"},
		{name: "lee_emoji_on", preset: PersonaPresets[4], useEmojis: true, golden: "testdata/agents_lee_emoji_on.golden"},
		{name: "lee_emoji_off", preset: PersonaPresets[4], useEmojis: false, golden: "testdata/agents_lee_emoji_off.golden"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := renderAgentConfig(c.preset, c.useEmojis, "en")
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
