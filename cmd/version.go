package cmd

import (
	"github.com/spf13/cobra"
)

// Version returns `version` command
func Version(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print version information",
		Long:    "Print the version, commit hash, and build date.",
		Example: "git-profile version",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println("Git Profile")
			cmd.Println("Version:", c.Version)
			cmd.Println("Commit hash:", c.CommitHash)
			cmd.Println("Compiled on", c.CompileDate)
		},
	}
}
