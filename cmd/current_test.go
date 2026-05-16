package cmd

import (
	"bytes"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/ui"
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
	is.Equal(b.String(), ui.RenderInline(ui.TitleStyle, "Current profile is:")+" "+ui.RenderInline(ui.NameStyle, "home")+"\n")
}
