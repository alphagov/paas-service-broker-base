module github.com/alphagov/paas-service-broker-base

require (
	code.cloudfoundry.org/diego-logging-client v0.0.0-20190918155030-cd01d2d2c629 // indirect
	code.cloudfoundry.org/lager v2.0.0+incompatible
	code.cloudfoundry.org/locket v0.0.0-20191127212858-571765e813ca
	code.cloudfoundry.org/tlsconfig v0.0.0-20191220232943-2819aba30e10 // indirect
	github.com/alphagov/paas-s3-broker v0.12.0 // indirect
	github.com/aws/aws-sdk-go v1.27.1 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190808214049-35bcce23fc5f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/hashicorp/consul/api v1.3.0 // indirect
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pivotal-cf/brokerapi v4.3.0+incompatible
	github.com/tedsuo/ifrit v0.0.0-20191009134036-9a97d0632f00
	golang.org/x/net v0.0.0-20190501004415-9ce7a6920f09
	google.golang.org/grpc v1.26.0 // indirect
)

go 1.13
