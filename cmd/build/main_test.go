package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {

	var tempDir string
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
	})

	DescribeTable("Running the detect command produces the correct error code",
		func(expectedErrorCode int, getArgs func() []string) {
			args := getArgs()
			command := exec.Command(buildCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(expectedErrorCode))
		},
		Entry("When the image.json file exists", 0, func() []string { return []string{"./testdata/validimage", "CHAINCODE_METADATA_DIR", tempDir} }),
		Entry("When the image.json file does not exist", 1, func() []string { return []string{"CHAINCODE_SOURCE_DIR", "CHAINCODE_METADATA_DIR", "BUILD_OUTPUT_DIR"} }),
		Entry("When too few arguments are provided", 1, func() []string { return []string{"CHAINCODE_SOURCE_DIR"} }),
		Entry("When too many arguments are provided", 1, func() []string {
			return []string{"CHAINCODE_SOURCE_DIR", "CHAINCODE_METADATA_DIR", "BUILD_OUTPUT_DIR", "UNEXPECTED_ARGUMENT"}
		}),
	)
})
