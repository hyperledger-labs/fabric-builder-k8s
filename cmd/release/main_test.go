package main_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	var tempDir string
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
	})

	DescribeTable("Running the release command produces the correct error code",
		func(expectedErrorCode int, getArgs func() []string) {
			args := getArgs()
			command := exec.Command(releaseCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(expectedErrorCode))
		},
		Entry("When there is no chaincode metadata", 0, func() []string {
			return []string{"BUILD_OUTPUT_DIR", "RELEASE_OUTPUT_DIR"}
		}),
		Entry("When there is chaincode metadata", 0, func() []string {
			return []string{"./testdata/buildwithindexes", tempDir}
		}),
		Entry("When too few arguments are provided", 1, func() []string {
			return []string{"BUILD_OUTPUT_DIR"}
		}),
		Entry("When too many arguments are provided", 1, func() []string {
			return []string{"BUILD_OUTPUT_DIR", "RELEASE_OUTPUT_DIR", "UNEXPECTED_ARGUMENT"}
		}),
	)

	It("should only copy .json CouchDB index definitions to the release output directory", func() {
		args := []string{"./testdata/buildwithindexes", tempDir}
		command := exec.Command(releaseCmdPath, args...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		indexPath := filepath.Join(tempDir, "statedb", "couchdb", "indexes", "indexOwner.json")
		Expect(indexPath).To(BeARegularFile())
		textPath := filepath.Join(tempDir, "statedb", "couchdb", "indexes", "test.txt")
		Expect(textPath).NotTo(BeAnExistingFile())
		subdirPath := filepath.Join(tempDir, "statedb", "couchdb", "indexes", "subdir", "indexOwner.json")
		Expect(subdirPath).NotTo(BeAnExistingFile())
	})
})
