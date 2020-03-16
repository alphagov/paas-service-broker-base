package broker_test

import (
	"testing"

	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/locket"
	"github.com/alphagov/paas-service-broker-base/testing/mock_locket_server"
)

var (
	mockLocket     *mock_locket_server.MockLocket
	locketFixtures mock_locket_server.LocketFixtures
)

func TestBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}

var _ = BeforeSuite(func() {
	var err error

	locketFixtures, err = mock_locket_server.SetupLocketFixtures()
	Expect(err).NotTo(HaveOccurred())
	mockLocket, err = mock_locket_server.New("keyBasedLock", locketFixtures.Filepath)
	Expect(err).NotTo(HaveOccurred())
	mockLocket.Start(mockLocket.Logger, mockLocket.ListenAddress, mockLocket.Certificate)

	locketLogger := lager.NewLogger("locket-test")
	Eventually(func() error {
		_, err := locket.NewClientSkipCertVerify(
			locketLogger,
			locket.ClientLocketConfig{
				LocketAddress:        mockLocket.ListenAddress,
				LocketCACertFile:     path.Join(locketFixtures.Filepath, "locket-server.cert.pem"),
				LocketClientCertFile: path.Join(locketFixtures.Filepath, "locket-client.cert.pem"),
				LocketClientKeyFile:  path.Join(locketFixtures.Filepath, "locket-client.key.pem"),
			},
		)
		return err
	}).Should(BeNil())
})

var _ = AfterSuite(func() {
	if mockLocket != nil {
		mockLocket.Stop()
	}
	locketFixtures.Cleanup()
})
