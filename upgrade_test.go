//go:build bdd
// +build bdd

package emoney_test

import (
	nt "github.com/e-money/em-ledger/networktest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tidwall/gjson"
)

var _ = Describe("Upgrade", func() {
	emcli := testnet.NewEmcli()

	var Authority = testnet.Keystore.Authority

	Describe("Authority manages issuers", func() {
		It("creates a new testnet", createNewTestnet)

		It("upgrade nodes and confirm", func() {
			const (
				name = "test-upg-0.2.0"
			)

			chainHeight, err := nt.GetHeight()
			Expect(err).ToNot(HaveOccurred())

			const upgDelta = 5
			upgHeight := chainHeight + upgDelta

			_, success, err := emcli.UpgSchedByHeight(Authority, name,
				upgHeight,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeTrue())

			bz, err := emcli.QueryUpgSched()
			Expect(err).ToNot(HaveOccurred())

			upgPlan := gjson.ParseBytes(bz).Get("plan")

			resName := upgPlan.Get("name").Str
			Expect(resName).To(Equal(name))
			resUpgHeight := upgPlan.Get("height").Int()
			Expect(resUpgHeight).To(BeEquivalentTo(upgHeight))

			// wait till the upgrade
			newHeight, err := nt.IncChain(upgDelta + 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(newHeight > chainHeight+upgDelta).To(BeTrue())

			// if we made it here, the upgrade succeeded
			bz, err = emcli.QueryUpgSched()
			Expect(err).ToNot(HaveOccurred())

			// assert that the upgrade plan zeroed out
			upgPlan = gjson.ParseBytes(bz).Get("plan")

			resName = upgPlan.Get("name").Str
			Expect(resName).To(HaveLen(0))
		})
	})
})
