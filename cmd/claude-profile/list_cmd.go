package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			profiles, err := profile.Discover(home)
			if err != nil {
				return err
			}
			active := profile.Active(home)
			if asJSON {
				data, _ := json.MarshalIndent(profiles, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "NAME\tPATH\tACTIVE")
			for _, p := range profiles {
				marker := ""
				if p.Name == active {
					marker = "●"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Name, p.Path, marker)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON")
	return cmd
}

func newCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			active := profile.Active(home)
			path := filepath.Join(home, ".claude")
			if active != "default" {
				path = filepath.Join(home, ".claude-"+active)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s)\n", active, path)
			return nil
		},
	}
}

func newWhichCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "which <name>",
		Short: "Print config dir path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			name := args[0]
			path := filepath.Join(home, ".claude")
			if name != "default" {
				path = filepath.Join(home, ".claude-"+name)
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func newTemplatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List embedded starter templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range templates.List() {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
}
