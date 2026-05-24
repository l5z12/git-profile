package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// Del returns `del` command
func Del(cfg storage) *cobra.Command {
	return delCommand(cfg, ui.SelectProfile)
}

//nolint:funlen
func delCommand(
	cfg storage,
	prompt func([]string, io.Reader, io.Writer) (string, error),
) *cobra.Command {
	return &cobra.Command{
		Use:     "del [profile] [key]",
		Aliases: []string{"rm"},
		Short:   "Delete a profile or a profile entry",
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

			switch len(args) {
			case 0:
				profile, err := prompt(cfg.Names(), cmd.InOrStdin(), cmd.OutOrStdout())
				if err != nil {
					if ui.IsAborted(err) {
						ui.PrintErrln(cmd, ui.ErrorStyle, "Interactive delete cancelled.")
						return
					}

					ui.PrintErrln(cmd, ui.ErrorStyle, "%s", err)
					os.Exit(1)
				}

				err = profileDelete(cmd, cfg, filename, profile)
				if err != nil {
					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
					os.Exit(1)
				}

				return
			case 1:
				profile := args[0]

				err := profileDelete(cmd, cfg, filename, profile)
				if err != nil {
					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
					os.Exit(1)
				}
			case 2:
				profile := args[0]
				if ok := cfg.Delete(profile, args[1]); !ok {
					ui.PrintErrln(cmd, ui.ErrorStyle, "There is no profile with given name")
					os.Exit(1)
				}

				err := cfg.Save(filename)
				if err != nil {
					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
					os.Exit(1)
				}

				ui.Println(cmd, ui.SuccessStyle, "Successfully removed `%s` from `%s` profile.\n", args[1], profile)
			}
		},
	}
}

func profileDelete(cmd *cobra.Command, cfg storage, filename string, profile string) error {
	if ok := cfg.DeleteProfile(profile); !ok {
		ui.PrintErrln(cmd, ui.ErrorStyle, "There is no profile with given name")
		os.Exit(1)
	}

	err := cfg.Save(filename)
	if err != nil {
		return err
	}

	ui.Println(cmd, ui.SuccessStyle, "Successfully removed `%s` profile.\n", profile)

	return nil
}
