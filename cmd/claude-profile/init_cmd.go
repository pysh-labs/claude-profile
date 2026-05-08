package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dmitriipyshinskii/claude-profile/internal/plugin"
	"github.com/dmitriipyshinskii/claude-profile/internal/render"
	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var fromFile, fromTemplate string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Bootstrap a profile from yaml or template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			data, err := loadSpec(fromFile, fromTemplate)
			if err != nil {
				return err
			}
			p, err := spec.Parse(data)
			if err != nil {
				return err
			}
			if err := spec.Interpolate(p); err != nil {
				return err
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			target := filepath.Join(home, ".claude-"+name)

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create %s with %d plugins\n", target, len(p.Plugins))
				return nil
			}

			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			if err := writeOwnedFiles(p, target, cmd.OutOrStdout()); err != nil {
				return err
			}
			if len(p.Plugins) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Installing %d plugins via claude CLI...\n", len(p.Plugins))
				failures, err := plugin.Install(target, p.Plugins)
				if err != nil {
					return err
				}
				for _, pl := range p.Plugins {
					failed := false
					for _, f := range failures {
						if f.Plugin == pl {
							fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %s\n      stderr: %s\n", pl, f.Stderr)
							failed = true
							break
						}
					}
					if !failed {
						fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s\n", pl)
					}
				}
				if len(failures) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "✗ %d of %d plugin installs failed.\n", len(failures), len(p.Plugins))
					return &PartialFailureError{Failed: len(failures), Total: len(p.Plugins)}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Profile %q ready.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fromFile, "file", "f", "", "Path to profile.yaml")
	cmd.Flags().StringVarP(&fromTemplate, "template", "t", "", "Embedded template name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print plan without writing")
	return cmd
}

func loadSpec(file, tmpl string) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	if tmpl != "" {
		return templates.Load(tmpl)
	}
	return nil, fmt.Errorf("provide either -f <file> or -t <template>")
}

func writeOwnedFiles(p *spec.Profile, target string, out io.Writer) error {
	settings, err := render.RenderSettings(p)
	if err != nil {
		return err
	}
	if err := render.WriteAtomic(filepath.Join(target, "settings.json"), settings); err != nil {
		return err
	}
	fmt.Fprintln(out, "✓ Wrote settings.json")

	if mcp, ok, err := render.RenderMCP(p); err != nil {
		return err
	} else if ok {
		if err := render.WriteAtomic(filepath.Join(target, ".mcp.json"), mcp); err != nil {
			return err
		}
		fmt.Fprintln(out, "✓ Wrote .mcp.json")
	}

	if md, ok := render.RenderClaudeMD(p); ok {
		if err := render.WriteAtomic(filepath.Join(target, "CLAUDE.md"), md); err != nil {
			return err
		}
		fmt.Fprintln(out, "✓ Wrote CLAUDE.md")
	}

	lock, err := render.RenderLock(p)
	if err != nil {
		return err
	}
	if err := render.WriteAtomic(filepath.Join(target, "profile.lock.json"), lock); err != nil {
		return err
	}
	fmt.Fprintln(out, "✓ Wrote profile.lock.json")
	return nil
}
