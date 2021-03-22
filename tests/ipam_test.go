package tests_test

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/prefix"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IPAM API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Address endpoint", func() {

		It("Should list all available addresses", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := address.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Prefix endpoint", func() {

		It("Should list all prefixes", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := prefix.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should create a new prefix and delete it later", func() {
			p := prefix.NewAPI(cli)
			ipV4 := 4
			networkMask := 24
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()

			By("Creating a new prefix")
			summary, err := p.Create(ctx, prefix.NewCreate(locationID, vlanID, ipV4, prefix.TypePrivate, networkMask))
			Expect(err).NotTo(HaveOccurred())

			var info prefix.Info
			By("Waiting for prefix to be 'Active'")
			Eventually(func() string {
				info, err = p.Get(ctx, summary.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.Vlans).NotTo(BeNil())
				return info.Status
			}, 15*time.Minute, 5*time.Second).Should(Equal("Active"))

			Expect(info.Vlans[0].ID).To(Equal(vlanID))

			By("Updating the prefix")
			_, err = p.Update(ctx, summary.ID, prefix.Update{CustomerDescription: "something else"})
			Expect(err).NotTo(HaveOccurred())

			By("Deleting the prefix")
			err = p.Delete(ctx, summary.ID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should reserve a random prefix for a given VLAN", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()

			By("Reserving a new IP address")
			res, err := address.NewAPI(cli).ReserveRandom(ctx, address.ReserveRandom{
				LocationID: locationID,
				VlanID:     vlanID,
				Count:      1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Data)).To(Equal(1))

		})
	})
})
