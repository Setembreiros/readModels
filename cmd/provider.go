package main

import (
	"context"
	"log"
	awsClients "readmodels/infrastructure/aws"
	database "readmodels/infrastructure/db"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Provider struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	env      string
}

func NewProvider(infoLog, errorLog *log.Logger, env string) *Provider {
	return &Provider{
		infoLog:  infoLog,
		errorLog: errorLog,
		env:      env,
	}
}

func (p Provider) ProvideDb() (*database.Database, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig()
	} else {
		cfg, err = provideAwsConfig()
	}
	if err != nil {
		p.errorLog.Fatalf("failed to load aws configuration %s", err)
		return nil, err
	}

	return &database.Database{
		Client: awsClients.NewDynamodbClient(cfg, p.infoLog, p.errorLog),
	}, nil
}

func provideAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}

func provideDevEnvironmentDbConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
}
