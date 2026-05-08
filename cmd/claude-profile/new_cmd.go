package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dmitriipyshinskii/claude-profile/internal/plugin"
	"github.com/dmitriipyshinskii/claude-profile/internal/render"
	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	var fromFile, fromTemplate string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new profile from yaml or template",
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
			if len(p.Marketplaces) > 0 {
				names := make([]string, 0, len(p.Marketplaces))
				for n := range p.Marketplaces {
					names = append(names, n)
				}
				sort.Strings(names)
				sources := make([]string, 0, len(names))
				labels := make(map[string]string, len(names))
				for _, n := range names {
					src := p.Marketplaces[n].Repo
					sources = append(sources, src)
					labels[src] = n
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Adding %d marketplace(s) via claude CLI...\n", len(sources))
				mfails, err := plugin.AddMarketplaces(target, sources)
				if err != nil {
					return err
				}
				mfailed := map[string]string{}
				for _, f := range mfails {
					mfailed[f.Source] = f.Stderr
				}
				for _, src := range sources {
					if msg, bad := mfailed[src]; bad {
						fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %s (%s)\n      stderr: %s\n", labels[src], src, msg)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s (%s)\n", labels[src], src)
					}
				}
				if len(mfails) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "✗ %d of %d marketplace adds failed.\n", len(mfails), len(sources))
					return &PartialFailureError{Failed: len(mfails), Total: len(sources)}
				}
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
			printActivationHint(cmd.OutOrStdout(), name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fromFile, "file", "f", "", "Path to profile.yaml")
	cmd.Flags().StringVarP(&fromTemplate, "template", "t", "", "Embedded template name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print plan without writing")
	return cmd
}

func printActivationHint(out io.Writer, profileName string) {
	shellName, rc := detectShellAndRC()
	alias := "claude-" + profileName

	if rcAlreadyActivated(rc) {
		fmt.Fprintf(out, "\nLaunch with: %s\n", alias)
		fmt.Fprintf(out, "(open a new shell or run `eval \"$(claude-profile init %s)\"` to refresh aliases)\n", shellName)
		return
	}

	loadCmd := fmt.Sprintf(`eval "$(claude-profile init %s)"`, shellName)
	if shellName == "fish" {
		loadCmd = "claude-profile init fish | source"
	}
	persistCmd := fmt.Sprintf(`echo '%s' >> %s`, loadCmd, rc)

	fmt.Fprintf(out, "\nLaunch with: %s\n", alias)
	fmt.Fprintf(out, "Activate aliases in this shell:\n  %s\n", loadCmd)
	fmt.Fprintf(out, "Persist for new shells:\n  %s\n", persistCmd)
}

func detectShellAndRC() (string, string) {
	home, _ := os.UserHomeDir()
	base := filepath.Base(os.Getenv("SHELL"))
	switch base {
	case "bash":
		return "bash", filepath.Join(home, ".bashrc")
	case "fish":
		return "fish", filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return "zsh", filepath.Join(home, ".zshrc")
	}
}

func rcAlreadyActivated(rc string) bool {
	data, err := os.ReadFile(rc)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "claude-profile init")
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
