package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	DescribeTable("Running the detect command produces the correct error code",
		func(expectedErrorCode int, args ...string) {
			command := exec.Command(detectCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(expectedErrorCode))
		},
		Entry("When the metadata contains a valid type", 0, "CHAINCODE_SOURCE_DIR", "./testdata/validtype"),
		Entry("When the metadata contains an invalid type", 1, "CHAINCODE_SOURCE_DIR", "./testdata/invalidtype"),
		Entry("When the metadata contents are invalid", 1, "CHAINCODE_SOURCE_DIR", "./testdata/invalidfile"),
		Entry("When the metadata does not exist", 1, "CHAINCODE_SOURCE_DIR", "CHAINCODE_METADATA_DIR"),
		Entry("When too few arguments are provided", 1, "CHAINCODE_SOURCE_DIR"),
		Entry(
			"When too many arguments are provided",
			1,
			"CHAINCODE_SOURCE_DIR",
			"CHAINCODE_METADATA_DIR",
			"UNEXPECTED_ARGUMENT",
		),
	)
})
