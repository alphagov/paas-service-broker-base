package broker_test

import (
	"testing"

	"fmt"
	"os/exec"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/locket"
	"github.com/phayes/freeport"
)

var (
	locketSession *gexec.Session
	locketAddress string
	fixturesPath  string
)

func TestBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}

var _ = BeforeSuite(func() {
	var err error

	locketPath, err := gexec.Build("github.com/alphagov/paas-service-broker-base/testing/mock_locket_server")
	Expect(err).NotTo(HaveOccurred())

	port, err := freeport.GetFreePort()
	Expect(err).NotTo(HaveOccurred())
	locketAddress = fmt.Sprintf("127.0.0.1:%d", port)

	fixturesPath, err = filepath.Abs("../testing/mock_locket_server/fixtures")
	Expect(err).NotTo(HaveOccurred())

	command := exec.Command(
		locketPath,
		"-listenAddress="+locketAddress,
		"-fixturesPath="+fixturesPath,
	)
	locketSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(locketSession).Should(gbytes.Say("listen"))

	locketLogger := lager.NewLogger("locket-test")
	Eventually(func() error {
		_, err := locket.NewClientSkipCertVerify(
			locketLogger,
			locket.ClientLocketConfig{
				LocketAddress:        locketAddress,
				LocketCACertFile:     path.Join(fixturesPath, "locket-server.cert.pem"),
				LocketClientCertFile: path.Join(fixturesPath, "locket-client.cert.pem"),
				LocketClientKeyFile:  path.Join(fixturesPath, "locket-client.key.pem"),
			},
		)
		return err
	}).Should(BeNil())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	locketSession.Kill()
})
