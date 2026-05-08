package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildFakeClaude(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "claude")
	cmd := exec.Command("go", "build", "-o", bin, "testdata/fake_claude.go")
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())
	return bin
}

func TestInstall_AllSucceed(t *testing.T) {
	bin := buildFakeClaude(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))

	logFile := filepath.Join(t.TempDir(), "calls.log")
	t.Setenv("FAKE_CLAUDE_LOG", logFile)

	failed, err := Install("/tmp/fake-cfg", []string{
		"a@m1",
		"b@m2",
	})
	require.NoError(t, err)
	require.Empty(t, failed)

	data, _ := os.ReadFile(logFile)
	calls := strings.Split(strings.TrimSpace(string(data)), "\n")
	require.Len(t, calls, 2)
	require.Equal(t, "plugin install a@m1", calls[0])
	require.Equal(t, "plugin install b@m2", calls[1])
}

func TestInstall_PartialFailure(t *testing.T) {
	bin := buildFakeClaude(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_CLAUDE_LOG", filepath.Join(t.TempDir(), "calls.log"))

	failed, err := Install("/tmp/fake-cfg", []string{"a@m1", "FAIL", "b@m2"})
	require.NoError(t, err, "Install returns no error; failures returned in slice")
	require.Len(t, failed, 1)
	require.Equal(t, "FAIL", failed[0].Plugin)
	require.Contains(t, failed[0].Stderr, "simulated failure")
}
