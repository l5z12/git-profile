package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

// Del returns `del` command
func Del(cfg storage) *cobra.Command {
	return delCommand(cfg, promptDeleteProfile)
}

//nolint:funlen
func delCommand(cfg storage, prompt func(storage, io.Writer) (string, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "del [profile] [key]",
		Aliases: []string{"rm"},
		Short:   "Delete an entry or a profile",
		Long: multiline(
			`Delete an entry from a profile or an entire profile.`,
			`Enter the "key" argument to remove the exact key from the profile.`,
		),
		Example: multiline(
			`git-profile del`,
			`git-profile del my-profile (delete the entire profile)`,
			`git-profile del my-profile user.name (deleting only a certain key)`,
		),
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) <= 2 { //nolint:mnd
				return nil
			}

			return fmt.Errorf("accepts at most 2 arguments")
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("config")

			if len(args) == 0 {
				profile, err := prompt(cfg, cmd.OutOrStdout())
				if err != nil {
					if err == huh.ErrUserAborted {
						cmd.Println("Interactive delete cancelled.")
						return
					}

					cmd.PrintErrln(err)
					os.Exit(1)
				}

				err = deleteProfile(cmd, cfg, filename, profile)
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}

				return
			}

			profile := args[0]

			if len(args) == 1 {
				err := deleteProfile(cmd, cfg, filename, profile)
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}
			} else {
				if ok := cfg.Delete(profile, args[1]); !ok {
					cmd.PrintErrln("There is no profile with given name")
					os.Exit(1)
				}

				err := cfg.Save(filename)
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}

				cmd.Printf("Successfully removed `%s` from `%s` profile.\n", args[1], profile)
			}
		},
	}
}

func promptDeleteProfile(cfg storage, _ io.Writer) (string, error) {
	if cfg.Len() == 0 {
		return "", fmt.Errorf("There are no available profiles")
	}

	names := cfg.Names()
	sort.Strings(names)

	var profile string

	err := huh.NewSelect[string]().
		Title("Select a profile to delete").
		Options(huh.NewOptions(names...)...).
		Value(&profile).
		WithKeyMap(huh.NewDefaultKeyMap()).
		Run()
	if err != nil {
		return "", err
	}

	return profile, nil
}

func deleteProfile(cmd *cobra.Command, cfg storage, filename string, profile string) error {
	if ok := cfg.DeleteProfile(profile); !ok {
		cmd.PrintErrln("There is no profile with given name")
		os.Exit(1)
	}

	err := cfg.Save(filename)
	if err != nil {
		return err
	}

	cmd.Printf("Successfully removed `%s` profile.\n", profile)

	return nil
}
