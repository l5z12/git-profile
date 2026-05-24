package cmd

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/matryer/is"
)

func TestCurrent(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
	}

	vcs := &vcsMock{
		IsRepositoryFunc: func() bool {
			return true
		},
		GetFunc: func(key string) (string, error) {
			return "home", nil
		},
	}

	var b bytes.Buffer

	cmd := Current(cfg, vcs)

	cmd.SetOut(&b)
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(ansi.Strip(b.String()), "Current profile: home\n")
}
