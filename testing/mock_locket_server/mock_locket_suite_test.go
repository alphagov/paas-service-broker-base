package mock_locket_server

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMockLocket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock Locket Suite")
}
