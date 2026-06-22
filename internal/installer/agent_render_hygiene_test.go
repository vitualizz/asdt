package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderAgentConfig_StripsHTMLComments is the AC1 guard: renderAgentConfig
// must emit no HTML comments. The template carries a source-only H2-ownership
// comment that stripHTMLComments removes before the body is wrapped for write.
// Covers every preset × {emoji on, off} so no persona reintroduces a comment.
func TestRenderAgentConfig_StripsHTMLComments(t *testing.T) {
	for _, preset := range PersonaPresets {
		for _, useEmojis := range []bool{true, false} {
			name := preset.ID
			if useEmojis {
				name += "_emoji_on"
			} else {
				name += "_emoji_off"
			}
			t.Run(name, func(t *testing.T) {
				got, err := renderAgentConfig(preset, useEmojis, "en")
				if err != nil {
					t.Fatalf("renderAgentConfig: %v", err)
				}
				if strings.Contains(got, "<!--") {
					t.Errorf("rendered config still contains an HTML comment marker '<!--':\n%s", got)
				}
			})
		}
	}
}

// TestWriteClaudeAgentConfig_PreservesAsdtMarkers is the AC2 end-to-end guard:
// the asdt block markers are themselves HTML comments, injected on the write
// path AFTER renderAgentConfig strips comments. Because the strip runs only
// inside renderAgentConfig (pre-wrap), the written CLAUDE.md must still carry
// BOTH markers — proving the strip never reaches the write path.
func TestWriteClaudeAgentConfig_PreservesAsdtMarkers(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	rendered, err := renderAgentConfig(PersonaPresets[0], true, "en")
	if err != nil {
		t.Fatalf("renderAgentConfig: %v", err)
	}

	if _, err := writeClaudeAgentConfig(rendered, AgentModeOverwrite); err != nil {
		t.Fatalf("writeClaudeAgentConfig: %v", err)
	}

	claudePath := filepath.Join(tmpHome, ".claude", "CLAUDE.md")
	data, readErr := os.ReadFile(claudePath)
	if readErr != nil {
		t.Fatalf("CLAUDE.md not found: %v", readErr)
	}
	content := string(data)
	if !strings.Contains(content, asdtBlockStart) {
		t.Errorf("written CLAUDE.md missing asdtBlockStart marker, got:\n%s", content)
	}
	if !strings.Contains(content, asdtBlockEnd) {
		t.Errorf("written CLAUDE.md missing asdtBlockEnd marker, got:\n%s", content)
	}
}
