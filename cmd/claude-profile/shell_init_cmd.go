package main

import (
	"fmt"
	"os"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/shell"
	"github.com/spf13/cobra"
)

func newShellInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "shell-init <shell>",
		Short:     "Print shell aliases for sourcing",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"zsh", "bash", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			profiles, err := profile.Discover(home)
			if err != nil {
				return err
			}
			out, err := shell.Generate(args[0], profiles)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
	}
}
