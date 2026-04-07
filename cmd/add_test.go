package cmd

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"charm.land/huh/v2"
	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestAdd(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := Add(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile", "key", "value"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(b.String()), "Successfully stored `key=value` to `profile` profile.")
}

func TestAddRejectsPartialArgs(t *testing.T) {
	is := is.New(t)

	cmd := Add(&storageMock{})
	cmd.SetArgs([]string{"profile"})

	err := cmd.Execute()

	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "accepts either no arguments for interactive mode or exactly 3 arguments"))
}

func TestAddInteractiveStoresOnlyChangedFilledValues(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: userNameKey, Value: "Jane Doe"},
				{Key: userEmailKey, Value: "old@example.com"},
			}, true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := addCommand(cfg, func(cfg storage, _ io.Reader, _ io.Writer) (addPromptResult, error) {
		return addPromptResult{
			Profile: "work",
			Values: map[string]string{
				userNameKey:       "Jane Doe",
				userEmailKey:      "new@example.com",
				userSigningKeyKey: "",
			},
		}, nil
	})
	cmd.SetOut(&b)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(len(cfg.StoreCalls()), 1)
	is.Equal(cfg.StoreCalls()[0].Profile, "work")
	is.Equal(cfg.StoreCalls()[0].Key, userEmailKey)
	is.Equal(cfg.StoreCalls()[0].Value, "new@example.com")
	is.Equal(len(cfg.SaveCalls()), 1)
	is.Equal(strings.TrimSpace(b.String()), "Successfully updated `work` profile.")
}

func TestAddInteractiveSkipsSaveWhenNothingChanged(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: userNameKey, Value: "Jane Doe"},
				{Key: userEmailKey, Value: "work@example.com"},
			}, true
		},
		SaveFunc: func(filename string) error {
			return errors.New("save should not be called")
		},
		StoreFunc: func(profile string, key string, value string) {
			t.Fatalf("Store should not be called")
		},
	}

	var b bytes.Buffer

	cmd := addCommand(cfg, func(cfg storage, _ io.Reader, _ io.Writer) (addPromptResult, error) {
		return addPromptResult{
			Profile: "work",
			Values: map[string]string{
				userNameKey:       "Jane Doe",
				userEmailKey:      "work@example.com",
				userSigningKeyKey: "",
			},
		}, nil
	})
	cmd.SetOut(&b)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(len(cfg.StoreCalls()), 0)
	is.Equal(len(cfg.SaveCalls()), 0)
	is.Equal(strings.TrimSpace(b.String()), "No profile changes to save.")
}

func TestAddInteractiveAbortDoesNotPrintSaveError(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		StoreFunc: func(profile string, key string, value string) {
			t.Fatalf("Store should not be called")
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

	cmd := addCommand(cfg, func(cfg storage, _ io.Reader, _ io.Writer) (addPromptResult, error) {
		return addPromptResult{}, huh.ErrUserAborted
	})
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(out.String()), "Interactive add cancelled.")
	is.Equal(errOut.String(), "")
}
