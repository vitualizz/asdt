package installer_test

import (
	"testing"

	"github.com/vitualizz/asdt/internal/installer"
)

func TestProviders_NonEmpty(t *testing.T) {
	if len(installer.Providers) == 0 {
		t.Fatal("Providers slice is empty")
	}
}

func TestProviders_EngramEntry(t *testing.T) {
	var found bool
	for _, p := range installer.Providers {
		if p.ID == installer.ProviderEngram {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("no provider with ID = ProviderEngram found in Providers")
	}
}

// TestProviders_BindEveryVerb guards the data every provider must carry for the
// installer to bind the prompts' memory interface: a config value for
// /asdt-init to write, a tool Name and Usage for every verb the prompts speak,
// and the Claude Code prefixes the executor agents are granted those tools
// under. A gap here would install a prompt with a verb it cannot call, or
// agents with no memory grant. The failure modes are TestBindMemory's cases.
func TestProviders_BindEveryVerb(t *testing.T) {
	for _, p := range installer.Providers {
		t.Run(p.ID, func(t *testing.T) {
			if err := p.Validate(); err != nil {
				t.Errorf("provider %q does not validate: %v", p.ID, err)
			}
		})
	}
}

// TestProviders_UniqueIdentity guards the keys the TUI and /asdt-init look
// providers up by: two providers sharing an ID, Name, or ConfigValue would
// collide in preflight rows, in the selection screen, or in the
// Prerequisites check that matches memory.provider against the binding.
func TestProviders_UniqueIdentity(t *testing.T) {
	fields := []struct {
		name string
		of   func(installer.ProviderDescriptor) string
	}{
		{"ID", func(p installer.ProviderDescriptor) string { return p.ID }},
		{"Name", func(p installer.ProviderDescriptor) string { return p.Name }},
		{"ConfigValue", func(p installer.ProviderDescriptor) string { return p.ConfigValue }},
	}
	for _, f := range fields {
		t.Run(f.name, func(t *testing.T) {
			seen := make(map[string]string, len(installer.Providers))
			for _, p := range installer.Providers {
				v := f.of(p)
				if prev, dup := seen[v]; dup {
					t.Errorf("providers %q and %q share %s %q", prev, p.ID, f.name, v)
				}
				seen[v] = p.ID
			}
		})
	}
}
