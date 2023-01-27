package util_test

import (
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {
	Describe("GetValidName", func() {
		It("should return a string with a maximum of 63 characters", func() {
			name := util.GetValidName("CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgMsp", "CongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaCongaOrgPeer0", "fabfabfabfabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(len(name)).To(BeNumerically("<=", 63))
		})

		It("should return a string which starts with an alphabetic character", func() {
			name := util.GetValidName("GreenCongaOrg", "GreenCongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(name).To(MatchRegexp("^[a-z]"))
		})

		It("should return a string which ends with an alphanumeric character", func() {
			name := util.GetValidName("BlueCongaOrg", "BlueCongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(name).To(MatchRegexp("[a-z0-9]$"))
		})

		It("should return a string which only contains lowercase alphanumeric characters or '-'", func() {
			name := util.GetValidName("BlueCongaOrg", "BlueCongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(name).To(MatchRegexp("^(?:[a-z0-9]|-)+$"))
		})

		It("should return different names for different input", func() {
			name1 := util.GetValidName("GreenCongaOrg", "GreenCongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			name2 := util.GetValidName("BlueCongaOrg", "BlueCongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(name1).NotTo(Equal(name2))
		})

		It("should return deterministic names", func() {
			name := util.GetValidName("CongaOrg", "CongaOrgPeer0", "fabcar:cffa266294278404e5071cb91150d550dc0bf855149908a170b1169d6160004b")
			Expect(name).To(Equal("cc-ocqvh9ir0mi0ef6urh12f3l0dar6csdmtjfhgbfvdp2d22u109r0"))
		})
	})
})
