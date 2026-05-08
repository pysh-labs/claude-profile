package statusline

import (
	"fmt"
	"strings"
)

type Input struct {
	Cwd        string
	Home       string
	Branch     string
	ModelName  string
	ContextPct float64
}

type ProfileMeta struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Color string `json:"color"`
}

const (
	reset    = "\x1b[0m"
	dim      = "\x1b[2m"
	cyanBold = "\x1b[1;36m"
	yellow   = "\x1b[0;33m"
)

var colors = map[string]string{
	"red":     "\x1b[1;31m",
	"green":   "\x1b[1;32m",
	"yellow":  "\x1b[1;33m",
	"blue":    "\x1b[1;34m",
	"magenta": "\x1b[1;35m",
	"cyan":    "\x1b[1;36m",
	"white":   "\x1b[1;37m",
}

func Render(in Input, p ProfileMeta) string {
	var parts []string

	if p.Label != "" {
		c := colors[p.Color]
		parts = append(parts, fmt.Sprintf("%s[%s]%s", c, p.Label, reset))
	}

	short := in.Cwd
	if in.Home != "" && strings.HasPrefix(short, in.Home) {
		short = "~" + strings.TrimPrefix(short, in.Home)
	}
	if short != "" {
		parts = append(parts, fmt.Sprintf("%s%s%s", cyanBold, short, reset))
	}

	if in.Branch != "" {
		parts = append(parts, fmt.Sprintf("%s⎇ %s%s", yellow, in.Branch, reset))
	}

	if in.ModelName != "" {
		parts = append(parts, fmt.Sprintf("%s%s%s", dim, in.ModelName, reset))
	}

	if in.ContextPct > 0 {
		parts = append(parts, ctxSegment(in.ContextPct))
	}

	return strings.Join(parts, "  ")
}

func ctxSegment(pct float64) string {
	var c string
	switch {
	case pct >= 80:
		c = colors["red"]
	case pct >= 50:
		c = colors["yellow"]
	default:
		c = colors["green"]
	}
	return fmt.Sprintf("%sctx:%.0f%%%s", c, pct, reset)
}
