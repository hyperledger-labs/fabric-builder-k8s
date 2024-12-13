// SPDX-License-Identifier: Apache-2.0

//go:build linux
// +build linux

package integration_test

import (
	"context"
	"testing"

	"github.com/hyperledger-labs/fabric-builder-k8s/test"
	"github.com/rogpeppe/go-internal/testscript"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestRunChaincode(t *testing.T) {
	testenv.Test(t, features.NewWithDescription(t.Name()+"Feature", "the builder should run chaincode in the specified namespace").
		Assess(t.Name()+"Assessment", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Helper()

			testscript.Run(t, test.NewTestscriptParams(t, "testdata/run_chaincode.txtar", testenv))

			return ctx
		}).Feature())
}

func TestRunChaincodeWithNamePrefix(t *testing.T) {
	testenv.Test(t, features.NewWithDescription(t.Name()+"Feature", "the builder should run chaincode using kubernetes object names with the specified prefix").
		Assess(t.Name()+"Assessment", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Helper()

			testscript.Run(t, test.NewTestscriptParams(t, "testdata/chaincode_name_prefix.txtar", testenv))

			return ctx
		}).Feature())
}

func TestRunChaincodeWithServiceAccount(t *testing.T) {
	testenv.Test(t, features.NewWithDescription(t.Name()+"Feature", "the builder should run chaincode in with the specified service account").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Helper()

			t.Logf("Creating service account before test %s", t.Name())
			serviceAccount := &v1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{Name: "chaincode", Namespace: cfg.Namespace()},
			}
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err := client.Resources().Create(ctx, serviceAccount); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess(t.Name()+"Assessment", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Helper()

			testscript.Run(t, test.NewTestscriptParams(t, "testdata/chaincode_service_account.txtar", testenv))

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Helper()

			t.Logf("Deleting service account after test %s", t.Name())
			serviceAccount := &v1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{Name: "chaincode", Namespace: cfg.Namespace()},
			}
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err := client.Resources().Delete(ctx, serviceAccount); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature())
}

func TestRunChaincodeWithAvailableNodeRole(t *testing.T) {
	testenv.Test(t, features.NewWithDescription(t.Name()+"Feature", "the builder should run chaincode on a dedicated kubernetes node").
		Assess(t.Name()+"Assessment", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Helper()

			testscript.Run(t, test.NewTestscriptParams(t, "testdata/dedicated_node_available.txtar", testenv))

			return ctx
		}).Feature())
}

func TestRunChaincodeWithoutAvailableNodeRole(t *testing.T) {
	testenv.Test(t, features.NewWithDescription(t.Name()+"Feature", "the builder should fail to run chaincode if a kubernetes node with the required role is not available").
		Assess(t.Name()+"Assessment", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Helper()

			testscript.Run(t, test.NewTestscriptParams(t, "testdata/dedicated_node_unavailable.txtar", testenv))

			return ctx
		}).Feature())
}
