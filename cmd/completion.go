package cmd

import (
	"os"

	"github.com/dotpilot/utils"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell auto-completion scripts for dotpilot.

To load completions:

Bash:
  $ source <(dotpilot completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ dotpilot completion bash > /etc/bash_completion.d/dotpilot
  # macOS:
  $ dotpilot completion bash > /usr/local/etc/bash_completion.d/dotpilot

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ dotpilot completion zsh > "${fpath[1]}/_dotpilot"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ dotpilot completion fish > ~/.config/fish/completions/dotpilot.fish

PowerShell:
  PS> dotpilot completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> dotpilot completion powershell > dotpilot.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}

		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to generate completion script")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}