package util_test

import (
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {
	Describe("GetValidRfc1035LabelName", func() {
		It("should return names with a maximum of 63 characters", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgMsp",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgPeer0", chaincodeData, 0)
			Expect(len(name)).To(Equal(63))
		})

		It("should return names with a maximum of 57 characters if a 6 character suffix is required", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgMsp",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgPeer0", chaincodeData, 6)
			Expect(len(name)).To(Equal(57))
		})

		It("should return names which starts with an alphabetic character", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "GreenCongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "GreenCongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(MatchRegexp("^[a-z]"))
		})

		It("should return names which end with an alphanumeric character", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "BlueCongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "BlueCongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(MatchRegexp("[a-z0-9]$"))
		})

		It("should return names which only contains lowercase alphanumeric characters or '-'", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "FAB/CAR*:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "BlueCongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "BlueCongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(MatchRegexp("^(?:[a-z0-9]|-)+$"))
		})

		It("should return different names for the same package IDs", func() {
			chaincodeData1 := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "GreenCongaOrg",
			}
			chaincodeData2 := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org2.example.org",
				MspID:       "BlueCongaOrg",
			}
			name1 := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "GreenCongaOrgPeer0", chaincodeData1, 0)
			name2 := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "BlueCongaOrgPeer0", chaincodeData2, 0)
			Expect(name1).NotTo(Equal(name2))
		})

		It("should return different names for different package IDs", func() {
			chaincodeData1 := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "RedCongaOrg",
			}
			chaincodeData2 := &util.ChaincodeJSON{
				ChaincodeID: "go-contract:6f98c4bb29414771312eddd1a813eef583df2121c235c4797792f141a46d4b45",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "RedCongaOrg",
			}
			name1 := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "RedCongaOrg", chaincodeData1, 0)
			name2 := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "RedCongaOrg", chaincodeData2, 0)
			Expect(name1).NotTo(Equal(name2))
		})

		It("should return deterministic names", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(Equal("hlf-k8sbuilder-ftw-fabcar-s6pwkq6bepi2e"))
		})

		It("should return names which start with the specified prefix and a safe version of the chaincode label", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "FAB/CAR*:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(HavePrefix("hlf-k8sbuilder-ftw" + "-fabcar-"))
		})

		It("should return names which end with a 13 character lowercase base32 encoded hash string", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaOrg",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(MatchRegexp("-[a-z2-7]{13}$"))
		})

		It("should return names with the full prefix and hash, and a truncated chaincode label", func() {
			chaincodeData := &util.ChaincodeJSON{
				ChaincodeID: "fabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcarfabfabfabfabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b",
				PeerAddress: "peer0.org1.example.com",
				MspID:       "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgMsp",
			}
			name := util.GetValidRfc1035LabelName("hlf-k8sbuilder-ftw", "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgPeer0", chaincodeData, 0)
			Expect(name).To(Equal("hlf-k8sbuilder-ftw-fabfabfabfabcarfabfabfabfabcar-b46p74k4ygwh6"))
		})
	})

	Describe("ParseAnnotations", func() {
		It("should return empty map for empty string", func() {
			result := util.ParseAnnotations("")
			Expect(result).To(BeEmpty())
		})

		It("should parse single annotation", func() {
			result := util.ParseAnnotations("sidecar.istio.io/inject=true")
			Expect(result).To(HaveLen(1))
			Expect(result["sidecar.istio.io/inject"]).To(Equal("true"))
		})

		It("should parse multiple annotations", func() {
			result := util.ParseAnnotations("sidecar.istio.io/inject=true,app=myapp,version=1.0")
			Expect(result).To(HaveLen(3))
			Expect(result["sidecar.istio.io/inject"]).To(Equal("true"))
			Expect(result["app"]).To(Equal("myapp"))
			Expect(result["version"]).To(Equal("1.0"))
		})

		It("should handle annotations with spaces", func() {
			result := util.ParseAnnotations(" sidecar.istio.io/inject = true , app = myapp ")
			Expect(result).To(HaveLen(2))
			Expect(result["sidecar.istio.io/inject"]).To(Equal("true"))
			Expect(result["app"]).To(Equal("myapp"))
		})

		It("should skip invalid entries without equals sign", func() {
			result := util.ParseAnnotations("sidecar.istio.io/inject=true,invalidentry,app=myapp")
			Expect(result).To(HaveLen(2))
			Expect(result["sidecar.istio.io/inject"]).To(Equal("true"))
			Expect(result["app"]).To(Equal("myapp"))
		})

		It("should skip empty entries", func() {
			result := util.ParseAnnotations("sidecar.istio.io/inject=true,,app=myapp")
			Expect(result).To(HaveLen(2))
			Expect(result["sidecar.istio.io/inject"]).To(Equal("true"))
			Expect(result["app"]).To(Equal("myapp"))
		})

		It("should handle annotations with empty values", func() {
			result := util.ParseAnnotations("sidecar.istio.io/inject=,app=myapp")
			Expect(result).To(HaveLen(2))
			Expect(result["sidecar.istio.io/inject"]).To(Equal(""))
			Expect(result["app"]).To(Equal("myapp"))
		})

		It("should handle annotations with equals signs in values", func() {
			result := util.ParseAnnotations("config=key=value,app=myapp")
			Expect(result).To(HaveLen(2))
			Expect(result["config"]).To(Equal("key=value"))
			Expect(result["app"]).To(Equal("myapp"))
		})

		It("should skip entries with empty keys", func() {
			result := util.ParseAnnotations("=value,app=myapp")
			Expect(result).To(HaveLen(1))
			Expect(result["app"]).To(Equal("myapp"))
		})
	})

})
