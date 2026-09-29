// Package ageagent lets programs that embed gopass as a library self-host
// the gopass age agent.
//
// The age backend auto-starts the agent by re-executing the host binary with
// the arguments `age agent start`. That only works if the host (a) registers
// a launcher on the context before initializing a store and (b) actually
// provides an `age agent start` command whose action is Serve. Hosts that do
// neither get graceful degradation: no autostart and a per-operation identity
// passphrase prompt instead of a fork bomb. See
// internal/backend/crypto/age.WithAgentLauncher.
//
// # Stability
//
// This package is best-effort stable. Additive changes (new exported symbols,
// new functional-option parameters) may appear in any release. Breaking changes
// (removal or signature change of an exported symbol, change of interface method
// sets or error semantics) require a [PKG-BREAK] entry in CHANGELOG.md and a
// minimum deprecation window of two minor releases or three months before the
// old symbol is removed. See docs/adr/A-12-pkg-api-stability.md for the full
// policy.
package ageagent

import (
	"context"
	"os/exec"

	"github.com/gopasspw/gopass/internal/backend/crypto/age"
)

// execCommand is indirection over os/exec.Command so tests can substitute a
// re-exec of the test binary and assert on the spawned child's environment
// without spawning a real agent.
var execCommand = exec.Command

// WithSelfLauncher registers an agent launcher that re-executes the current
// binary as `<argv0> age agent start` and detaches it. Use it if (and only if)
// the calling program provides an `age agent start` command backed by Serve;
// call it once during startup, before any store is initialized, so every
// age.New call site sees it.
func WithSelfLauncher(ctx context.Context) context.Context {
	return WithLauncher(ctx, launch)
}

// WithLauncher registers a custom agent launcher. Prefer this over
// WithSelfLauncher when the agent is started by a separate executable instead
// of a re-exec of the host binary.
func WithLauncher(ctx context.Context, l func(context.Context) error) context.Context {
	return age.WithAgentLauncher(ctx, l)
}
