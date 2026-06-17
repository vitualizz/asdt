package panels

import "github.com/charmbracelet/lipgloss"

// FocusBorderStyle returns a border style with per-side colors when focused.
func FocusBorderStyle(focused bool) lipgloss.Style {
	if focused {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTopForeground(ColorPrimary).
			BorderBottomForeground(ColorInactive).
			BorderLeftForeground(ColorSecondary).
			BorderRightForeground(ColorInactive)
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorInactive)
}
