//go:build windows

package main

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CREATE_NO_WINDOW runs the daemon on a hidden console; DETACHED_PROCESS gives
// it no console at all, which makes every git/agent child allocate its own
// visible console window. These two flags are mutually exclusive, so the
// detached daemon must use CREATE_NO_WINDOW and must not use DETACHED_PROCESS.
const (
	createNoWindowFlag  = 0x08000000
	detachedProcessFlag = 0x00000008
)

func TestDetachedDaemonRunsOnHiddenConsole(t *testing.T) {
	assert := assert.New(t)

	comSpec := os.Getenv("ComSpec")
	if comSpec == "" {
		comSpec = "cmd.exe"
	}

	var launched *exec.Cmd
	err := startDetachedDaemon(context.Background(), detachedDaemonOptions{
		Executable: comSpec,
		Args:       []string{"/c", "exit", "0"},
		AfterStart: func(cmd *exec.Cmd) { launched = cmd },
	})
	require.NoError(t, err)
	require.NotNil(t, launched, "AfterStart must run with the launched command")
	require.NotNil(t, launched.SysProcAttr, "detached daemon must set process attributes")

	flags := launched.SysProcAttr.CreationFlags
	assert.NotZero(flags&createNoWindowFlag,
		"daemon must run on a hidden console (CREATE_NO_WINDOW)")
	assert.Zero(flags&detachedProcessFlag,
		"daemon must not use DETACHED_PROCESS (children would open visible console windows)")
}
