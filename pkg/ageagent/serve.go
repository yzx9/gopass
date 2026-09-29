package ageagent

import (
	"context"
	"fmt"

	"github.com/gopasspw/gopass/internal/backend/crypto/age/agent"
	"github.com/gopasspw/gopass/internal/out"
)

// Serve creates the age agent and runs it in the foreground, blocking until it
// is shut down by SIGINT/SIGTERM or a client quit. It announces itself through
// the ctx-configured output, mirroring the gopass CLI action; when spawned by
// WithSelfLauncher the child's stdout is discarded, so this cannot pollute
// e.g. a native-messaging channel. Embedders must expose it under an
// `age agent start` command; see the package documentation.
func Serve(ctx context.Context) error {
	out.Printf(ctx, "Starting age agent ...")

	ag, err := agent.New()
	if err != nil {
		return fmt.Errorf("failed to create age agent: %w", err)
	}

	return ag.Run(ctx)
}
