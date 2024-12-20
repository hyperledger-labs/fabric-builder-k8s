package main_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("should return an error if the CORE_PEER_ID environment variable is not set", func() {
		args := []string{"BUILD_OUTPUT_DIR", "RUN_METADATA_DIR"}
		command := exec.Command(runCmdPath, args...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(1))
		Eventually(
			session.Err,
		).Should(gbytes.Say(`run \[\d+\]: Expected CORE_PEER_ID environment variable`))
	})

	DescribeTable("Running the run command with the wrong arguments produces the correct error",
		func(args ...string) {
			command := exec.Command(runCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Eventually(
				session.Err,
			).Should(gbytes.Say(`run \[\d+\]: Expected BUILD_OUTPUT_DIR and RUN_METADATA_DIR arguments`))
		},
		Entry("When too few arguments are provided", "BUILD_OUTPUT_DIR"),
		Entry(
			"When too many arguments are provided",
			"BUILD_OUTPUT_DIR",
			"RUN_METADATA_DIR",
			"UNEXPECTED_ARGUMENT",
		),
	)

	DescribeTable("Running the run command produces the correct error for invalid FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX environment variable values",
		func(kubeNamePrefixValue, expectedErrorMessage string) {
			args := []string{"BUILD_OUTPUT_DIR", "RUN_METADATA_DIR"}
			command := exec.Command(runCmdPath, args...)
			command.Env = append(os.Environ(),
				"CORE_PEER_ID=core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789",
				"FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX="+kubeNamePrefixValue,
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Eventually(
				session.Err,
			).Should(gbytes.Say(expectedErrorMessage))
		},
		Entry("When the FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX is too long", "long-prefix-is-looooooooooooooooooooong", `run \[\d+\]: The FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX environment variable must be a maximum of 30 characters`),
		Entry("When the FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX contains invalid characters", "invalid/PREFIX*", `run \[\d+\]: The FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX environment variable must be a valid DNS-1035 label: a DNS-1035 label must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character`),
		Entry("When the FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX starts with a number", "1prefix", `run \[\d+\]: The FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX environment variable must be a valid DNS-1035 label: a DNS-1035 label must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character`),
		Entry("When the FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX starts with a dash", "-prefix", `run \[\d+\]: The FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX environment variable must be a valid DNS-1035 label: a DNS-1035 label must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character`),
	)
})
