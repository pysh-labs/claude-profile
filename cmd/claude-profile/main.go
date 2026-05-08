package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCmd(v string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claude-profile",
		Short:   "Profile-as-code for Claude Code",
		Version: v,
	}
	cmd.SetVersionTemplate("claude-profile {{.Version}}\n")
	return cmd
}

func main() {
	if err := newRootCmd(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
