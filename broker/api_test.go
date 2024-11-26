package broker_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"

	"code.cloudfoundry.org/lager/v3"
	. "github.com/alphagov/paas-service-broker-base/broker"
	"github.com/alphagov/paas-service-broker-base/provider/fakes"
	broker_tester "github.com/alphagov/paas-service-broker-base/testing"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/domain"
	"github.com/pivotal-cf/brokerapi/domain/apiresponses"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Broker API", func() {
	var (
		err error

		instanceID  string
		orgGUID     string
		spaceGUID   string
		service1    string
		plan1       string
		validConfig Config

		username string
		password string

		logger       lager.Logger
		fakeProvider *fakes.FakeAsyncProvider

		broker       *Broker
		brokerAPI    http.Handler
		brokerTester broker_tester.BrokerTester
	)

	BeforeEach(func() {
		instanceID = "instanceID"
		orgGUID = "org-guid"
		spaceGUID = "space-guid"
		service1 = "service1"
		plan1 = "plan1"
		validConfig = Config{
			API: API{
				BasicAuthUsername: username,
				BasicAuthPassword: password,
				Locket: &LocketConfig{
					Address:        mockLocket.ListenAddress,
					CACertFile:     path.Join(locketFixtures.Filepath, "locket-server.cert.pem"),
					ClientCertFile: path.Join(locketFixtures.Filepath, "locket-client.cert.pem"),
					ClientKeyFile:  path.Join(locketFixtures.Filepath, "locket-client.key.pem"),
					SkipVerify:     true,
					RetryInterval:  time.Millisecond * 1,
				},
			},
			Catalog: Catalog{apiresponses.CatalogResponse{
				Services: []domain.Service{
					{
						ID:            service1,
						Name:          service1,
						PlanUpdatable: true,
						Plans: []domain.ServicePlan{
							{
								ID:   plan1,
								Name: plan1,
							},
						},
					},
				},
			},
			},
		}
		logger = lager.NewLogger("broker-api")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
		fakeProvider = &fakes.FakeAsyncProvider{}
		broker, err = New(validConfig, fakeProvider, logger)
		Expect(err).NotTo(HaveOccurred())
		brokerAPI = NewAPI(broker, logger, validConfig)

		brokerTester = broker_tester.New(brokerapi.BrokerCredentials{
			Username: validConfig.API.BasicAuthUsername,
			Password: validConfig.API.BasicAuthPassword,
		}, brokerAPI)
	})

	It("serves a healthcheck endpoint", func() {
		res := brokerTester.Get("/healthcheck", url.Values{})
		Expect(res.Code).To(Equal(http.StatusOK))
	})

	Describe("Services", func() {
		It("serves the catalog", func() {
			res := brokerTester.Services()
			Expect(res.Code).To(Equal(http.StatusOK))

			catalogResponse := apiresponses.CatalogResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &catalogResponse)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(catalogResponse.Services)).To(Equal(1))
			Expect(catalogResponse.Services[0].ID).To(Equal(service1))
			Expect(len(catalogResponse.Services[0].Plans)).To(Equal(1))
			Expect(catalogResponse.Services[0].Plans[0].ID).To(Equal(plan1))
		})
	})

	Describe("Provision", func() {
		It("accepts a provision request", func() {
			fakeProvider.ProvisionReturns(&domain.ProvisionedServiceSpec{
				DashboardURL:  "dashboardURL",
				OperationData: "operationData",
				IsAsync:       true,
			}, nil)
			res := brokerTester.Provision(
				instanceID,
				broker_tester.RequestBody{
					ServiceID:        service1,
					PlanID:           plan1,
					OrganizationGUID: orgGUID,
					SpaceGUID:        spaceGUID,
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			provisioningResponse := apiresponses.ProvisioningResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &provisioningResponse)
			Expect(err).NotTo(HaveOccurred())

			expectedResponse := apiresponses.ProvisioningResponse{
				DashboardURL:  "dashboardURL",
				OperationData: "operationData",
			}
			Expect(provisioningResponse).To(Equal(expectedResponse))
		})

		It("responds with an internal server error if the provider errors", func() {
			fakeProvider.ProvisionReturns(&domain.ProvisionedServiceSpec{}, errors.New("some provisioning error"))
			res := brokerTester.Provision(
				instanceID,
				broker_tester.RequestBody{
					ServiceID:        service1,
					PlanID:           plan1,
					OrganizationGUID: orgGUID,
					SpaceGUID:        spaceGUID,
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
		})

		It("rejects requests for synchronous provisioning", func() {
			res := brokerTester.Provision(
				instanceID,
				broker_tester.RequestBody{
					ServiceID:        service1,
					PlanID:           plan1,
					OrganizationGUID: orgGUID,
					SpaceGUID:        spaceGUID,
				},
				false,
			)
			Expect(res.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	Describe("Deprovision", func() {
		It("accepts a deprovision request", func() {
			fakeProvider.DeprovisionReturns(&domain.DeprovisionServiceSpec{
				OperationData: "operationData",
				IsAsync:       true,
			}, nil)
			res := brokerTester.Deprovision(instanceID, service1, plan1, true)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			deprovisionResponse := apiresponses.DeprovisionResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &deprovisionResponse)
			Expect(err).NotTo(HaveOccurred())

			expectedResponse := apiresponses.DeprovisionResponse{
				OperationData: "operationData",
			}
			Expect(deprovisionResponse).To(Equal(expectedResponse))
		})

		It("responds with an internal server error if the provider errors", func() {
			fakeProvider.DeprovisionReturns(nil, errors.New("some deprovisioning error"))
			res := brokerTester.Deprovision(instanceID, service1, plan1, true)
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
		})

		It("rejects requests for synchronous deprovisioning", func() {
			res := brokerTester.Deprovision(instanceID, service1, plan1, false)
			Expect(res.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	Describe("Bind", func() {
		var (
			bindingID string
			appGUID   string
		)

		BeforeEach(func() {
			bindingID = "bindingID"
			appGUID = "appGUID"
		})

		It("creates a binding", func() {
			fakeProvider.BindReturns(&domain.Binding{Credentials: "secrets"}, nil)
			res := brokerTester.Bind(
				instanceID,
				bindingID,
				broker_tester.RequestBody{
					ServiceID: service1,
					PlanID:    plan1,
					AppGUID:   appGUID,
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusCreated))

			binding := domain.Binding{}
			err := json.Unmarshal(res.Body.Bytes(), &binding)
			Expect(err).NotTo(HaveOccurred())

			expectedBinding := domain.Binding{
				Credentials: "secrets",
			}
			Expect(binding).To(Equal(expectedBinding))
		})

		It("responds with an internal server error if the provider errors", func() {
			fakeProvider.BindReturns(nil, errors.New("some binding error"))
			res := brokerTester.Bind(
				instanceID,
				bindingID,
				broker_tester.RequestBody{
					ServiceID: service1,
					PlanID:    plan1,
					AppGUID:   appGUID,
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("Unbind", func() {
		var bindingID string

		BeforeEach(func() {
			bindingID = "bindingID"
		})

		It("unbinds", func() {
			fakeProvider.UnbindReturns(&domain.UnbindSpec{
				OperationData: "operation data",
			}, nil)
			res := brokerTester.Unbind(
				instanceID,
				service1,
				plan1,
				bindingID,
				true,
			)
			Expect(res.Code).To(Equal(http.StatusOK))
		})

		It("responds with an internal server error if the provider errors", func() {
			fakeProvider.UnbindReturns(nil, errors.New("some unbinding error"))
			res := brokerTester.Unbind(
				instanceID,
				service1,
				plan1,
				bindingID,
				true,
			)
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("Update", func() {
		It("accepts an update request", func() {
			fakeProvider.UpdateReturns(&domain.UpdateServiceSpec{
				OperationData: "operationData",
				IsAsync:       true,
			}, nil)
			res := brokerTester.Update(
				instanceID,
				broker_tester.RequestBody{
					ServiceID: service1,
					PlanID:    plan1,
					PreviousValues: &broker_tester.RequestBody{
						PlanID: "plan2",
					},
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			updateResponse := apiresponses.UpdateResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &updateResponse)
			Expect(err).NotTo(HaveOccurred())

			expectedResponse := apiresponses.UpdateResponse{
				OperationData: "operationData",
			}
			Expect(updateResponse).To(Equal(expectedResponse))
		})

		It("responds with an internal server error if the provider errors", func() {
			fakeProvider.UpdateReturns(nil, errors.New("some update error"))
			res := brokerTester.Update(
				instanceID,
				broker_tester.RequestBody{
					ServiceID: service1,
					PlanID:    plan1,
					PreviousValues: &broker_tester.RequestBody{
						PlanID: "plan2",
					},
				},
				true,
			)
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
		})

		It("rejects requests for synchronous updating", func() {
			res := brokerTester.Update(
				instanceID,
				broker_tester.RequestBody{
					ServiceID: service1,
					PlanID:    plan1,
					PreviousValues: &broker_tester.RequestBody{
						PlanID: "plan2",
					},
				},
				false,
			)
			Expect(res.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	Describe("LastOperation", func() {
		It("provides the state of the operation", func() {
			fakeProvider.LastOperationReturns(&domain.LastOperation{
				State:       domain.Succeeded,
				Description: "description",
			}, nil)
			res := brokerTester.LastOperation(instanceID, "", "", "")
			Expect(res.Code).To(Equal(http.StatusOK))

			lastOperationResponse := apiresponses.LastOperationResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &lastOperationResponse)
			Expect(err).NotTo(HaveOccurred())

			expectedResponse := apiresponses.LastOperationResponse{
				State:       domain.Succeeded,
				Description: "description",
			}
			Expect(lastOperationResponse).To(Equal(expectedResponse))
		})

		It("responds with an internal server error if the provider errors", func() {
			lastOperationError := errors.New("some last operation error")
			fakeProvider.LastOperationReturns(nil, lastOperationError)
			res := brokerTester.LastOperation(instanceID, "", "", "")
			Expect(res.Code).To(Equal(http.StatusInternalServerError))

			lastOperationResponse := apiresponses.LastOperationResponse{}
			err := json.Unmarshal(res.Body.Bytes(), &lastOperationResponse)
			Expect(err).NotTo(HaveOccurred())

			expectedResponse := apiresponses.LastOperationResponse{
				State:       "",
				Description: lastOperationError.Error(),
			}
			Expect(lastOperationResponse).To(Equal(expectedResponse))
		})
	})
})
