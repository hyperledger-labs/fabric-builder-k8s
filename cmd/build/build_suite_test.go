package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

//nolint:gochecknoglobals // not sure how to avoid this
var buildCmdPath string

func TestBuild(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Suite")
}

var _ = BeforeSuite(func() {
	var err error
	buildCmdPath, err = gexec.Build("github.com/hyperledger-labs/fabric-builder-k8s/cmd/build")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
