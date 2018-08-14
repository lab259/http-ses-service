[![CircleCI](https://circleci.com/gh/lab259/http-ses-service.svg?style=shield)](https://circleci.com/gh/lab259/http-ses-service)
[![codecov](https://codecov.io/gh/lab259/http-ses-service/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/http-ses-service)
[![GoDoc](https://godoc.org/github.com/lab259/http-ses-service?status.svg)](http://godoc.org/github.com/lab259/http-ses-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/http-ses-service)](https://goreportcard.com/report/github.com/lab259/http-ses-service)

# http-ses-service

The http-ses-service is the [lab259/http](//github.com/lab259/http) service for
the Amazon Simple Email Service (SES).

It wraps the logic of dealing with credentials and keeping the session.

## Dependencies

It depends on the [lab259/http](//github.com/lab259/http) (and its dependencies,
of course) itself and the [aws/aws-sdk-go](//github.com/aws/aws-sdk-go) library.

## Installation

First, fetch the library to the repository.

	go get github.com/lab259/http-ses-service

## Usage

The service is designed to be "extended" and not used directly.

**srv.go**

```Go
package mail

import (
	"github.com/lab259/http"
	"github.com/lab259/http-ses-service"
)

type MailSQSService struct {
	sessrv.SQSService
}

func (service *MailSQSService) LoadConfiguration() (interface{}, error) {
	var configuration sessrv.SQSServiceConfiguration

	configurationLoader := http.NewFileConfigurationLoader("/etc/mail")
	configurationUnmarshaler := &http.ConfigurationUnmarshalerYaml{}

	config, err := configurationLoader.Load("ses.yml")
	if err != nil {
		return err
	}
	return configurationUnmarshaler.Unmarshal(config, dst)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

```

**example.go**
```Go
// ...

var mq rscsrsv.MailSQSService

func init() {
	configuration, err := mq.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	err = mq.ApplyConfiguration(configuration)
	if err != nil {
		panic(err)
	}
	err = mq.Start()
	if err != nil {
		panic(err)
	}
}

// sendEmail sends an email using the SES
func sendEmail() {
	err := service.RunWithSES(func (s *ses.SES) error {

		// use the `s` variable to implement sending an email.
		// More info on: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-ses-with-go-sdk.html

		return nil
	})
	if err != nil {
		panic(err)
	}
	// ...
}

// ...
```