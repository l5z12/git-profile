package cmd

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func trim(s string) string {
	return strings.Trim(strings.TrimSpace(ansi.Strip(s)), "\n")
}
