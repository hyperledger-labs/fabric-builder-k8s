package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	DescribeTable("Running the run command produces the correct error code",
		func(expectedErrorCode int, args ...string) {
			command := exec.Command(runCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(expectedErrorCode))
		},
		Entry("When too few arguments are provided", 1, "BUILD_OUTPUT_DIR"),
		Entry(
			"When too many arguments are provided",
			1,
			"BUILD_OUTPUT_DIR",
			"RUN_METADATA_DIR",
			"UNEXPECTED_ARGUMENT",
		),
	)
})
