# paas-service-broker-base

This repository contains two GoLang libraries:
 
1. A basic framework for creating service brokers on GOV.UK PaaS
2. A test harness for testing brokers

## Using the framework library
The framework library provides 3 key interfaces and components

`provider.ServiceProvider`. An interface defining the methods expected of a service broker, as defined in the 
[Open Service Broker API spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md). Service brokers 
implement the interface to integrate with the rest of `paas-service-broker-base`

`broker.New(Config, provider.ServiceProvider, lager.Logger) *Broker`. Creates a new `broker.Broker`, which takes care
of the boilerplate code for logging, and checking that services and plans exist, and then delegates to the  `ServiceProvider` 
implementation.

`broker.NewAPI(brokerapi.ServiceBroker, lager.Logger, Config) http.Handler`. Creates a new `http.Handler` which listens
for all the paths defined in the [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md) and 
routes them to the correct action on the `broker.Broker` instance.

Combining these 3 parts, a consumer will have a `main.go` file a little like

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	
	"code.cloudfoundry.org/lager/v3"
	"github.com/alphagov/paas-my-service-broker/"
	"github.com/alphagov/paas-service-broker-base/broker"
)

func main() {
	// Read the configuration file
	cfg := readConfigFile()
	
	// Create a new logger
	logger := lager.NewLogger("my-service-broker")
	logger.RegisterSink(lager.NewWriterSin(os.StdOut, cfg.API.LogLevel))
	
	// Create a new service broker
	myProvider := NewProviderInstance()
	serviceBroker := broker.New(cfg, myProvider, logger)
	brokerAPI := broker.NewAPI(serviceBroker, logger, cfg)
	
	// Create an HTTP listener
	listener, err := net.Listen("tcp", ":" + cfg.API.Port)
	if err != nil {
		log.Fatalf("Error listening on port %s: %s, config.API.Port, err)
	}
	
	// Begin listening
	fmt.Println("MyServiceBroker started on port " + config.API.Port + "...")
	http.Serve(listener, brokerAPI)
}
```

## Using broker tester
The broker tester provides a test harness for testing a broker via its HTTP interface, using an in-memory instance of 
the broker.

It has a single major function `testing.New(brokerapi.BrokerCredentials, http.Handler) BrokerTester` which returns a instance of `BrokerTester`, whose API closely mimics that of the Open Service Broker API spec. `BrokerTester` takes care of marshalling the method arguments in to an HTTP request, which is handled by an in-memory service broker API.

To create a broker test harness, consumers create a broker API instance as above, and pass it to `testing.New`. See the example test case below.

```go

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"github.com/alphagov/paas-service-broker-base/broker"
	brokertesting "github.com/alphagov/paas-service-broker-base/testing"
	
	"github.com/pivotal-cf/brokerapi"
)

const ASYNC_ALLOWED = true

var _ = Describe("test suite", function() {
	It("provisions stuff", func(){
		brokerTester := init()
		
		res := brokerTester.Provision("instance-id", brokertesting.RequestBody{
			ServiceId: "service",
			PlanID: "plan",
		}, ASYNC_ALLOWED)
		
		defer func(){
			res := brokerTester.Deprovision("instance-id", "service", "plan", ASYNC_ALLOWED)
			Expect(res.Code).To(Equal(http.StatusOK))
		}()
		
		Expect(res.Code).To(Equal(http.StatusCreated))
	})
})

func init() brokertesting.BrokerTester {
	cfg := readTestConfigFile()
	
	logger := lager.NewLogger("my-service-broker")
	logger.RegisterSink(lager.NewWriterSinl(GinkgoWriter, cfg.API.LagerLogLevel))
	
	myProvider := NewMyProvider()
	serviceBroker := broker.New(cfg, myProvider, logger)
	brokerAPI := broker.NewAPI(serviceBroker, logger, cfg)
	
	credentials := brokerapi.BrokerCredentials{
	  Username: cfg.BasicAuthUsername,
	  Password: cfg.BasicAuthPassword
	}
	return brokertesting.New(credentials, brokerAPI)
}
```
