package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderLock_Shape(t *testing.T) {
	p := &spec.Profile{
		Metadata:   spec.Metadata{Name: "personal"},
		Statusline: spec.Statusline{Label: "personal", Color: "green"},
	}
	data, err := RenderLock(p)
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "personal", got["name"])
	require.Equal(t, "personal", got["label"])
	require.Equal(t, "green", got["color"])
}
