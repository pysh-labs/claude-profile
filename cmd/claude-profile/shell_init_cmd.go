package main

import (
	"fmt"
	"os"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/shell"
	"github.com/spf13/cobra"
)

func newShellInitCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "shell-init <shell>",
		Short: "Print shell aliases for sourcing",
		Long: `Print shell aliases for switching profiles.

Supported shells: zsh, bash, fish.

By default, only managed profiles (default + dirs with profile.lock.json)
get aliases — unmanaged ~/.claude-* directories belonging to other tools
are skipped to avoid clobbering their commands. Use --all to include them.`,
		Example: `  # one-shot evaluation
  eval "$(claude-profile shell-init zsh)"

  # persist in your shell config
  claude-profile shell-init zsh >> ~/.zshrc`,
		Args:      cobra.ExactArgs(1),
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
