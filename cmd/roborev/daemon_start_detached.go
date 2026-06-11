package main

import (
	"context"

	kitdaemon "go.kenn.io/kit/daemon"
)

// startDetachedDaemon launches the daemon as a detached background process via
// go.kenn.io/kit/daemon.StartDetached. On Windows the child runs on its own
// hidden console (CREATE_NO_WINDOW) so neither it nor its console-subsystem
// descendants (git, agent CLIs) open visible console windows; on other
// platforms it detaches from the caller's process group.
func startDetachedDaemon(ctx context.Context, opts detachedDaemonOptions) error {
	return kitdaemon.StartDetached(ctx, kitdaemon.StartDetachedOptions{
		Executable:      opts.Executable,
		Args:            opts.Args,
		Env:             opts.Env,
		Stdout:          opts.Stdout,
		Stderr:          opts.Stderr,
		RefuseEphemeral: opts.RefuseEphemeral,
		AfterStart:      opts.AfterStart,
	})
}
