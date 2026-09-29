//go:build !windows

package ageagent

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/gopasspw/gopass/internal/backend/crypto/age"
)

// TestLaunchStampsSpawnGuard proves launch() stamps age.SpawnGuardEnv on the
// spawned child and re-execs the resolved executable path. It overrides
// execCommand to re-exec the test binary as TestLaunchHelper, which exits
// non-zero unless it inherited SpawnGuardEnv=1. This locks the
// defense-in-depth layer: if the cmd.Env line in launch() is ever dropped,
// the helper fails and this test catches it.
func TestLaunchStampsSpawnGuard(t *testing.T) {
	// Carried into the child via launch()'s os.Environ(); tells the helper to run.
	t.Setenv("GP_AGENTLAUNCH_TEST_HELPER", "1")

	var spawned *exec.Cmd
	var spawnedName string
	orig := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		spawnedName = name
		spawned = exec.Command(os.Args[0], append([]string{"-test.run=TestLaunchHelper", "--", name}, args...)...)

		return spawned
	}
	t.Cleanup(func() { execCommand = orig })

	// The test waits on the child itself to observe its exit status, so take
	// the reaping responsibility away from launch().
	origWait := waitChild
	waitChild = func(*exec.Cmd) {}
	t.Cleanup(func() { waitChild = origWait })

	if err := launch(context.Background()); err != nil {
		t.Fatalf("launch returned error: %v", err)
	}
	if spawned == nil {
		t.Fatal("launch did not build a command")
	}

	// launch must re-exec the resolved executable path, not a bare os.Args[0]
	// that a PATH lookup could shadow.
	wantExe, err := os.Executable()
	if err != nil {
		wantExe = os.Args[0]
	}
	if spawnedName != wantExe {
		t.Fatalf("launch re-exec'd %q, want the resolved executable %q", spawnedName, wantExe)
	}

	if err := spawned.Wait(); err != nil {
		t.Fatalf("spawned child did not see %s=1 (launch must stamp it): %v", age.SpawnGuardEnv, err)
	}
}

// TestLaunchHelper is the child side of TestLaunchStampsSpawnGuard. It runs only
// when re-exec'd with GP_AGENTLAUNCH_TEST_HELPER=1, and fails unless it inherited
// SpawnGuardEnv=1 from launch()'s cmd.Env. Skips during normal test runs.
func TestLaunchHelper(t *testing.T) {
	if os.Getenv("GP_AGENTLAUNCH_TEST_HELPER") != "1" {
		t.Skip()
	}
	if os.Getenv(age.SpawnGuardEnv) != "1" {
		t.Fatalf("%s not set on spawned child", age.SpawnGuardEnv)
	}
}
