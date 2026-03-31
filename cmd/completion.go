package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Completion returns `completion` command
func Completion(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `Generate shell completion script for Git Profile.

To load completions:

Bash:

  $ source <(git-profile completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ git-profile completion bash > /etc/bash_completion.d/git-profile
  # macOS:
  $ git-profile completion bash > $(brew --prefix)/etc/bash_completion.d/git-profile

Zsh:

  # Create completion directory and generate file:
  $ mkdir -p ~/.zsh/completions
  $ git-profile completion zsh > ~/.zsh/completions/_git-profile

  # Add to ~/.zshrc if not already present:
  $ echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
  $ echo 'autoload -U compinit && compinit' >> ~/.zshrc

  # Reload shell or run: source ~/.zshrc

Fish:

  $ git-profile completion fish | source

  # To load completions for each session, execute once:
  $ git-profile completion fish > ~/.config/fish/completions/git-profile.fish
`,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = rootCmd.GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				_ = rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				_ = rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}
