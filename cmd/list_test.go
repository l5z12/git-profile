package cmd

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
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
	is.Equal(ansi.Strip(b.String()), "Available profiles:\n\nhome\nuser.email: work@example.com\n")
}
