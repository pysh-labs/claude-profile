package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd_VersionFlag(t *testing.T) {
	cmd := newRootCmd("0.0.1-test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--version"})

	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "0.0.1-test")
}
