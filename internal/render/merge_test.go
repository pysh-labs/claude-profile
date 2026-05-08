package render

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeepMerge_OverridesWin(t *testing.T) {
	base := map[string]any{"a": 1, "b": "base"}
	over := map[string]any{"b": "over", "c": true}

	got := DeepMerge(base, over)
	require.Equal(t, 1, got["a"])
	require.Equal(t, "over", got["b"])
	require.Equal(t, true, got["c"])
}

func TestDeepMerge_NestedObjectsRecurse(t *testing.T) {
	base := map[string]any{
		"perms": map[string]any{"allow": []any{"x"}, "deny": []any{"y"}},
	}
	over := map[string]any{
		"perms": map[string]any{"allow": []any{"z"}},
	}

	got := DeepMerge(base, over)
	perms := got["perms"].(map[string]any)
	require.Equal(t, []any{"z"}, perms["allow"], "arrays replace, not concat")
	require.Equal(t, []any{"y"}, perms["deny"])
}
