package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderMCP_OneServer(t *testing.T) {
	p := &spec.Profile{
		MCPServers: map[string]spec.MCPServer{
			"atl": {Command: "docker", Args: []string{"run", "x"}, Env: map[string]string{"K": "v"}},
		},
	}
	data, ok, err := RenderMCP(p)
	require.NoError(t, err)
	require.True(t, ok)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))
	servers := got["mcpServers"].(map[string]any)
	atl := servers["atl"].(map[string]any)
	require.Equal(t, "docker", atl["command"])
}

func TestRenderMCP_NoServersReturnsFalse(t *testing.T) {
	p := &spec.Profile{}
	_, ok, err := RenderMCP(p)
	require.NoError(t, err)
	require.False(t, ok)
}
