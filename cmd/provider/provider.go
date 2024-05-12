package provider

import (
	"context"
	"log"
	awsClients "readmodels/infrastructure/aws"
	database "readmodels/infrastructure/db"
	"readmodels/infrastructure/kafka"
	"readmodels/internal/events"
	"readmodels/internal/events/handlers"
	"readmodels/internal/userprofile"

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

func (p Provider) ProvideEventBus() *events.EventBus {
	eventBus := events.NewEventBus()

	return eventBus
}

func (p Provider) ProvideSubscriptions(database *database.Database) []events.EventSubscription {
	return []events.EventSubscription{
		{
			EventType: "UserWasRegisteredEvent",
			Handler:   handlers.NewUserWasRegisteredEventHandler(p.infoLog, p.errorLog, userprofile.UserProfileRepository(*database)),
		},
	}
}

func (p Provider) ProvideKafkaConsumer(eventBus *events.EventBus) (*kafka.KafkaConsumer, error) {
	var brokers []string

	if p.env == "development" {
		brokers = []string{
			"localhost:9093",
		}
	} else {
		brokers = []string{
			"localhost:9093",
		}
	}

	return kafka.NewKafkaConsumer(brokers, eventBus, p.infoLog, p.errorLog)
}

func (p Provider) ProvideDb(ctx context.Context) (*database.Database, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig(ctx)
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		p.errorLog.Fatalf("failed to load aws configuration %s", err)
		return nil, err
	}

	return &database.Database{
		Client: awsClients.NewDynamodbClient(cfg, p.infoLog, p.errorLog),
	}, nil
}

func provideAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func provideDevEnvironmentDbConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
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
