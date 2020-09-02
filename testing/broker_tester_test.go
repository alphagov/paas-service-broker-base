package testing_test

import (
	"github.com/alphagov/paas-service-broker-base/testing"
	"github.com/alphagov/paas-service-broker-base/testing/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/fake_http_handler.go net/http.Handler

var _ = Describe("Broker Tester", func() {
	Context("makes HTTP requests to the in-memory broker API", func() {
		var credentials brokerapi.BrokerCredentials
		var brokerAPI fakes.FakeHandler
		var brokerTester testing.BrokerTester

		BeforeEach(func() {
			credentials = brokerapi.BrokerCredentials{Username: "user", Password: "password"}
			brokerAPI = fakes.FakeHandler{}
			brokerTester = testing.New(credentials, &brokerAPI)
		})

		It("'Services' makes a GET request to the right place", func() {
			brokerTester.Services()

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("GET"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/catalog"))
		})

		It("'Provision' makes a PUT request to the right place", func() {
			brokerTester.Provision("instance_id", testing.RequestBody{}, false)

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("PUT"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id"))
		})

		It("'Deprovision' makes a DELETE request to the right place", func() {
			brokerTester.Deprovision("instance_id", "service_id", "plan_id", false)

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("DELETE"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id"))
		})

		It("'Bind' makes a PUT request to the right place", func() {
			brokerTester.Bind("instance_id", "binding_id", testing.RequestBody{}, false)

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("PUT"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id/service_bindings/binding_id"))
		})

		It("'Unbind' makes a DELETE request to the right place", func() {
			brokerTester.Unbind("instance_id", "service_id", "plan_id", "binding_id", false)

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("DELETE"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id/service_bindings/binding_id"))
		})

		It("'Update' makes a PATCH request to the right place", func() {
			brokerTester.Update("instance_id", testing.RequestBody{}, false)

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("PATCH"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id"))
		})

		It("'LastOperation' makes a GET request to the right place", func() {
			brokerTester.LastOperation("instance_id", "service_id", "plan_id", "operation")

			Expect(brokerAPI.ServeHTTPCallCount()).To(Equal(1))
			Expect(getRequestMethod(&brokerAPI)).To(Equal("GET"))
			Expect(getRequestPath(&brokerAPI)).To(Equal("/v2/service_instances/instance_id/last_operation"))
		})
	})
})

func getRequestMethod(handler *fakes.FakeHandler) string {
	_, req := handler.ServeHTTPArgsForCall(0)
	return req.Method
}

func getRequestPath(handler *fakes.FakeHandler) string {
	_, req := handler.ServeHTTPArgsForCall(0)
	return req.URL.Path
}
