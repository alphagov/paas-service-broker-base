package mock_locket_server_test

import (
	"io/ioutil"
	"os"

	"github.com/alphagov/paas-service-broker-base/testing/mock_locket_server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Locket Fixtures", func() {
	It("Creates and removes Locket fixtures", func() {
		f, err := mock_locket_server.SetupLocketFixtures()
		Expect(err).NotTo(HaveOccurred())
		_, err = os.Stat(f.Filepath + "/locket-client.cert.pem")
		Expect(err).NotTo(HaveOccurred())
		clientCert, err := ioutil.ReadFile(f.Filepath + "/locket-client.cert.pem")
		Expect(err).ToNot(HaveOccurred())
		Expect(string(clientCert)).To(Equal(mock_locket_server.LocketClientCertPEM))
		f.Cleanup()
		_, err = os.Stat(f.Filepath + "/locket-client.cert.pem")
		Expect(err).To(HaveOccurred())
	})
})
