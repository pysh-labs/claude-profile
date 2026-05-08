package tcolor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnabled_BufferIsNotTerminal(t *testing.T) {
	var buf bytes.Buffer
	require.False(t, Enabled(&buf))
}

func TestEnabled_RespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	require.False(t, Enabled(nil))
}

func TestWrap_PassesThroughWhenDisabled(t *testing.T) {
	var buf bytes.Buffer
	require.Equal(t, "●", Wrap(&buf, Green, "●"))
}

func TestWrapName_UnknownColorFallsBackToGreen(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	require.Equal(t, "x", WrapName(nil, "puce", "x"))
}
