package render

import (
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderClaudeMD_PassThrough(t *testing.T) {
	p := &spec.Profile{ClaudeMD: "Branch prefix: feature/.\n"}
	data, ok := RenderClaudeMD(p)
	require.True(t, ok)
	require.Equal(t, "Branch prefix: feature/.\n", string(data))
}

func TestRenderClaudeMD_EmptyReturnsFalse(t *testing.T) {
	p := &spec.Profile{}
	_, ok := RenderClaudeMD(p)
	require.False(t, ok)
}
