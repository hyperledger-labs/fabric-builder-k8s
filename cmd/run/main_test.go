package main_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitfield/script"
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

	It(
		"should start a chaincode pod using the supplied configuration environment variables",
		Label("kind"),
		func() {
			homedir, err := os.UserHomeDir()
			Expect(err).NotTo(HaveOccurred())

			args := []string{"./testdata/validimage", "./testdata/validchaincode"}
			command := exec.Command(runCmdPath, args...)
			command.Env = append(os.Environ(),
				fmt.Sprintf("KUBECONFIG_PATH=%s/.kube/config", homedir),
				"CORE_PEER_ID=core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789",
				"FABRIC_K8S_BUILDER_DEBUG=true",
				"FABRIC_K8S_BUILDER_NAMESPACE=chaincode",
				"FABRIC_K8S_BUILDER_SERVICE_ACCOUNT=chaincode",
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).ShouldNot(gexec.Exit())
			Eventually(
				session.Err,
			).Should(gbytes.Say(`run \[\d+\] DEBUG: FABRIC_K8S_BUILDER_NAMESPACE=chaincode`))
			Eventually(
				session.Err,
			).Should(gbytes.Say(`run \[\d+\] DEBUG: FABRIC_K8S_BUILDER_SERVICE_ACCOUNT=chaincode`))
			Eventually(
				session.Err,
			).Should(gbytes.Say(`run \[\d+\]: Running chaincode ID CHAINCODE_LABEL:6f98c4bb29414771312eddd1a813eef583df2121c235c4797792f141a46d4b45 in kubernetes pod chaincode/hlfcc-chaincodelabel-f887209uhojj2`))

			pipe := script.Exec(
				"kubectl wait --for=condition=ready pod --timeout=120s --namespace=chaincode -l fabric-builder-k8s-cclabel=CHAINCODE_LABEL",
			)
			_, err = pipe.Stdout()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipe.ExitStatus()).To(Equal(0))

			descArgs := []string{
				"describe",
				"pod",
				"--namespace=chaincode",
				"-l",
				"fabric-builder-k8s-cclabel=CHAINCODE_LABEL",
			}
			descCommand := exec.Command("kubectl", descArgs...)
			descSession, err := gexec.Start(descCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(descSession).Should(gexec.Exit(0))
			Eventually(descSession.Out).Should(gbytes.Say(`Namespace:\s+chaincode`))
			Eventually(
				descSession.Out,
			).Should(gbytes.Say(`fabric-builder-k8s-ccid:\s+CHAINCODE_LABEL:6f98c4bb29414771312eddd1a813eef583df2121c235c4797792f141a46d4b45`))
			Eventually(descSession.Out).Should(gbytes.Say(`fabric-builder-k8s-mspid:\s+MSPID`))
			Eventually(descSession.Out).Should(gbytes.Say(`fabric-builder-k8s-peeraddress:\s+PEER_ADDRESS`))
			Eventually(descSession.Out).Should(gbytes.Say(`fabric-builder-k8s-peerid:\s+core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789`))
			Eventually(
				descSession.Out,
			).Should(gbytes.Say(`CORE_CHAINCODE_ID_NAME:\s+CHAINCODE_LABEL:6f98c4bb29414771312eddd1a813eef583df2121c235c4797792f141a46d4b45`))
			Eventually(descSession.Out).Should(gbytes.Say(`CORE_PEER_ADDRESS:\s+PEER_ADDRESS`))
			Eventually(descSession.Out).Should(gbytes.Say(`CORE_PEER_LOCALMSPID:\s+MSPID`))
		},
	)
})
