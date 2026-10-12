package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// asdtBlockStart and asdtBlockEnd wrap the ASDT agent config block inside
	// each assistant's global config file. ASDT owns only the text between
	// them; everything outside the markers belongs to the user and is never
	// rewritten, so updates stay idempotent and removal stays clean.
	asdtBlockStart = "<!-- asdt:agent-config -->"
	asdtBlockEnd   = "<!-- /asdt:agent-config -->"
)

// claudeAgentDir returns the directory where Claude Code's global config lives (~/.claude).
func claudeAgentDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

// openCodeConfigDir returns the OpenCode config directory, respecting XDG_CONFIG_HOME.
func openCodeConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "opencode")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "opencode")
}

// AgentConfigAdapters lists the global config file of every supported
// assistant. Both carry the same marker block; the shared upsert/remove logic
// below is the only code that edits them.
var AgentConfigAdapters = []AgentConfigAdapterDescriptor{
	{
		AssistantID: AssistantClaudeCode,
		ConfigPath:  func() string { return filepath.Join(claudeAgentDir(), "CLAUDE.md") },
	},
	{
		AssistantID: AssistantOpenCode,
		ConfigPath:  func() string { return filepath.Join(openCodeConfigDir(), "AGENTS.md") },
	},
}

// ConfigExists reports whether the assistant's global config file exists with
// any content — ASDT's block or the user's own. The TUI asks before editing an
// existing file.
func (a AgentConfigAdapterDescriptor) ConfigExists() bool {
	info, err := os.Stat(a.ConfigPath())
	return err == nil && info.Size() > 0
}

// HasBlock reports whether the assistant's global config file carries an ASDT
// marker (well-formed or not), so the dashboard can offer removal.
func (a AgentConfigAdapterDescriptor) HasBlock() bool {
	data, err := os.ReadFile(a.ConfigPath())
	if err != nil {
		return false
	}
	s := string(data)
	return len(markerLines(s, asdtBlockStart)) > 0 || len(markerLines(s, asdtBlockEnd)) > 0
}

// markerLines returns the byte offset of every line in content that holds
// marker and nothing else (a trailing "\r" is tolerated). A marker quoted inside
// a line of the user's prose is not a marker, so it never counts toward a block.
func markerLines(content, marker string) []int {
	var offsets []int
	for off := 0; off < len(content); {
		line := content[off:]
		next := len(content)
		if i := strings.IndexByte(line, '\n'); i >= 0 {
			line = line[:i]
			next = off + i + 1
		}
		if strings.TrimSuffix(line, "\r") == marker {
			offsets = append(offsets, off)
		}
		off = next
	}
	return offsets
}

// agentBlockSpan locates ASDT's block in content. Only a marker alone on its
// own line counts (see markerLines). found is false when neither marker is
// present. Any partial or malformed state (one marker without the other,
// either marker more than once, end before start) is an error: the caller
// writes nothing rather than guess which text is ASDT's.
// start is the index of the start marker; end is the index just past the end
// marker.
func agentBlockSpan(content string) (start, end int, found bool, err error) {
	starts := markerLines(content, asdtBlockStart)
	ends := markerLines(content, asdtBlockEnd)
	if len(starts) == 0 && len(ends) == 0 {
		return 0, 0, false, nil
	}
	if len(starts) != 1 || len(ends) != 1 {
		return 0, 0, false, fmt.Errorf("malformed ASDT block: %q appears %d time(s) and %q %d time(s) on a line of its own, want exactly 1 each — fix the markers by hand, then retry", asdtBlockStart, len(starts), asdtBlockEnd, len(ends))
	}
	start, endIdx := starts[0], ends[0]
	if endIdx < start {
		return 0, 0, false, fmt.Errorf("malformed ASDT block: %q precedes %q — fix the markers by hand, then retry", asdtBlockEnd, asdtBlockStart)
	}
	return start, endIdx + len(asdtBlockEnd), true, nil
}

// upsertAgentBlock returns content with ASDT's block set to rendered. An
// existing block is replaced in place; otherwise the block is added at the end,
// after one blank line, and closed by a newline. Text outside the markers is
// returned byte-identical.
func upsertAgentBlock(content, rendered string) (string, error) {
	block := asdtBlockStart + "\n" + strings.TrimRight(rendered, "\n") + "\n" + asdtBlockEnd

	start, end, found, err := agentBlockSpan(content)
	if err != nil {
		return "", err
	}
	if found {
		return content[:start] + block + content[end:], nil
	}

	switch {
	case content == "":
		return block + "\n", nil
	case strings.HasSuffix(content, "\n"):
		return content + "\n" + block + "\n", nil
	default:
		return content + "\n\n" + block + "\n", nil
	}
}

// removeAgentBlock returns content without ASDT's block, the newline that
// closes it, and the blank line that separates it from the user's text. found
// is false (content unchanged) when there is no block.
func removeAgentBlock(content string) (string, bool, error) {
	start, end, found, err := agentBlockSpan(content)
	if err != nil || !found {
		return content, false, err
	}

	before, after := content[:start], content[end:]
	after = strings.TrimPrefix(after, "\n") // the block's own line ending
	switch {
	case strings.HasSuffix(before, "\n\n"):
		before = before[:len(before)-1] // the blank separator line
	case before == "":
		after = strings.TrimPrefix(after, "\n") // block opened the file: drop the blank line after it
	}
	return before + after, true, nil
}

// writeAgentBlock upserts rendered into the file at path, creating it (and its
// directory) when absent. The write is atomic and symlink-safe (see
// atomicReplaceFile). Nothing is written when the markers are malformed or the
// file already holds exactly this block.
func writeAgentBlock(path, rendered string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	updated, err := upsertAgentBlock(string(existing), rendered)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if existing != nil && updated == string(existing) {
		return nil
	}
	return atomicReplaceFile(path, []byte(updated))
}

// RemoveAgentConfig removes ASDT's block from the global config file of each
// given assistant, leaving the user's text byte-identical to what it was
// before install, except that a missing final newline is added: install
// cannot record whether the newline before its block was the user's, so
// removal keeps it. The write is atomic and symlink-safe. The file itself is
// never deleted, even when nothing is left in it. An assistant whose file has
// no block, or that has no adapter, is reported as Skipped; malformed markers
// are an error and the file is left untouched. A removal also clears the
// persona recorded in the install metadata.
func RemoveAgentConfig(assistants []AssistantDescriptor) []AgentConfigResult {
	results := make([]AgentConfigResult, len(assistants))
	for i, a := range assistants {
		results[i] = AgentConfigResult{AssistantID: a.ID}
		adapter, ok := AgentConfigAdapterFor(a.ID)
		if !ok {
			results[i].Skipped = true
			continue
		}
		path := adapter.ConfigPath()
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			results[i].Skipped = true
			continue
		}
		if err != nil {
			results[i].Err = err
			continue
		}
		updated, found, err := removeAgentBlock(string(data))
		if err != nil {
			results[i].Err = fmt.Errorf("%s: %w", path, err)
			continue
		}
		if !found {
			results[i].Skipped = true
			continue
		}
		if err := atomicReplaceFile(path, []byte(updated)); err != nil {
			results[i].Err = err
			continue
		}
		results[i].Written = []string{path}

		// Best-effort: the persona no longer applies, so the dashboard stops
		// showing it.
		if meta, merr := ReadInstallMeta(a); merr == nil {
			meta.Persona, meta.Emojis = "", ""
			_ = WriteInstallMeta(a, meta)
		}
	}
	return results
}
