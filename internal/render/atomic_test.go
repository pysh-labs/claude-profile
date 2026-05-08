package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteAtomic_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")

	require.NoError(t, WriteAtomic(path, []byte("hello")))

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "hello", string(got))
}

func TestWriteAtomic_OverwriteLeavesNoTmp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")

	require.NoError(t, WriteAtomic(path, []byte("v1")))
	require.NoError(t, WriteAtomic(path, []byte("v2")))

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "v2", string(got))

	entries, _ := os.ReadDir(dir)
	require.Len(t, entries, 1, "should not leave .tmp behind")
}
