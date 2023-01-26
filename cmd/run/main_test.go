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
			).Should(gbytes.Say(`run \[\d+\]: Running chaincode ID CHAINCODE_ID in kubernetes pod chaincode/cc-mspid-core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789chai`))

			pipe := script.Exec(
				"kubectl wait --for=condition=ready pod --timeout=120s --namespace=chaincode -l fabric-builder-k8s-peerid=core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789",
			)
			_, err = pipe.Stdout()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipe.ExitStatus()).To(Equal(0))

			descArgs := []string{
				"describe",
				"pod",
				"--namespace=chaincode",
				"-l",
				"fabric-builder-k8s-peerid=core-peer-id-abcdefghijklmnopqrstuvwxyz-0123456789",
			}
			descCommand := exec.Command("kubectl", descArgs...)
			descSession, err := gexec.Start(descCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(descSession).Should(gexec.Exit(0))
			Eventually(descSession.Out).Should(gbytes.Say(`Namespace:\s+chaincode`))
			Eventually(descSession.Out).Should(gbytes.Say(`fabric-builder-k8s-mspid=MSPID`))
			Eventually(
				descSession.Out,
			).Should(gbytes.Say(`fabric-builder-k8s-ccid:\s+CHAINCODE_ID`))
			Eventually(descSession.Out).Should(gbytes.Say(`CORE_CHAINCODE_ID_NAME:\s+CHAINCODE_ID`))
			Eventually(descSession.Out).Should(gbytes.Say(`CORE_PEER_ADDRESS:\s+PEER_ADDRESS`))
			Eventually(descSession.Out).Should(gbytes.Say(`CORE_PEER_LOCALMSPID:\s+MSPID`))
		},
	)
})
