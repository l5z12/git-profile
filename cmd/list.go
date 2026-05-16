package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// List returns `list` command
func List(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List profiles",
		Long:    "Display the list of available profiles.",
		Example: "git-profile list",
		Run: func(cmd *cobra.Command, _ []string) {
			if cfg.Len() == 0 {
				cmd.Println(`There are no available profiles.`)
				cmd.Println(`To add a new profile, use the following examples:`)
				cmd.Println(`  git-profile add my-profile user.name "John Doe"`)
				cmd.Println(`  git-profile add my-profile user.email work@example.com`)

				return
			}

			cmd.Println(ui.RenderInline(ui.TitleStyle, "Available profiles:"))

			for _, name := range cfg.Names() {
				cmd.Println("\n" + ui.RenderInline(ui.NameStyle, name))

				profile, _ := cfg.Lookup(name)
				for _, entry := range profile {
					cmd.Printf("%s: %s\n", ui.RenderInline(ui.KeyStyle, entry.Key), ui.RenderInline(ui.ValueStyle, entry.Value))
				}
			}
		},
	}
}
