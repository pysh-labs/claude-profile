package main

import (
	"fmt"
	"os"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/shell"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "init <shell>",
		Short: "Activate shell integration (prints aliases for sourcing)",
		Long: `Activate shell integration. Prints two things:
  - per-profile aliases (claude-X) for one-shot invocations
  - a wrapper for "claude-profile use <name>" that exports
    CLAUDE_CONFIG_DIR for the rest of the shell session

Supported shells: zsh, bash, fish.

By default, only managed profiles (default + dirs with profile.lock.json)
get aliases — unmanaged ~/.claude-* directories belonging to other tools
are skipped to avoid clobbering their commands. Use --all to include them.`,
		Example: `  # activate in current session
  eval "$(claude-profile init zsh)"

  # persist in your shell config
  echo 'eval "$(claude-profile init zsh)"' >> ~/.zshrc`,
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"zsh", "bash", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			profiles, err := profile.Discover(home)
			if err != nil {
				return err
			}
			if !all {
				kept := profiles[:0]
				for _, p := range profiles {
					if p.Name == "default" || p.Managed {
						kept = append(kept, p)
					}
				}
				profiles = kept
			}
			out, err := shell.Generate(args[0], profiles)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "Include unmanaged ~/.claude-* directories")
	return cmd
}
