// SPDX-License-Identifier: Apache-2.0

//go:build linux
// +build linux

package integration_test

import (
	"os"
	"testing"

	"github.com/hyperledger-labs/fabric-builder-k8s/cmd"
	"github.com/hyperledger-labs/fabric-builder-k8s/test"
	"github.com/rogpeppe/go-internal/testscript"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
)

//nolint:gochecknoglobals // not sure how to avoid this
var (
	testenv env.Environment
)

func TestMain(m *testing.M) {
	envCfg := envconf.New().WithRandomNamespace()
	testenv = env.NewWithConfig(envCfg)
	clusterName := envconf.RandomName("test-cluster", 16)

	testenv.Setup(
		envfuncs.CreateCluster(kind.NewProvider(), clusterName),
		envfuncs.CreateNamespace(envCfg.Namespace()),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace(envCfg.Namespace()),
		envfuncs.DestroyCluster(clusterName),
	)

	wm := test.NewWrappedM(m, testenv)
	os.Exit(testscript.RunMain(wm, map[string]func() int{
		"run": cmd.Run,
	}))
}
