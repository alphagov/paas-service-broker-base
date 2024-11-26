package broker_test

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"sync"
	"time"

	"code.cloudfoundry.org/lager/v3"
	. "github.com/alphagov/paas-service-broker-base/broker"
	"github.com/alphagov/paas-service-broker-base/provider"
	"github.com/alphagov/paas-service-broker-base/provider/fakes"
	"github.com/pivotal-cf/brokerapi/domain"
	"github.com/pivotal-cf/brokerapi/domain/apiresponses"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Broker", func() {

	const asyncAllowed = true

	var (
		validConfig Config
		instanceID  string
		orgGUID     string
		spaceGUID   string
		plan1       domain.ServicePlan
		plan2       domain.ServicePlan
		service1    domain.Service
	)

	BeforeEach(func() {
		instanceID = "instanceID"
		orgGUID = "org-guid"
		spaceGUID = "space-guid"
		plan1 = domain.ServicePlan{
			ID:   "plan1",
			Name: "plan1",
		}
		plan2 = domain.ServicePlan{
			ID:   "plan2",
			Name: "plan2",
		}
		service1 = domain.Service{
			ID:            "service1",
			Name:          "service1",
			PlanUpdatable: true,
			Plans:         []domain.ServicePlan{plan1, plan2},
		}
		validConfig = Config{
			Catalog: Catalog{
				apiresponses.CatalogResponse{
					Services: []domain.Service{service1},
				},
			},
			API: API{
				Locket: &LocketConfig{
					Address:        mockLocket.ListenAddress,
					CACertFile:     path.Join(locketFixtures.Filepath, "locket-server.cert.pem"),
					ClientCertFile: path.Join(locketFixtures.Filepath, "locket-client.cert.pem"),
					ClientKeyFile:  path.Join(locketFixtures.Filepath, "locket-client.key.pem"),
					SkipVerify:     true,
					RetryInterval:  time.Millisecond * 1,
				},
			},
		}
	})

	Describe("Provision", func() {
		var validProvisionDetails domain.ProvisionDetails

		BeforeEach(func() {
			validProvisionDetails = domain.ProvisionDetails{
				ServiceID:        service1.ID,
				PlanID:           plan1.ID,
				OrganizationGUID: orgGUID,
				SpaceGUID:        spaceGUID,
			}
		})

		It("logs a debug message when provision begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(log).To(gbytes.Say("provision-start"))
		})

		It("errors if async isn't allowed", func() {
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			asyncAllowed := false

			_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, asyncAllowed)

			Expect(err).To(Equal(apiresponses.ErrAsyncRequired))
		})

		It("errors if the service is not in the catalog", func() {
			config := validConfig
			config.Catalog = Catalog{Catalog: apiresponses.CatalogResponse{}}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(err).To(MatchError("Error: service " + service1.ID + " not found in the catalog"))
		})

		It("errors if the plan is not in the catalog", func() {
			config := validConfig
			config.Catalog.Catalog.Services[0].Plans = []domain.ServicePlan{}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(err).To(MatchError("Error: plan " + plan1.ID + " not found in service " + service1.ID))
		})

		It("sets a deadline by which the provision request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(fakeProvider.ProvisionCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.ProvisionArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(fakeProvider.ProvisionCallCount()).To(Equal(1))
			_, provisionData := fakeProvider.ProvisionArgsForCall(0)

			expectedProvisionData := provider.ProvisionData{
				InstanceID: instanceID,
				Details:    validProvisionDetails,
				Service:    service1,
				Plan:       plan1,
			}

			Expect(provisionData).To(Equal(expectedProvisionData))
		})

		It("errors if provisioning fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.ProvisionReturns(nil, errors.New("ERROR PROVISIONING"))

			_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(err).To(MatchError("ERROR PROVISIONING"))
		})

		Context("when provider does not implement async", func() {
			It("does not allow returning async true", func() {
				fakeProvider := &fakes.FakeProvisioner{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).NotTo(HaveOccurred())
				fakeProvider.ProvisionReturns(&domain.ProvisionedServiceSpec{
					IsAsync: true,
				}, nil)

				_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

				Expect(err).To(MatchError(ErrAsyncProvisionNotImplemented))
			})
		})

		Context("when provider does not implement provisioning", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when provisioning succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.ProvisionReturns(&domain.ProvisionedServiceSpec{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(log).To(gbytes.Say("provision-success"))
		})

		It("returns the provisioned service spec", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.ProvisionReturns(&domain.ProvisionedServiceSpec{
				DashboardURL:  "dashboard URL",
				OperationData: "operation data",
				IsAsync:       true,
			}, nil)

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			Expect(b.Provision(context.Background(), instanceID, validProvisionDetails, true)).
				To(Equal(domain.ProvisionedServiceSpec{
					IsAsync:       true,
					DashboardURL:  "dashboard URL",
					OperationData: "operation data",
				}))
		})

		It("gets a lock and releases it once it's created", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			locket := &fakes.FakeLocketClient{}

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			b.LocketClient = locket

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)
			Expect(locket.LockCallCount()).To(Equal(1))
			_, lockCallOne, _ := locket.LockArgsForCall(0)
			Expect(lockCallOne.Resource.Key).To(ContainSubstring("broker/instanceID"))

			Expect(locket.ReleaseCallCount()).To(Equal(1))
			_, releaseReqOne, _ := locket.ReleaseArgsForCall(0)
			Expect(releaseReqOne.Resource.Key).To(Equal(lockCallOne.Resource.Key))
		})

		It("waits for a lock and releases it once it's created", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			locket := &fakes.FakeLocketClient{}
			locket.LockReturnsOnCall(0, nil, status.Errorf(codes.AlreadyExists, "lock-collision"))
			locket.LockReturnsOnCall(1, nil, nil)

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			b.LocketClient = locket

			_, _ = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(locket.LockCallCount()).To(Equal(2))

			_, lockCallOne, _ := locket.LockArgsForCall(0)
			Expect(lockCallOne.Resource.Key).To(ContainSubstring("broker/instanceID"))

			_, lockCallTwo, _ := locket.LockArgsForCall(1)
			Expect(lockCallTwo.Resource.Key).To(ContainSubstring("broker/instanceID"))

			Expect(locket.ReleaseCallCount()).To(Equal(1))
			_, releaseReqOne, _ := locket.ReleaseArgsForCall(0)
			Expect(releaseReqOne.Resource.Key).To(Equal(lockCallOne.Resource.Key))
		})

		It("fails after waiting for many locks", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			locket := &fakes.FakeLocketClient{}
			locket.LockReturns(nil, status.Errorf(codes.AlreadyExists, "lock-collision"))

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			b.LocketClient = locket

			_, err = b.Provision(context.Background(), instanceID, validProvisionDetails, true)

			Expect(err).To(HaveOccurred())

			Expect(locket.LockCallCount()).To(Equal(31))

			_, lockCallOne, _ := locket.LockArgsForCall(0)
			Expect(lockCallOne.Resource.Key).To(ContainSubstring("broker/instanceID"))

			Expect(locket.ReleaseCallCount()).To(Equal(0))
		})
	})

	Describe("Deprovision", func() {
		var validDeprovisionDetails domain.DeprovisionDetails

		BeforeEach(func() {
			validDeprovisionDetails = domain.DeprovisionDetails{
				ServiceID: service1.ID,
				PlanID:    plan1.ID,
			}
		})

		It("logs a debug message when deprovision begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(log).To(gbytes.Say("deprovision-start"))
		})

		It("errors if async isn't allowed", func() {
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			asyncAllowed := false

			_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, asyncAllowed)

			Expect(err).To(Equal(apiresponses.ErrAsyncRequired))
		})

		It("errors if the service is not in the catalog", func() {
			config := validConfig
			config.Catalog = Catalog{Catalog: apiresponses.CatalogResponse{}}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(err).To(MatchError("Error: service " + service1.ID + " not found in the catalog"))
		})

		It("errors if the plan is not in the catalog", func() {
			config := validConfig
			config.Catalog.Catalog.Services[0].Plans = []domain.ServicePlan{}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(err).To(MatchError("Error: plan " + plan1.ID + " not found in service " + service1.ID))
		})

		It("sets a deadline by which the deprovision request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(fakeProvider.DeprovisionCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.DeprovisionArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(fakeProvider.DeprovisionCallCount()).To(Equal(1))
			_, deprovisionData := fakeProvider.DeprovisionArgsForCall(0)

			expectedDeprovisionData := provider.DeprovisionData{
				InstanceID: instanceID,
				Service:    service1,
				Plan:       plan1,
				Details:    validDeprovisionDetails,
			}

			Expect(deprovisionData).To(Equal(expectedDeprovisionData))
		})

		It("errors if deprovisioning fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.DeprovisionReturns(nil, errors.New("ERROR DEPROVISIONING"))

			_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(err).To(MatchError("ERROR DEPROVISIONING"))
		})

		Context("when provider does not implement async", func() {
			It("does not allow returning async true", func() {
				fakeProvider := &fakes.FakeProvisioner{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).NotTo(HaveOccurred())
				fakeProvider.DeprovisionReturns(&domain.DeprovisionServiceSpec{
					IsAsync: true,
				}, nil)

				_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

				Expect(err).To(MatchError(ErrAsyncProvisionNotImplemented))
			})
		})

		Context("when provider does not implement provisioning", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when deprovisioning succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.DeprovisionReturns(&domain.DeprovisionServiceSpec{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)

			Expect(log).To(gbytes.Say("deprovision-success"))
		})

		It("returns the deprovisioned service spec", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.DeprovisionReturns(&domain.DeprovisionServiceSpec{
				OperationData: "operation data",
				IsAsync:       true,
			}, nil)
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			Expect(b.Deprovision(context.Background(), instanceID, validDeprovisionDetails, true)).
				To(Equal(domain.DeprovisionServiceSpec{
					IsAsync:       true,
					OperationData: "operation data",
				}))
		})
	})

	Describe("Bind", func() {
		var (
			bindingID        string
			appGUID          string
			bindResource     *domain.BindResource
			validBindDetails domain.BindDetails
		)

		BeforeEach(func() {
			bindingID = "bindingID"
			appGUID = "appGUID"
			bindResource = &domain.BindResource{
				AppGuid: appGUID,
			}
			validBindDetails = domain.BindDetails{
				AppGUID:      appGUID,
				PlanID:       plan1.ID,
				ServiceID:    service1.ID,
				BindResource: bindResource,
			}
		})

		It("logs a debug message when binding begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(log).To(gbytes.Say("binding-start"))
		})

		It("sets a deadline by which the binding request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(fakeProvider.BindCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.BindArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(fakeProvider.BindCallCount()).To(Equal(1))
			_, bindData := fakeProvider.BindArgsForCall(0)

			expectedBindData := provider.BindData{
				InstanceID:   instanceID,
				BindingID:    bindingID,
				Details:      validBindDetails,
				AsyncAllowed: asyncAllowed,
			}

			Expect(bindData).To(Equal(expectedBindData))
		})

		It("errors if binding fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.BindReturns(nil, errors.New("ERROR BINDING"))

			_, err = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(err).To(MatchError("ERROR BINDING"))
		})

		Context("when provider does not implement async binding", func() {
			It("does not allow returning async true", func() {
				fakeProvider := &fakes.FakeBinder{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).NotTo(HaveOccurred())
				fakeProvider.BindReturns(&domain.Binding{
					IsAsync: true,
				}, nil)

				_, err = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

				Expect(err).To(MatchError(ErrAsyncBindNotImplemented))
			})
		})

		Context("when provider does not implement binding", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when binding succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.BindReturns(&domain.Binding{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(log).To(gbytes.Say("binding-success"))
		})

		It("returns the binding", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.BindReturns(&domain.Binding{
				Credentials: "some-value-of-interface{}-type",
				IsAsync:     true,
			}, nil)
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			Expect(b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)).
				To(Equal(domain.Binding{
					Credentials: "some-value-of-interface{}-type",
					IsAsync:     true,
				}))
		})

		It("gets a lock and releases it at the end", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			locket := &fakes.FakeLocketClient{}

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			b.LocketClient = locket

			_, _ = b.Bind(context.Background(), instanceID, bindingID, validBindDetails, true)

			Expect(locket.LockCallCount()).To(Equal(1))
			_, lockCallOne, _ := locket.LockArgsForCall(0)
			Expect(lockCallOne.Resource.Key).To(ContainSubstring("broker/instanceID"))

			Expect(locket.ReleaseCallCount()).To(Equal(1))
			_, releaseReqOne, _ := locket.ReleaseArgsForCall(0)
			Expect(releaseReqOne.Resource.Key).To(Equal(lockCallOne.Resource.Key))
		})
	})

	Describe("GetBinding", func() {
		var (
			bindingID string
		)

		BeforeEach(func() {
			bindingID = "bindingID"
		})

		It("logs a debug message when get binding begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.GetBinding(context.Background(), instanceID, bindingID)

			Expect(log).To(gbytes.Say("get-binding-start"))
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.GetBinding(context.Background(), instanceID, bindingID)

			Expect(fakeProvider.GetBindingCallCount()).To(Equal(1))
			_, data := fakeProvider.GetBindingArgsForCall(0)

			expectedData := provider.GetBindData{
				InstanceID: instanceID,
				BindingID:  bindingID,
			}

			Expect(data).To(Equal(expectedData))
		})

		It("errors if get binding fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.GetBindingReturns(nil, errors.New("ERROR BINDING"))
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.GetBinding(context.Background(), instanceID, bindingID)

			Expect(err).To(MatchError("ERROR BINDING"))
		})

		Context("when provider does not implement async binding", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.GetBinding(context.Background(), instanceID, bindingID)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when binding succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.GetBindingReturns(&domain.GetBindingSpec{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.GetBinding(context.Background(), instanceID, bindingID)

			Expect(log).To(gbytes.Say("get-binding-success"))
		})

		It("returns the binding", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.GetBindingReturns(&domain.GetBindingSpec{
				Credentials: "some-value-of-interface{}-type",
			}, nil)
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			Expect(b.GetBinding(context.Background(), instanceID, bindingID)).To(Equal(domain.GetBindingSpec{
				Credentials: "some-value-of-interface{}-type",
			}))
		})
	})

	Describe("Unbind", func() {
		var (
			bindingID          string
			validUnbindDetails domain.UnbindDetails
		)

		BeforeEach(func() {
			bindingID = "bindingID"
			validUnbindDetails = domain.UnbindDetails{
				PlanID:    plan1.ID,
				ServiceID: service1.ID,
			}
		})

		It("logs a debug message when unbinding begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(log).To(gbytes.Say("unbinding-start"))
		})

		It("sets a deadline by which the unbinding request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(fakeProvider.UnbindCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.UnbindArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(fakeProvider.UnbindCallCount()).To(Equal(1))
			_, unbindData := fakeProvider.UnbindArgsForCall(0)

			expectedUnbindData := provider.UnbindData{
				InstanceID:   instanceID,
				BindingID:    bindingID,
				Details:      validUnbindDetails,
				AsyncAllowed: asyncAllowed,
			}

			Expect(unbindData).To(Equal(expectedUnbindData))
		})

		It("errors if unbinding fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.UnbindReturns(nil, errors.New("ERROR UNBINDING"))

			_, err = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(err).To(MatchError("ERROR UNBINDING"))
		})

		Context("when provider does not implement async binding", func() {
			It("does not allow returning async true", func() {
				fakeProvider := &fakes.FakeBinder{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).NotTo(HaveOccurred())
				fakeProvider.UnbindReturns(&domain.UnbindSpec{
					IsAsync: true,
				}, nil)

				_, err = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

				Expect(err).To(MatchError(ErrAsyncBindNotImplemented))
			})
		})

		Context("when provider does not implement binding", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when unbinding succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.UnbindReturns(&domain.UnbindSpec{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(log).To(gbytes.Say("unbinding-success"))
		})

		It("gets a lock and releases it at the end", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			locket := &fakes.FakeLocketClient{}

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			b.LocketClient = locket

			_, _ = b.Unbind(context.Background(), instanceID, bindingID, validUnbindDetails, true)

			Expect(locket.LockCallCount()).To(Equal(1))
			_, lockCallOne, _ := locket.LockArgsForCall(0)
			Expect(lockCallOne.Resource.Key).To(ContainSubstring("broker/instanceID"))

			Expect(locket.ReleaseCallCount()).To(Equal(1))
			_, releaseReqOne, _ := locket.ReleaseArgsForCall(0)
			Expect(releaseReqOne.Resource.Key).To(Equal(lockCallOne.Resource.Key))
		})
	})

	Describe("Update", func() {
		var updatePlanDetails domain.UpdateDetails

		BeforeEach(func() {
			updatePlanDetails = domain.UpdateDetails{
				ServiceID: service1.ID,
				PlanID:    plan2.ID,
				PreviousValues: domain.PreviousValues{
					ServiceID: service1.ID,
					PlanID:    plan1.ID,
					OrgID:     orgGUID,
					SpaceID:   spaceGUID,
				},
			}
		})

		Describe("Updatability", func() {
			Context("when the plan is not updatable", func() {
				var updateParametersDetails domain.UpdateDetails

				BeforeEach(func() {
					validConfig.Catalog.Catalog.Services[0].PlanUpdatable = false

					updateParametersDetails = domain.UpdateDetails{
						ServiceID:     service1.ID,
						PlanID:        plan1.ID,
						RawParameters: json.RawMessage(`{"new":"parameter"}`),
						PreviousValues: domain.PreviousValues{
							ServiceID: service1.ID,
							PlanID:    plan1.ID,
							OrgID:     orgGUID,
							SpaceID:   spaceGUID,
						},
					}
				})

				It("returns an error when changing the plan", func() {
					b, err := New(validConfig, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
					Expect(err).NotTo(HaveOccurred())

					Expect(updatePlanDetails.PlanID).NotTo(Equal(updatePlanDetails.PreviousValues.PlanID))
					_, err = b.Update(context.Background(), instanceID, updatePlanDetails, true)

					Expect(err).To(Equal(apiresponses.ErrPlanChangeNotSupported))
				})

				It("accepts the update request when just changing parameters", func() {
					fakeProvider := &fakes.FakeAsyncProvider{}
					fakeProvider.UpdateReturns(&domain.UpdateServiceSpec{}, nil)
					b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
					Expect(err).NotTo(HaveOccurred())

					Expect(updateParametersDetails.PlanID).To(Equal(updateParametersDetails.PreviousValues.PlanID))
					_, err = b.Update(context.Background(), instanceID, updateParametersDetails, true)

					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		It("logs a debug message when update begins", func() {
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(log).To(gbytes.Say("update-start"))
		})

		It("errors if async isn't allowed", func() {
			b, err := New(validConfig, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			asyncAllowed := false

			_, err = b.Update(context.Background(), instanceID, updatePlanDetails, asyncAllowed)

			Expect(err).To(Equal(apiresponses.ErrAsyncRequired))
		})

		It("errors if the service is not in the catalog", func() {
			config := validConfig
			config.Catalog = Catalog{Catalog: apiresponses.CatalogResponse{}}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(err).To(MatchError("Error: service " + service1.ID + " not found in the catalog"))
		})

		It("errors if the plan is not in the catalog", func() {
			config := validConfig
			config.Catalog.Catalog.Services[0].Plans = []domain.ServicePlan{}
			b, err := New(config, &fakes.FakeAsyncProvider{}, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, err = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(err).To(MatchError("Error: plan " + plan2.ID + " not found in service " + service1.ID))
		})

		It("sets a deadline by which the update request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(fakeProvider.UpdateCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.UpdateArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(fakeProvider.UpdateCallCount()).To(Equal(1))
			_, updateData := fakeProvider.UpdateArgsForCall(0)

			expectedUpdateData := provider.UpdateData{
				InstanceID: instanceID,
				Details:    updatePlanDetails,
				Service:    service1,
				Plan:       plan2,
			}

			Expect(updateData).To(Equal(expectedUpdateData))
		})

		It("errors if update fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.UpdateReturns(nil, errors.New("ERROR UPDATING"))

			_, err = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(err).To(MatchError("ERROR UPDATING"))
		})

		Context("when provider does not implement updating", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.Update(context.Background(), instanceID, updatePlanDetails, true)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when updating succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.UpdateReturns(&domain.UpdateServiceSpec{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.Update(context.Background(), instanceID, updatePlanDetails, true)

			Expect(log).To(gbytes.Say("update-success"))
		})

		It("returns the update service spec", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.UpdateReturns(&domain.UpdateServiceSpec{
				OperationData: "operation data",
				IsAsync:       true,
			}, nil)

			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Update(context.Background(), instanceID, updatePlanDetails, true)).
				To(Equal(domain.UpdateServiceSpec{
					IsAsync:       true,
					OperationData: "operation data",
				}))
		})
	})

	Describe("LastOperation", func() {
		var pollDetails domain.PollDetails

		BeforeEach(func() {
			pollDetails = domain.PollDetails{
				OperationData: `{"operation_type": "provision"}`,
			}
		})

		It("logs a debug message when the last operation check begins", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.LastOperationReturns(&domain.LastOperation{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastOperation(context.Background(), instanceID, pollDetails)

			Expect(log).To(gbytes.Say("last-operation-start"))
		})

		It("sets a deadline by which the last operation request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastOperation(context.Background(), instanceID, pollDetails)

			Expect(fakeProvider.LastOperationCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.LastOperationArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastOperation(context.Background(), instanceID, pollDetails)

			Expect(fakeProvider.LastOperationCallCount()).To(Equal(1))
			_, lastOperationData := fakeProvider.LastOperationArgsForCall(0)

			expectedLastOperationData := provider.LastOperationData{
				InstanceID:  instanceID,
				PollDetails: pollDetails,
			}

			Expect(lastOperationData).To(Equal(expectedLastOperationData))
		})

		It("errors if last operation fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.LastOperationReturns(nil, errors.New("ERROR LAST OPERATION"))

			_, err = b.LastOperation(context.Background(), instanceID, pollDetails)

			Expect(err).To(MatchError("ERROR LAST OPERATION"))
		})

		Context("when provider does not implement last operation", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.LastOperation(context.Background(), instanceID, pollDetails)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when last operation check succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.LastOperationReturns(&domain.LastOperation{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastOperation(context.Background(), instanceID, pollDetails)

			Expect(log).To(gbytes.Say("last-operation-success"))
		})

		It("returns the last operation status", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.LastOperationReturns(&domain.LastOperation{
				State:       domain.Succeeded,
				Description: "Provision successful",
			}, nil)

			Expect(b.LastOperation(context.Background(), instanceID, pollDetails)).
				To(Equal(domain.LastOperation{
					State:       domain.Succeeded,
					Description: "Provision successful",
				}))
		})
	})

	Describe("LastBindingOperation", func() {
		var pollDetails domain.PollDetails
		var bindingID string

		BeforeEach(func() {
			pollDetails = domain.PollDetails{
				OperationData: `{"operation_type": "provision"}`,
			}
			bindingID = "bindingID"
		})

		It("logs a debug message when the last binding operation check begins", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.LastBindingOperationReturns(&domain.LastOperation{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)

			Expect(log).To(gbytes.Say("last-binding-operation-start"))
		})

		It("sets a deadline by which the last operation request should complete", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)

			Expect(fakeProvider.LastBindingOperationCallCount()).To(Equal(1))
			receivedContext, _ := fakeProvider.LastBindingOperationArgsForCall(0)

			_, hasDeadline := receivedContext.Deadline()

			Expect(hasDeadline).To(BeTrue())
		})

		It("passes the correct data to the Provider", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)

			Expect(fakeProvider.LastBindingOperationCallCount()).To(Equal(1))
			_, lastBindingOperationData := fakeProvider.LastBindingOperationArgsForCall(0)

			expectedLastBindingOperationData := provider.LastBindingOperationData{
				InstanceID:  instanceID,
				BindingID:   bindingID,
				PollDetails: pollDetails,
			}

			Expect(lastBindingOperationData).To(Equal(expectedLastBindingOperationData))
		})

		It("errors if last binding operation fails", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.LastBindingOperationReturns(nil, errors.New("ERROR LAST OPERATION"))

			_, err = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)

			Expect(err).To(MatchError("ERROR LAST OPERATION"))
		})

		Context("when provider does not implement last operation", func() {
			It("will return error", func() {
				var fakeProvider interface{}
				b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
				Expect(err).ToNot(HaveOccurred())

				_, err = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)
				Expect(err).To(MatchError(ErrNotImplemented))
			})
		})

		It("logs a debug message when last operation check succeeds", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			fakeProvider.LastBindingOperationReturns(&domain.LastOperation{}, nil)
			logger := lager.NewLogger("broker")
			log := gbytes.NewBuffer()
			logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
			b, err := New(validConfig, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			_, _ = b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)

			Expect(log).To(gbytes.Say("last-binding-operation-success"))
		})

		It("returns the last operation status", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			b, err := New(validConfig, fakeProvider, lager.NewLogger("broker"))
			Expect(err).NotTo(HaveOccurred())
			fakeProvider.LastBindingOperationReturns(&domain.LastOperation{
				State:       domain.Succeeded,
				Description: "Provision successful",
			}, nil)

			Expect(b.LastBindingOperation(context.Background(), instanceID, bindingID, pollDetails)).
				To(Equal(domain.LastOperation{
					State:       domain.Succeeded,
					Description: "Provision successful",
				}))
		})
	})

	Describe("Locking", func() {

		var (
			configWithLocket    Config
			configWithoutLocket Config
		)

		BeforeEach(func() {
			configWithLocket = Config{
				Catalog: Catalog{
					apiresponses.CatalogResponse{
						Services: []domain.Service{service1},
					},
				},
				API: API{
					Locket: &LocketConfig{
						Address:        mockLocket.ListenAddress,
						CACertFile:     path.Join(locketFixtures.Filepath, "locket-server.cert.pem"),
						ClientCertFile: path.Join(locketFixtures.Filepath, "locket-client.cert.pem"),
						ClientKeyFile:  path.Join(locketFixtures.Filepath, "locket-client.key.pem"),
						SkipVerify:     true,
					},
				},
			}
			configWithoutLocket = Config{
				Catalog: Catalog{
					apiresponses.CatalogResponse{
						Services: []domain.Service{service1},
					},
				},
				API: API{},
			}
		})

		It("should use SimpleLock if no locket config provided", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			logger := lager.NewLogger("broker")
			b, err := New(configWithoutLocket, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.LocketClient).ToNot(BeNil())
			Expect(b.LocketClient).To(BeAssignableToTypeOf(&SimpleLock{}))
		})

		It("should not use SimpleLock if locket config provided", func() {
			fakeProvider := &fakes.FakeAsyncProvider{}
			logger := lager.NewLogger("broker")
			b, err := New(configWithLocket, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.LocketClient).ToNot(BeNil())
			Expect(b.LocketClient).NotTo(BeAssignableToTypeOf(&SimpleLock{}))
		})

		DescribeTable("should lock and unlock", func(cfg Config) {
			fakeProvider := &fakes.FakeAsyncProvider{}
			logger := lager.NewLogger("broker")
			b, err := New(cfg, fakeProvider, logger)
			Expect(err).NotTo(HaveOccurred())
			s := "original"
			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				defer wg.Done()
				By("g1 getting lock")
				lock, err := b.ObtainServiceLock(context.Background(), instanceID, 30)
				Expect(err).NotTo(HaveOccurred())
				defer b.ReleaseServiceLock(context.Background(), lock)
				By("g1 got lock")

				g1original := s
				s = "goroutine-1"
				time.Sleep(1 * time.Second)
				s = g1original

				By("g1 done")
			}()
			go func() {
				defer wg.Done()
				By("g2 getting lock")
				lock, err := b.ObtainServiceLock(context.Background(), instanceID, 30)
				Expect(err).NotTo(HaveOccurred())
				defer b.ReleaseServiceLock(context.Background(), lock)
				By("g2 got lock")

				g2original := s
				s = "goroutine-2"
				time.Sleep(1 * time.Second)
				s = g2original

				By("g2 done")
			}()

			wg.Wait()

			Expect(s).To(Equal("original"))

		},
			Entry("with a locket config", configWithLocket),
			Entry("without a locket config", configWithoutLocket),
		)

	})
})
