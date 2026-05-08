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
		Long: `Render the per-profile statusline.

This is an internal command invoked by Claude Code itself. It reads a
JSON payload from stdin describing the current session and prints the
statusline. The wiring lives in settings.json under statusLine.command;
you do not normally call this directly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if f, ok := cmd.InOrStdin().(*os.File); ok {
				if fi, err := f.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
					return fmt.Errorf("statusline reads JSON from stdin (invoked by Claude Code, not interactively); pipe a payload or skip running it directly")
				}
			}
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
