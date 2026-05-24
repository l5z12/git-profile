package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

const (
	currentProfileKey  = "current-profile.name"
	defaultProfileName = "default"
)

// Current returns `current` command
func Current(cfg storage, v vcs) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Aliases: []string{"c"},
		Short:   "Show the current profile",
		Long:    "Show the selected profile for the current repository.",
		Example: "git-profile current",
		PreRun: func(cmd *cobra.Command, _ []string) {
			check(cmd, cfg, v)
		},
		Run: func(cmd *cobra.Command, _ []string) {
			profile, err := v.Get(currentProfileKey)
			if len(profile) == 0 || err != nil {
				profile = defaultProfileName
			}

			ui.Println(cmd, ui.SuccessStyle, "Current profile: %s", profile)
		},
	}
}
