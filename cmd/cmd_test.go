// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"os"
	"testing"

	"github.com/hyperledger-labs/fabric-builder-k8s/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		// "build": cmd.Build,
		"detect": cmd.Detect,
		// "release": cmd.Release,
		// "run": cmd.Run,
	}))
}

func TestBuildCommand(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/build",
	})
}

func TestDetectCommand(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/detect",
	})
}

func TestReleaseCommand(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/release",
	})
}

func TestRunCommand(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/run",
	})
}
