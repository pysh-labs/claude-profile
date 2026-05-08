package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmitriipyshinskii/claude-profile/internal/statusline"
	"github.com/spf13/cobra"
)

func newStatuslineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "statusline",
		Short: "Render statusline (called by Claude Code)",
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := statusline.ParseStdin(cmd.InOrStdin())
			if err != nil {
				return err
			}
			configDir := os.Getenv("CLAUDE_CONFIG_DIR")
			if configDir == "" {
				home, _ := os.UserHomeDir()
				configDir = filepath.Join(home, ".claude")
			}
			meta := statusline.LoadProfileMeta(configDir)
			fmt.Fprint(cmd.OutOrStdout(), statusline.Render(in, meta))
			return nil
		},
	}
}
