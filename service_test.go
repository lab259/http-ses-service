package sqssrv

import (
	"github.com/jamillosantos/macchiato"
	"github.com/lab259/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"testing"
	"github.com/aws/aws-sdk-go/service/ses"
)

func TestService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "Redigo Test Suite")
}

var _ = Describe("SQSService", func() {
	It("should fail loading a configuration", func() {
		var service SESService
		configuration, err := service.LoadConfiguration()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("not implemented"))
		Expect(configuration).To(BeNil())
	})

	It("should fail applying configuration", func() {
		var service SESService
		err := service.ApplyConfiguration(map[string]interface{}{
			"address": "localhost",
		})
		Expect(err).To(Equal(http.ErrWrongConfigurationInformed))
	})

	It("should apply the configuration using a pointer", func() {
		var service SESService
		err := service.ApplyConfiguration(&SESServiceConfiguration{
			Endpoint: "endpoint",
			Region:   "region",
			Secret:   "secret",
			Key:      "key",
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Region).To(Equal("region"))
		Expect(service.Configuration.Secret).To(Equal("secret"))
		Expect(service.Configuration.Key).To(Equal("key"))
	})

	It("should apply the configuration using a copy", func() {
		var service SESService
		err := service.ApplyConfiguration(SESServiceConfiguration{
			Endpoint: "endpoint",
			Region: "region",
			Secret: "secret",
			Key:    "key",
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Region).To(Equal("region"))
		Expect(service.Configuration.Endpoint).To(Equal("endpoint"))
		Expect(service.Configuration.Secret).To(Equal("secret"))
		Expect(service.Configuration.Key).To(Equal("key"))
	})

	validConfiguration := SESServiceConfiguration{
		Endpoint: "http://localhost:4576",
	}

	It("should start the service", func() {
		var service SESService
		Expect(service.ApplyConfiguration(&validConfiguration)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		Expect(service.RunWithSES(func(client *ses.SES) error {
			Expect(client).NotTo(BeNil())
			return nil
		})).To(BeNil())
	})

	It("should restart the service", func() {
		var service SESService
		Expect(service.ApplyConfiguration(&validConfiguration)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Restart()).To(BeNil())
		Expect(service.RunWithSES(func(client *ses.SES) error {
			Expect(client).NotTo(BeNil())
			return nil
		})).To(BeNil())
	})


	It("should fail try to run a command without starting the service", func() {
		var service SESService
		Expect(service.RunWithSES(nil)).To(Equal(http.ErrServiceNotRunning))
	})

})
