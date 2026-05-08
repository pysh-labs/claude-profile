package statusline

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender_FullInputWithProfile(t *testing.T) {
	in := Input{
		Cwd:        "/Users/me/Projects/x",
		Home:       "/Users/me",
		Branch:     "main",
		ModelName:  "Opus 4.7",
		ContextPct: 42,
	}
	p := ProfileMeta{Name: "personal", Label: "personal", Color: "green"}

	got := Render(in, p)
	require.Contains(t, got, "[personal]")
	require.Contains(t, got, "~/Projects/x")
	require.Contains(t, got, "main")
	require.Contains(t, got, "Opus 4.7")
	require.Contains(t, got, "ctx:42%")
	require.True(t, strings.Contains(got, "\x1b[1;32m"), "green ANSI for prefix")
}

func TestRender_NoProfile_NoPrefix(t *testing.T) {
	in := Input{Cwd: "/x", Home: "/x"}
	got := Render(in, ProfileMeta{})
	require.NotContains(t, got, "]", "prefix bracket absent when no profile")
}

func TestRender_NoBranch_SkipBranch(t *testing.T) {
	in := Input{Cwd: "/x", Home: "/x", ModelName: "M"}
	got := Render(in, ProfileMeta{})
	require.NotContains(t, got, "⎇")
}

func TestRender_CtxColors(t *testing.T) {
	cases := []struct {
		pct  float64
		ansi string
	}{
		{30, "\x1b[1;32m"},
		{60, "\x1b[1;33m"},
		{85, "\x1b[1;31m"},
	}
	for _, tc := range cases {
		got := Render(Input{Cwd: "/x", Home: "/x", ContextPct: tc.pct}, ProfileMeta{})
		require.Contains(t, got, tc.ansi, "pct=%v", tc.pct)
	}
}

func TestRender_HomeSubstitution(t *testing.T) {
	got := Render(
		Input{Cwd: "/Users/me/Projects/x", Home: "/Users/me"},
		ProfileMeta{},
	)
	require.Contains(t, got, "~/Projects/x")
	require.NotContains(t, got, "/Users/me/Projects/x")
}
