package cmd

import (
	"bytes"
	"io"
	"testing"

	"charm.land/huh/v2"
	"github.com/matryer/is"
)

func TestDel(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		NamesFunc: func() []string {
			return []string{"work"}
		},
		DeleteProfileFunc: func(profile string) bool {
			return true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Del(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully removed `profile` profile.")
}

func TestDelKey(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		DeleteFunc: func(profile string, value string) bool {
			return true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Del(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile", "user.name"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully removed `user.name` from `profile` profile.")
}

func TestDelInteractive(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		NamesFunc: func() []string {
			return []string{"work"}
		},
		DeleteProfileFunc: func(profile string) bool {
			return true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
	}

	var b bytes.Buffer

	cmd := delCommand(cfg, func(_ []string, _ io.Reader, _ io.Writer) (string, error) {
		return "work", nil
	})

	cmd.SetOut(&b)
	cmd.SetArgs(nil)
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(len(cfg.DeleteProfileCalls()), 1)
	is.Equal(cfg.DeleteProfileCalls()[0].Profile, "work")
	is.Equal(len(cfg.SaveCalls()), 1)
	is.Equal(trim(b.String()), "Successfully removed `work` profile.")
}

func TestDelInteractiveCancelled(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		NamesFunc: func() []string {
			return []string{"work"}
		},
		DeleteProfileFunc: func(profile string) bool {
			t.Fatalf("DeleteProfile should not be called")
			return false
		},
		SaveFunc: func(filename string) error {
			t.Fatalf("Save should not be called")
			return nil
		},
	}

	var (
		out    bytes.Buffer
		errOut bytes.Buffer
	)

	cmd := delCommand(cfg, func(_ []string, _ io.Reader, _ io.Writer) (string, error) {
		return "", huh.ErrUserAborted
	})

	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(nil)
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(errOut.String()), "Interactive delete cancelled.")
}
