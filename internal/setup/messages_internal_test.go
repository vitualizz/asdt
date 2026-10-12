package setup

import (
	"fmt"
	"testing"

	"github.com/vitualizz/asdt/internal/i18n"
	"github.com/vitualizz/asdt/internal/installer"
)

// TestNewerAvailable exercises the unexported newerAvailable comparison
// helper directly (white-box) since it must remain unexported per design.
func TestNewerAvailable(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{name: "dev build", current: "dev", latest: "v9.9.9", want: false},
		{name: "empty current", current: "", latest: "v1.2.3", want: false},
		{name: "equal normalized", current: "0.3.0", latest: "v0.3.0", want: false},
		{name: "newer", current: "v0.2.0", latest: "v0.3.0", want: true},
		{name: "empty latest", current: "v0.2.0", latest: "", want: false},
		{name: "whitespace tolerance", current: " v0.2.0 ", latest: "v0.3.0", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newerAvailable(tt.current, tt.latest)
			if got != tt.want {
				t.Errorf("newerAvailable(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}

// TestMemoryMetaName covers the dashboard's Memory value: a known provider's
// Name, the localized "none" for a "no memory" install, and the raw value for
// an ID no known provider has.
func TestMemoryMetaName(t *testing.T) {
	s := i18n.English.Dashboard
	for _, tt := range []struct{ value, want string }{
		{installer.ProviderEngram, "Engram"},
		{installer.MemoryNone, s.ValueNoMemory},
		{"retired", "retired"},
	} {
		if got := memoryMetaName(tt.value, s); got != tt.want {
			t.Errorf("memoryMetaName(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

// TestCodeIntelMetaName mirrors TestMemoryMetaName for the dashboard's
// code-intelligence part.
func TestCodeIntelMetaName(t *testing.T) {
	s := i18n.English.Dashboard
	for _, tt := range []struct{ value, want string }{
		{installer.ProviderCodegraph, "Codegraph"},
		{installer.CodeIntelNone, s.ValueNoCodeIntel},
		{"retired", "retired"},
	} {
		if got := codeIntelMetaName(tt.value, s); got != tt.want {
			t.Errorf("codeIntelMetaName(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

// TestProviderProbeCmd_LocalizedNotFound asserts the not-found preflight
// detail of both provider kinds comes from the catalog, so it follows the
// chosen language.
func TestProviderProbeCmd_LocalizedNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	detect := installer.ProviderDetect{Binary: "asdt-no-such-binary", InstallURL: "https://example.invalid/ghost"}
	for _, code := range []string{"en", "es", "es-ES"} {
		s := i18n.ForCode(code).Installer
		for kind, notFound := range map[string]string{"memory": s.DetailProviderNotFound, "code intelligence": s.DetailCodeIntelNotFound} {
			msg, ok := providerProbeCmd("Ghost", detect, notFound)().(EnvironmentCheckProgressMsg)
			if !ok {
				t.Fatalf("%s %s: providerProbeCmd returned %T", code, kind, msg)
			}
			want := fmt.Sprintf(notFound, detect.InstallURL)
			if msg.RowLabel != "Ghost" || msg.Detail != want || !msg.SoftWarn {
				t.Errorf("%s %s: row %q detail %q (softWarn %v), want Ghost %q (soft warning)", code, kind, msg.RowLabel, msg.Detail, msg.SoftWarn, want)
			}
		}
	}
}
