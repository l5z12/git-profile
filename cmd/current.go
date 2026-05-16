package cmd

import (
	"os"

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
		Run: func(cmd *cobra.Command, _ []string) {
			if cfg.Len() == 0 || !v.IsRepository() {
				os.Exit(1)
			}

			profile, err := v.Get(currentProfileKey)
			if len(profile) == 0 || err != nil {
				profile = defaultProfileName
			}

			cmd.Printf("%s %s\n", ui.RenderInline(ui.TitleStyle, "Current profile is:"), ui.RenderInline(ui.NameStyle, profile))
		},
	}
}
