package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/tcolor"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var asJSON, all bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles",
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
			active := profile.Active(home)
			if asJSON {
				data, _ := json.MarshalIndent(profiles, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}
			out := cmd.OutOrStdout()
			var buf bytes.Buffer
			tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
			if all {
				fmt.Fprintln(tw, "NAME\tPATH\tMANAGED\tACTIVE")
			} else {
				fmt.Fprintln(tw, "NAME\tPATH\tACTIVE")
			}
			activeColor := ""
			for _, p := range profiles {
				marker := ""
				if p.Name == active {
					marker = "●"
					activeColor = p.Color
				}
				if all {
					managed := "no"
					if p.Managed {
						managed = "yes"
					}
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", p.Name, p.Path, managed, marker)
				} else {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Name, p.Path, marker)
				}
			}
			if err := tw.Flush(); err != nil {
				return err
			}
			rendered := buf.String()
			if tcolor.Enabled(out) {
				rendered = strings.Replace(rendered, "●", tcolor.WrapName(out, activeColor, "●"), 1)
			}
			fmt.Fprint(out, rendered)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON")
	cmd.Flags().BoolVar(&all, "all", false, "Include unmanaged ~/.claude-* directories")
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
			color := ""
			if profiles, err := profile.Discover(home); err == nil {
				for _, p := range profiles {
					if p.Name == active {
						color = p.Color
						break
					}
				}
			}
			out := cmd.OutOrStdout()
			dot := tcolor.WrapName(out, color, "●")
			fmt.Fprintf(out, "%s %s (%s)\n", dot, active, path)
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
