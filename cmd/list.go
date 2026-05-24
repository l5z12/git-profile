package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// List returns `list` command
func List(cfg storage, v vcs) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List profiles",
		Long:    "Display the list of available profiles.",
		Example: "git-profile list",
		Run: func(cmd *cobra.Command, _ []string) {
			check(cmd, cfg, v)

			ui.Println(cmd, ui.InfoStyle, "Available profiles:")

			for _, name := range cfg.Names() {
				ui.Println(cmd, ui.SuccessStyle, "- %s:", name)

				profile, _ := cfg.Lookup(name)
				for _, entry := range profile {
					ui.Println(cmd, ui.DefaultStyle, "  %s: %s", entry.Key, entry.Value)
				}
			}
		},
	}
}
