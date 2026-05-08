package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpolate_ReplacesEnvVars(t *testing.T) {
	t.Setenv("FOO", "bar")
	t.Setenv("BAZ", "qux")

	p := &Profile{
		MCPServers: map[string]MCPServer{
			"s": {
				Command: "x",
				Env: map[string]string{
					"A": "${FOO}",
					"B": "literal",
					"C": "${BAZ}-suffix",
				},
			},
		},
	}

	err := Interpolate(p)
	require.NoError(t, err)
	require.Equal(t, "bar", p.MCPServers["s"].Env["A"])
	require.Equal(t, "literal", p.MCPServers["s"].Env["B"])
	require.Equal(t, "qux-suffix", p.MCPServers["s"].Env["C"])
}

func TestInterpolate_MissingVarFails(t *testing.T) {
	p := &Profile{
		MCPServers: map[string]MCPServer{
			"s": {Command: "x", Env: map[string]string{"A": "${MISSING_VAR_XYZ}"}},
		},
	}

	err := Interpolate(p)
	require.Error(t, err)
	require.Contains(t, err.Error(), "MISSING_VAR_XYZ")
	require.Contains(t, err.Error(), "mcp_servers.s.env.A")
}
