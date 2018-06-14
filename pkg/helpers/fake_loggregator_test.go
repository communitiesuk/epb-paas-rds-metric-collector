package helpers_test

import (
	"time"

	loggregator "code.cloudfoundry.org/go-loggregator"

	"github.com/alphagov/paas-rds-metric-collector/pkg/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IngressClient", func() {
	var (
		server *helpers.FakeLoggregatorIngressServer
		client *loggregator.IngressClient
	)

	BeforeEach(func() {
		var err error
		server, err = helpers.NewFakeLoggregatorIngressServer(
			"../../fixtures/server.crt",
			"../../fixtures/server.key",
			"../../fixtures/CA.crt",
		)
		Expect(err).NotTo(HaveOccurred())

		err = server.Start()
		Expect(err).NotTo(HaveOccurred())

		tlsConfig, err := loggregator.NewIngressTLSConfig(
			"../../fixtures/CA.crt",
			"../../fixtures/server.crt",
			"../../fixtures/server.key",
		)
		Expect(err).NotTo(HaveOccurred())

		client, err = loggregator.NewIngressClient(
			tlsConfig,
			loggregator.WithAddr(server.Addr),
			loggregator.WithTag("origin", "rds-metrics-collector"),
		)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.Stop()
	})

	It("should receive one envelope from one metric", func() {
		client.EmitGauge(
			loggregator.WithGaugeValue("test", 1, "s"),
		)

		Eventually(
			server.Receivers,
			5*time.Second,
		).Should(Receive())
	})

	PIt("should receive more than one envelope", func() {
		client.EmitGauge(
			loggregator.WithGaugeValue("test", 1, "s"),
		)

		Eventually(
			server.Receivers,
			5*time.Second,
		).Should(Receive())

		time.Sleep(200 * time.Millisecond)

		client.EmitGauge(
			loggregator.WithGaugeValue("test", 2, "s"),
		)

		Eventually(
			server.Receivers,
			5*time.Second,
		).Should(Receive())
	})
})