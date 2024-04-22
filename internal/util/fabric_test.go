package util_test

import (
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fabric", func() {
	DescribeTable("NewChaincodePackageID return a new ChaincodePackageID with the expected label and hash values",
		func(chaincodeID, expectedLabel, expectedHash string) {
			packageID := util.NewChaincodePackageID(chaincodeID)
			Expect(packageID.Label).To(Equal(expectedLabel), "The ChaincodePackageID should include the expected label")
			Expect(packageID.Hash).To(Equal(expectedHash), "The ChaincodePackageID should include the expected hash")
		},
		Entry("When the chaincode ID only contains one colon", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b", "fabcar", "cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b"),
		// The rest are a bit of a guess since I'm not sure the package ID format is defined in detail anywhere
		Entry("When the chaincode ID contains more than one colon", "fab:car:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b", "fab:car", "cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b"),
		Entry("When the chaincode ID contains a double colon", "fab::car:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b", "fab::car", "cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b"),
		Entry("When the chaincode ID is an empty string", "", "", ""),
		Entry("When the chaincode ID does not contain a colon", "fabcar", "", ""),
	)
})
