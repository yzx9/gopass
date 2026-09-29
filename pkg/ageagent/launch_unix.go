//go:build !windows

package ageagent

import (
	"context"
	"os"
	"syscall"

	"github.com/gopasspw/gopass/internal/backend/crypto/age"
)

func launch(_ context.Context) error {
	// Prefer the resolved executable path over os.Args[0]: a host started by
	// PATH lookup (e.g. a browser wrapper running a bare binary name) must
	// re-exec itself, not whatever shadows that name in PATH at spawn time.
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}

	cmd := execCommand(exe, "age", "agent", "start")
	cmd.Env = append(os.Environ(), age.SpawnGuardEnv+"=1")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	waitChild(cmd)

	return nil
}
