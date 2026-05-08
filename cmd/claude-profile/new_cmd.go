package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dmitriipyshinskii/claude-profile/internal/plugin"
	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/render"
	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/dmitriipyshinskii/claude-profile/internal/tcolor"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func okMark(w io.Writer) string   { return tcolor.Wrap(w, tcolor.Green, "✓") }
func failMark(w io.Writer) string { return tcolor.Wrap(w, tcolor.Red, "✗") }

func newNewCmd() *cobra.Command {
	var fromFile, fromTemplate string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new profile from yaml or template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !profile.SafeName(name) {
				return fmt.Errorf("invalid profile name %q: must start with a letter or digit and contain only [A-Za-z0-9_-]", name)
			}
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
			originalName := p.Metadata.Name
			p.Metadata.Name = name
			if p.Statusline.Label == originalName {
				p.Statusline.Label = name
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

			if err := os.MkdirAll(target, 0o700); err != nil {
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
				out := cmd.OutOrStdout()
				for _, src := range sources {
					if msg, bad := mfailed[src]; bad {
						fmt.Fprintf(out, "  %s %s (%s)\n      stderr: %s\n", failMark(out), labels[src], src, msg)
					} else {
						fmt.Fprintf(out, "  %s %s (%s)\n", okMark(out), labels[src], src)
					}
				}
				if len(mfails) > 0 {
					fmt.Fprintf(out, "%s %d of %d marketplace adds failed.\n", failMark(out), len(mfails), len(sources))
					return &PartialFailureError{Failed: len(mfails), Total: len(sources)}
				}
			}
			if len(p.Plugins) > 0 {
				out := cmd.OutOrStdout()
				fmt.Fprintf(out, "Installing %d plugins via claude CLI...\n", len(p.Plugins))
				failures, err := plugin.Install(target, p.Plugins)
				if err != nil {
					return err
				}
				for _, pl := range p.Plugins {
					failed := false
					for _, f := range failures {
						if f.Plugin == pl {
							fmt.Fprintf(out, "  %s %s\n      stderr: %s\n", failMark(out), pl, f.Stderr)
							failed = true
							break
						}
					}
					if !failed {
						fmt.Fprintf(out, "  %s %s\n", okMark(out), pl)
					}
				}
				if len(failures) > 0 {
					fmt.Fprintf(out, "%s %d of %d plugin installs failed.\n", failMark(out), len(failures), len(p.Plugins))
					return &PartialFailureError{Failed: len(failures), Total: len(p.Plugins)}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s Profile %q ready.\n", okMark(cmd.OutOrStdout()), name)
			printActivationHint(cmd.OutOrStdout(), name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fromFile, "file", "f", "", "Path to profile.yaml")
	cmd.Flags().StringVarP(&fromTemplate, "template", "t", "", "Embedded template name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print plan without writing")
	cmd.MarkFlagsMutuallyExclusive("file", "template")
	cmd.MarkFlagsOneRequired("file", "template")
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
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.Contains(trimmed, "claude-profile init") {
			return true
		}
	}
	return false
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
	if err := render.WriteAtomicSecret(filepath.Join(target, "settings.json"), settings); err != nil {
		return err
	}
	fmt.Fprintln(out, okMark(out), "Wrote settings.json")

	if mcp, ok, err := render.RenderMCP(p); err != nil {
		return err
	} else if ok {
		if err := render.WriteAtomicSecret(filepath.Join(target, ".mcp.json"), mcp); err != nil {
			return err
		}
		fmt.Fprintln(out, okMark(out), "Wrote .mcp.json")
	}

	if md, ok := render.RenderClaudeMD(p); ok {
		if err := render.WriteAtomic(filepath.Join(target, "CLAUDE.md"), md); err != nil {
			return err
		}
		fmt.Fprintln(out, okMark(out), "Wrote CLAUDE.md")
	}

	lock, err := render.RenderLock(p)
	if err != nil {
		return err
	}
	if err := render.WriteAtomicSecret(filepath.Join(target, "profile.lock.json"), lock); err != nil {
		return err
	}
	fmt.Fprintln(out, okMark(out), "Wrote profile.lock.json")
	return nil
}
