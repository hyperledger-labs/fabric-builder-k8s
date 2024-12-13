// SPDX-License-Identifier: Apache-2.0

//go:build linux
// +build linux

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger-labs/fabric-builder-k8s/cmd"
	"github.com/hyperledger-labs/fabric-builder-k8s/test"
	"github.com/rogpeppe/go-internal/testscript"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
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
		envfuncs.CreateClusterWithConfig(kind.NewProvider(), clusterName, "testdata/kind-config.yaml"),
		envfuncs.CreateNamespace(envCfg.Namespace()),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace(envCfg.Namespace()),
		// envfuncs.ExportClusterLogs(kindClusterName, "./logs"),
		envfuncs.DestroyCluster(clusterName),
	)

	testenv.AfterEachTest(func(ctx context.Context, cfg *envconf.Config, t *testing.T) (context.Context, error) { //nolint:thelper // *testing.T must be last param for TestEnvFunc
		t.Helper()

		t.Logf("Deleting jobs after test %s", t.Name())

		client, err := cfg.NewClient()
		if err != nil {
			return ctx, fmt.Errorf("delete jobs func: %w", err)
		}

		jobs := new(batchv1.JobList)

		err = client.Resources(cfg.Namespace()).List(ctx, jobs)
		if err != nil {
			return ctx, fmt.Errorf("delete jobs func: %w", err)
		}

		for _, job := range jobs.Items {
			if err := client.Resources().Delete(ctx, &job, resources.WithDeletePropagation(string(metav1.DeletePropagationBackground))); err != nil {
				return ctx, fmt.Errorf("delete jobs func: %w", err)
			}
		}

		return ctx, nil
	})

	wm := test.NewWrappedM(m, testenv)
	os.Exit(testscript.RunMain(wm, map[string]func() int{
		"run": cmd.Run,
	}))
}
