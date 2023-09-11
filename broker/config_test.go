package broker_test

import (
	"strings"

	"code.cloudfoundry.org/lager"
	. "github.com/alphagov/paas-service-broker-base/broker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var configSource string

	Describe("Mandatory fields", func() {
		It("requires a basic auth username", func() {
			configSource = `
				{
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: basic auth username required"))
		})

		It("requires a basic auth password", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: basic auth password required"))
		})

		It("requires a catalog", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: catalog required"))
		})

		It("requires a locket address if locket config provided", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: locket address required"))
		})

		It("does not require locket config", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).ToNot(HaveOccurred())
		})

		It("requires all lts fields if tls config provided", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"tls": {}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: TLS certificate required"))

		})
	})

	Describe("Log levels", func() {

		It("helps convert log level string values into lager.LogLevel values", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			config, err := NewConfig(strings.NewReader(configSource))
			Expect(err).NotTo(HaveOccurred())

			lagerLogLevel, err := config.API.ConvertLogLevel()
			Expect(err).NotTo(HaveOccurred())
			Expect(lagerLogLevel).To(Equal(lager.DEBUG))
		})

		It("errors if the log level doesn't map to a Lager log level", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"log_level": "debuggery",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: log level debuggery does not map to a Lager log level"))
		})
	})

	Describe("Default values", func() {
		It("sets a default port", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"log_level": "debug",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			config, err := NewConfig(strings.NewReader(configSource))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.API.Port).To(Equal(DefaultPort))
		})

		It("sets a default log_level", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			config, err := NewConfig(strings.NewReader(configSource))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.API.LogLevel).To(Equal(DefaultLogLevel))
		})

		It("sets a default context_timeout_seconds", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"port": "8080",
					"catalog": {"services": [{"name": "service1", "plans": [{"name": "plan1"}]}]},
					"locket": {"address": "my.locker.server"}
				}
			`
			config, err := NewConfig(strings.NewReader(configSource))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.API.ContextTimeout()).To(Equal(DefaultContextTimeout))
		})
	})

	Describe("Catalog", func() {
		It("requires at least one service", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"catalog": {"services": []},
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: at least one service is required"))
		})

		It("requires at least one plan", func() {
			configSource = `
				{
					"basic_auth_username":"username",
					"basic_auth_password":"1234",
					"catalog": {"services": [
						{"name": "service1", "plans": [{"name": "plan1"}]},
						{"name": "service2", "plans": []}
					]},
					"locket": {"address": "my.locker.server"}
				}
			`
			_, err := NewConfig(strings.NewReader(configSource))
			Expect(err).To(MatchError("Config error: no plans found for service service2"))
		})
	})
})
