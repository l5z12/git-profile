package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	// TitleStyle is a bold cyan style for titles.
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Cyan)
	// NameStyle is a bold green style for names.
	NameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Green)
	// KeyStyle is a yellow style for keys.
	KeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Yellow)
	// ValueStyle is a blue style for values.
	ValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Blue)
)

// RenderInline renders a string with the given style and removes any trailing newlines.
func RenderInline(style lipgloss.Style, value string) string {
	return strings.TrimRight(style.Render(value), "\n")
}
