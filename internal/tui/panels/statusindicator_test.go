package panels_test

import (
	"strings"
	"testing"

	"github.com/vitualizz/asdt/internal/tui/panels"
)

func TestFocusBorderStyleFocusedRender(t *testing.T) {
	style := panels.FocusBorderStyle(true)
	rendered := style.Render("test")
	if !strings.Contains(rendered, "test") {
		t.Errorf("FocusBorderStyle(true) render lost content, got: %q", rendered)
	}
}

func TestFocusBorderStyleUnfocusedRender(t *testing.T) {
	style := panels.FocusBorderStyle(false)
	rendered := style.Render("test")
	if !strings.Contains(rendered, "test") {
		t.Errorf("FocusBorderStyle(false) render lost content, got: %q", rendered)
	}
}
