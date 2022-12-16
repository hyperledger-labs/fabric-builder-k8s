package main_test

import (
	"os"
	"testing"
	"time"

	"github.com/bitfield/script"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

//nolint:gochecknoglobals // not sure how to avoid this
var (
	includeKindTests bool
	runCmdPath       string
)

func TestRun(t *testing.T) {
	RegisterFailHandler(Fail)

	suiteConfig, _ := GinkgoConfiguration()

	kindEnv := os.Getenv("INCLUDE_KIND_TESTS")
	if kindEnv == "true" {
		includeKindTests = true
	} else {
		includeKindTests = false
	}

	if !includeKindTests {
		if suiteConfig.LabelFilter == "" {
			suiteConfig.LabelFilter = "!kind"
		} else {
			suiteConfig.LabelFilter = "(" + suiteConfig.LabelFilter + ") && !kind"
		}
	}

	RunSpecs(t, "Run Suite", suiteConfig)
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(2 * time.Second)

	var err error
	runCmdPath, err = gexec.Build("github.com/hyperledger-labs/fabric-builder-k8s/cmd/run")
	Expect(err).NotTo(HaveOccurred())

	if includeKindTests {
		script.Exec("kind delete cluster --name fabric-builder-k8s-test")

		pipe := script.Exec("kind create cluster --name fabric-builder-k8s-test")
		_, err = pipe.Stdout()
		Expect(err).NotTo(HaveOccurred())
		Expect(pipe.ExitStatus()).To(Equal(0))

		pipe = script.Exec("kubectl create namespace chaincode")
		_, err = pipe.Stdout()
		Expect(err).NotTo(HaveOccurred())
		Expect(pipe.ExitStatus()).To(Equal(0))

		pipe = script.Exec("kubectl create serviceaccount chaincode --namespace=chaincode")
		_, err = pipe.Stdout()
		Expect(err).NotTo(HaveOccurred())
		Expect(pipe.ExitStatus()).To(Equal(0))
	}
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	if includeKindTests {
		_, err := script.Exec("kind delete cluster --name fabric-builder-k8s-test").Stdout()
		Expect(err).NotTo(HaveOccurred())
	}
})
