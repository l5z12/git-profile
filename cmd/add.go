package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/config"
)

const (
	userNameKey       = "user.name"
	userEmailKey      = "user.email"
	userSigningKeyKey = "user.signingkey"
)

type addPromptResult struct {
	Profile string
	Values  map[string]string
}

// Add returns `add` command
func Add(cfg storage) *cobra.Command {
	return addCommand(cfg, runAddPrompt)
}

func addCommand(cfg storage, prompt func(storage, io.Reader, io.Writer) (addPromptResult, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "add [profile] [key] [value]",
		Aliases: []string{"set"},
		Short:   "Add an entry to a profile",
		Long:    "Adds an entry to a profile or updates an existing profile.",
		Example: multiline(
			`git-profile add`,
			`git-profile add my-profile user.email work@example.com`,
			`git-profile add my-profile user.name "John Doe"`,
			`git-profile add my-profile user.signingkey AAAAAAAA`,
		),
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 || len(args) == 3 {
				return nil
			}

			return fmt.Errorf("accepts either no arguments for interactive mode or exactly 3 arguments")
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("config")

			if len(args) == 3 { //nolint:mnd
				err := storeSingleEntry(cmd, cfg, filename, args[0], args[1], args[2])
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}

				return
			}

			err := storeInteractiveEntries(cmd, cfg, filename, prompt)
			if err != nil {
				if err == huh.ErrUserAborted {
					cmd.Println("Interactive add cancelled.")
					return
				}

				cmd.PrintErrln("Unable to save config file:", err)
				os.Exit(1)
			}
		},
	}
}

func storeSingleEntry(
	cmd *cobra.Command,
	cfg storage,
	filename string,
	profile string,
	key string,
	value string,
) error {
	cfg.Store(profile, key, value)

	err := cfg.Save(filename)
	if err != nil {
		return err
	}

	cmd.Printf("Successfully stored `%s=%s` to `%s` profile.\n", key, value, profile)

	return nil
}

func storeInteractiveEntries(
	cmd *cobra.Command,
	cfg storage,
	filename string,
	prompt func(storage, io.Reader, io.Writer) (addPromptResult, error),
) error {
	result, err := prompt(cfg, cmd.InOrStdin(), cmd.OutOrStdout())
	if err != nil {
		return err
	}

	currentValues := map[string]string{}
	if entries, ok := cfg.Lookup(result.Profile); ok {
		currentValues = entriesToMap(entries)
	}

	changed := 0

	for _, key := range []string{userNameKey, userEmailKey, userSigningKeyKey} {
		value := strings.TrimSpace(result.Values[key])
		if value == "" {
			continue
		}

		if currentValues[key] == value {
			continue
		}

		cfg.Store(result.Profile, key, value)

		changed++
	}

	if changed == 0 {
		cmd.Println("No profile changes to save.")
		return nil
	}

	err = cfg.Save(filename)
	if err != nil {
		return err
	}

	cmd.Printf("Successfully updated `%s` profile.\n", result.Profile)

	return nil
}

func runAddPrompt(cfg storage, in io.Reader, out io.Writer) (addPromptResult, error) {
	result := addPromptResult{
		Values: make(map[string]string),
	}

	err := huh.NewForm(
		huh.NewGroup(huh.
			NewInput().
			Title("Enter an existing profile name or a new one.").
			Value(&result.Profile).
			Validate(func(value string) error {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("profile name cannot be empty")
				}

				return nil
			}),
		),
	).WithKeyMap(huh.NewDefaultKeyMap()).WithInput(in).WithOutput(out).Run()
	if err != nil {
		return addPromptResult{}, err
	}

	result.Profile = strings.TrimSpace(result.Profile)

	var (
		name       string
		email      string
		signingKey string
	)

	if entries, ok := cfg.Lookup(result.Profile); ok {
		values := entriesToMap(entries)
		name = values[userNameKey]
		email = values[userEmailKey]
		signingKey = values[userSigningKeyKey]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Value for user.name").
				Value(&name),
			huh.NewInput().
				Title("Value for user.email").
				Value(&email),
			huh.NewInput().
				Title("Value for user.signingkey").
				Value(&signingKey),
		),
	).WithKeyMap(huh.NewDefaultKeyMap()).WithInput(in).WithOutput(out).Run()
	if err != nil {
		return addPromptResult{}, err
	}

	result.Values[userNameKey] = name
	result.Values[userEmailKey] = email
	result.Values[userSigningKeyKey] = signingKey

	return result, nil
}

func entriesToMap(entries []config.Entry) map[string]string {
	values := make(map[string]string, len(entries))

	for _, entry := range entries {
		values[entry.Key] = entry.Value
	}

	return values
}
