package sqssrv

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/lab259/http"
)

// SESServiceConfiguration is the configuration for the `SESService`
type SESServiceConfiguration struct {
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

type credentialsFromStruct struct {
	credentials *SESServiceConfiguration
}

func NewCredentialsFromStruct(credentials *SESServiceConfiguration) *credentialsFromStruct {
	return &credentialsFromStruct{
		credentials: credentials,
	}
}

func (c *credentialsFromStruct) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     c.credentials.Key,
		SecretAccessKey: c.credentials.Secret,
	}, nil
}

func (*credentialsFromStruct) IsExpired() bool {
	return false
}

// SESService is the service which manages a service queue on the AWS.
type SESService struct {
	running       bool
	awsSES        *ses.SES
	Configuration SESServiceConfiguration
}

// LoadConfiguration returns a not implemented error.
//
// This happens because this method should be implemented on the struct that
// will use this implementation.
func (service *SESService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ApplyConfiguration applies a given configuration to the service.
func (service *SESService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case SESServiceConfiguration:
		service.Configuration = c
		return nil
	case *SESServiceConfiguration:
		service.Configuration = *c
		return nil
	}
	return http.ErrWrongConfigurationInformed
}

// Restart stops and then starts the service again.
func (service *SESService) Restart() error {
	if service.running {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start initializes the aws client with the configuration.
func (service *SESService) Start() error {
	if !service.running {
		conf := aws.Config{
			Credentials: credentials.NewCredentials(NewCredentialsFromStruct(&service.Configuration)),
			Region:      aws.String(service.Configuration.Region),
		}

		sess, err := session.NewSessionWithOptions(session.Options{
			Config: conf,
		})
		if err != nil {
			return err
		}
		service.awsSES = ses.New(sess)
		service.running = true
	}
	return nil
}

// Stop erases the aws client reference.
func (service *SESService) Stop() error {
	if service.running {
		service.awsSES = nil
		service.running = false
	}
	return nil
}

// RunWithSQS runs a handler passing the reference of a `sqs.SQS` client.
func (service *SESService) RunWithSES(handler func(client *ses.SES) error) error {
	if service.running {
		return handler(service.awsSES)
	}
	return http.ErrServiceNotRunning
}
