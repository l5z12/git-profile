package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

var (
	// DefaultStyle is the default style for general text.
	DefaultStyle = lipgloss.NewStyle()
	// InfoStyle is a bold blue style for informational messages.
	InfoStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Blue)
	// SuccessStyle is a bold green style for success messages.
	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Green)
	// ErrorStyle is a bold red style for error messages.
	ErrorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Red)
)

// Println renders text with the given style and prints it to the defined output, fallback to Stderr if not set.
func Println(cmd *cobra.Command, style lipgloss.Style, msg string, agrs ...any) {
	cmd.Println(renderInline(style, fmt.Sprintf(msg, agrs...)))
}

// PrintErrln renders text with the given style and prints it to stderr.
func PrintErrln(cmd *cobra.Command, style lipgloss.Style, msg string, agrs ...any) {
	cmd.PrintErrln(renderInline(style, fmt.Sprintf(msg, agrs...)))
}

func renderInline(style lipgloss.Style, value string) string {
	return strings.TrimRight(style.Render(value), "\n")
}
