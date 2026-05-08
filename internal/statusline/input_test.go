package statusline

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStdin_AllFields(t *testing.T) {
	json := `{
		"workspace": {"current_dir": "/x"},
		"model": {"display_name": "Opus 4.7"},
		"context_window": {"used_percentage": 42.5}
	}`
	in, err := ParseStdin(strings.NewReader(json))
	require.NoError(t, err)
	require.Equal(t, "/x", in.Cwd)
	require.Equal(t, "Opus 4.7", in.ModelName)
	require.InDelta(t, 42.5, in.ContextPct, 0.01)
}

func TestParseStdin_FallbackCwd(t *testing.T) {
	json := `{"cwd": "/legacy"}`
	in, err := ParseStdin(strings.NewReader(json))
	require.NoError(t, err)
	require.Equal(t, "/legacy", in.Cwd)
}

func TestParseStdin_InvalidJSON(t *testing.T) {
	_, err := ParseStdin(strings.NewReader("not json"))
	require.Error(t, err)
}
