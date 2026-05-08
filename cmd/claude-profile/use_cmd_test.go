package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCmd_ErrorsWithoutShellIntegration(t *testing.T) {
	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"use", "personal"})

	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "shell integration")
}
