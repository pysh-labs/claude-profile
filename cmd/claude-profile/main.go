package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCmd(v string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "claude-profile",
		Short:         "Profile-as-code for Claude Code",
		Version:       v,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetVersionTemplate("claude-profile {{.Version}}\n")
	cmd.AddCommand(
		newStatuslineCmd(),
		newInitCmd(),
		newNewCmd(),
		newListCmd(),
		newCurrentCmd(),
		newWhichCmd(),
		newTemplatesCmd(),
	)
	return cmd
}

func main() {
	err := newRootCmd(version).Execute()
	if err == nil {
		return
	}
	var pf *PartialFailureError
	if errors.As(err, &pf) {
		os.Exit(3)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
