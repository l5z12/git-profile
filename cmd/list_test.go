package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/ui"
)

func TestList(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		NamesFunc: func() []string {
			return []string{"home"}
		},
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: "user.email", Value: "work@example.com"},
			}, true
		},
	}

	var b bytes.Buffer

	cmd := List(cfg)

	cmd.SetOut(&b)
	err := cmd.Execute()

	is.NoErr(err)
	output := b.String()
	is.True(strings.Contains(output, ui.RenderInline(ui.TitleStyle, "Available profiles:")))
	is.True(strings.Contains(output, ui.RenderInline(ui.NameStyle, "home")))
	is.True(strings.Contains(output, ui.RenderInline(ui.KeyStyle, "user.email")+": "+ui.RenderInline(ui.ValueStyle, "work@example.com")))
	is.Equal(ansi.Strip(output), `Available profiles:

home
user.email: work@example.com
`)
}
