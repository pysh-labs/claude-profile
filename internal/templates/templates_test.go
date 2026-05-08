package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList_ReturnsAll(t *testing.T) {
	got := List()
	require.ElementsMatch(t, []string{"personal", "solo-dev", "enterprise"}, got)
}

func TestLoad_Personal(t *testing.T) {
	data, err := Load("personal")
	require.NoError(t, err)
	require.Contains(t, string(data), "name: personal")
}

func TestLoad_Unknown(t *testing.T) {
	_, err := Load("missing")
	require.Error(t, err)
}
