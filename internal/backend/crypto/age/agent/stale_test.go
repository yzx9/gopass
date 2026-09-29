package agent

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsStaleSocket(t *testing.T) {
	dir := t.TempDir()

	// A missing file is not stale (nothing to clean up).
	require.False(t, isStaleSocket(filepath.Join(dir, "missing.sock")))

	// A live listener is not stale.
	live := filepath.Join(dir, "live.sock")
	ln, err := net.Listen("unix", live)
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()
	require.False(t, isStaleSocket(live))

	// A leftover file with no listener behind it is stale.
	stale := filepath.Join(dir, "stale.sock")
	require.NoError(t, os.WriteFile(stale, nil, 0o600))
	require.True(t, isStaleSocket(stale))
}
