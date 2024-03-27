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

		assetPrivateDataCollectionIndexPath := filepath.Join(tempDir, "statedb", "couchdb", "collections", "assetCollection", "indexes", "indexOwner.json")
		Expect(assetPrivateDataCollectionIndexPath).To(BeARegularFile(), "Private data index should be copied")

		fabCarPrivateDataCollectionIndexPath := filepath.Join(tempDir, "statedb", "couchdb", "collections", "fabCarCollection", "indexes", "indexOwner.json")
		Expect(fabCarPrivateDataCollectionIndexPath).To(BeARegularFile(), "Private data index should be copied")

		textPath := filepath.Join(tempDir, "statedb", "couchdb", "indexes", "test.txt")
		Expect(textPath).NotTo(BeAnExistingFile(), "Unexpected files should not be copied")

		subdirPath := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"indexes",
			"subdir",
			"indexOwner.json",
		)
		Expect(subdirPath).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")

		privateDataCollectionSubdirPath := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"collections",
			"fabCarCollection",
			"subdir",
			"indexes",
			"indexOwner.json",
		)
		Expect(privateDataCollectionSubdirPath).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")

		collectionsdCollectionPath := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"collectionsd",
			"fabCarCollection",
			"indexes",
			"indexOwner.json",
		)
		Expect(collectionsdCollectionPath).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")

		indexedCollectionSubdirPath := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"indexed",
			"indexes",
			"indexOwner.json",
		)
		Expect(indexedCollectionSubdirPath).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")

		rootIndexOwnerJSONFile := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"indexOwner.json",
		)
		Expect(rootIndexOwnerJSONFile).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")

		roottestTXTFile := filepath.Join(
			tempDir,
			"statedb",
			"couchdb",
			"test.txt",
		)
		Expect(roottestTXTFile).NotTo(BeAnExistingFile(), "Files outside indexes directory should not be copied")
	})
})
