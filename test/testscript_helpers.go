// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

const (
	waitInterval     = 2 * time.Second
	shortWaitTimeout = 30 * time.Second
	longWaitTimeout  = 5 * time.Minute
	jobinfoArgs      = 2
	podinfoArgs      = 2
)

type WrappedM struct {
	m       *testing.M
	testenv env.Environment
}

func (w WrappedM) Run() int {
	return w.testenv.Run(w.m)
}

func NewWrappedM(m *testing.M, testenv env.Environment) WrappedM {
	return WrappedM{
		m:       m,
		testenv: testenv,
	}
}

func getChaincodePackageLabels(script *testscript.TestScript, args []string) (string, string) {
	cclabel := args[0]
	cchash := args[1]

	packageHashBytes, err := hex.DecodeString(cchash)
	if err != nil {
		script.Fatalf("error decoding chaincode package hash %v: %v", cchash, err)
	}

	encodedPackageHash := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(packageHashBytes)

	return cclabel, encodedPackageHash
}

func getConfig(script *testscript.TestScript) *envconf.Config {
	testenv, ok := script.Value("testenv").(env.Environment)
	if !ok {
		script.Logf("could not get testenv")
	}

	cfg := testenv.EnvConf()

	return cfg
}

func waitForChaincodeJob(script *testscript.TestScript, cfg *envconf.Config, cclabel string, cchash string) *batchv1.Job {
	script.Logf("Waiting for job to be created for chaincode %s", cclabel)

	jobs := &batchv1.JobList{}

	err := wait.For(conditions.New(cfg.Client().Resources()).ResourceListN(jobs, 1,
		resources.WithLabelSelector(
			labels.FormatLabels(map[string]string{
				"fabric-builder-k8s-cclabel": cclabel,
				"fabric-builder-k8s-cchash":  cchash,
			}))), wait.WithInterval(waitInterval), wait.WithTimeout(shortWaitTimeout))
	if err != nil {
		script.Fatalf("failed waiting for job to be created for chaincode %s: %v", cclabel, err)
	}

	job := &jobs.Items[0]

	return job
}

func waitForChaincodePod(script *testscript.TestScript, cfg *envconf.Config, jobname string) *v1.Pod {
	script.Logf("Waiting for pod to be created for chaincode job %s", jobname)

	pods := &v1.PodList{}

	err := wait.For(conditions.New(cfg.Client().Resources()).ResourceListN(pods, 1,
		resources.WithLabelSelector(
			labels.FormatLabels(map[string]string{
				"batch.kubernetes.io/job-name": jobname,
			}))), wait.WithInterval(waitInterval), wait.WithTimeout(shortWaitTimeout))
	if err != nil {
		script.Fatalf("failed waiting for pod to be created for chaincode job %s: %v", jobname, err)
	}

	pod := &pods.Items[0]
	podname := pod.GetName()

	script.Logf("Waiting for pod %s to start", podname)

	err = wait.For(conditions.New(cfg.Client().Resources()).PodReady(pod), wait.WithInterval(waitInterval), wait.WithTimeout(longWaitTimeout))
	if err != nil {
		script.Fatalf("failed to wait for chaincode pod %s to reach Ready condition: %v", podname, err)
	}

	return pod
}

func jobInfoCmd(script *testscript.TestScript, _ bool, args []string) {
	if len(args) != jobinfoArgs {
		script.Fatalf("usage: jobinfo chaincode_label chaincode_hash")
	}

	cclabel, cchash := getChaincodePackageLabels(script, args)

	cfg := getConfig(script)

	job := waitForChaincodeJob(script, cfg, cclabel, cchash)

	var err error
	_, err = script.Stdout().Write([]byte(fmt.Sprintf("Job name: %s\n", job.GetName())))
	script.Check(err)

	_, err = script.Stdout().Write([]byte(fmt.Sprintf("Job namespace: %s\n", job.GetNamespace())))
	script.Check(err)

	_, err = script.Stdout().Write([]byte("Job labels:\n"))
	script.Check(err)

	for k, v := range job.GetLabels() {
		_, err = script.Stdout().Write([]byte(fmt.Sprintf("%s=%s\n", k, v)))
		script.Check(err)
	}

	_, err = script.Stdout().Write([]byte("Job annotations:\n"))
	script.Check(err)

	for k, v := range job.GetAnnotations() {
		_, err = script.Stdout().Write([]byte(fmt.Sprintf("%s=%s\n", k, v)))
		script.Check(err)
	}
}

func podInfoCmd(script *testscript.TestScript, _ bool, args []string) {
	if len(args) != podinfoArgs {
		script.Fatalf("usage: podinfo chaincode_label chaincode_hash")
	}

	cclabel, cchash := getChaincodePackageLabels(script, args)

	cfg := getConfig(script)

	job := waitForChaincodeJob(script, cfg, cclabel, cchash)
	jobname := job.GetName()

	pod := waitForChaincodePod(script, cfg, jobname)
	podname := pod.GetName()

	var err error
	_, err = script.Stdout().Write([]byte(fmt.Sprintf("Pod name: %s\n", podname)))
	script.Check(err)

	_, err = script.Stdout().Write([]byte(fmt.Sprintf("Pod namespace: %s\n", pod.GetNamespace())))
	script.Check(err)

	_, err = script.Stdout().Write([]byte(fmt.Sprintf("Pod service account: %s\n", pod.Spec.ServiceAccountName)))
	script.Check(err)

	_, err = script.Stdout().Write([]byte("Pod labels:\n"))
	script.Check(err)

	for k, v := range pod.GetLabels() {
		_, err = script.Stdout().Write([]byte(fmt.Sprintf("%s=%s\n", k, v)))
		script.Check(err)
	}

	_, err = script.Stdout().Write([]byte("Pod annotations:\n"))
	script.Check(err)

	for k, v := range pod.GetAnnotations() {
		_, err = script.Stdout().Write([]byte(fmt.Sprintf("%s=%s\n", k, v)))
		script.Check(err)
	}
}

func setupTestscriptEnv(t *testing.T, tsenv *testscript.Env, testenv env.Environment) error {
	t.Helper()

	tsenv.Setenv("KUBECONFIG_PATH", testenv.EnvConf().KubeconfigFile())
	tsenv.Setenv("TESTENV_NAMESPACE", testenv.EnvConf().Namespace())
	tsenv.Values["testenv"] = testenv

	return nil
}

func NewTestscriptParams(t *testing.T, scriptfile string, testenv env.Environment) testscript.Params {
	t.Helper()

	keepWorkDirs, _ := strconv.ParseBool(os.Getenv("KEEP_TESTSCRIPT_DIRS"))

	params := testscript.Params{
		Files:               []string{scriptfile},
		RequireExplicitExec: true,
		Setup:               func(e *testscript.Env) error { return setupTestscriptEnv(t, e, testenv) },
		TestWork:            keepWorkDirs,
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"jobinfo": jobInfoCmd,
			"podinfo": podInfoCmd,
		},
	}

	return params
}
