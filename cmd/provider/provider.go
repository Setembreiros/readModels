package provider

import (
	"context"
	"log"
	awsClients "readmodels/infrastructure/aws"
	"readmodels/infrastructure/kafka"
	"readmodels/internal/api"
	"readmodels/internal/bus"
	database "readmodels/internal/db"
	"readmodels/internal/userprofile"
	"readmodels/internal/userprofile/handlers"

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

func (p *Provider) ProvideApiEndpoint(database *database.Database) *api.Api {
	return api.NewApiEndpoint(p.infoLog, p.errorLog, p.ProvideApiControllers(database))
}

func (p *Provider) ProvideApiControllers(database *database.Database) []api.Controller {
	return []api.Controller{
		userprofile.NewUserProfileController(p.infoLog, p.errorLog, userprofile.UserProfileRepository(*database)),
	}
}

func (p *Provider) ProvideEventBus() *bus.EventBus {
	eventBus := bus.NewEventBus()

	return eventBus
}

func (p *Provider) ProvideSubscriptions(database *database.Database) *[]bus.EventSubscription {
	return &[]bus.EventSubscription{
		{
			EventType: "UserWasRegisteredEvent",
			Handler:   handlers.NewUserWasRegisteredEventHandler(p.infoLog, p.errorLog, userprofile.UserProfileRepository(*database)),
		},
	}
}

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
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

func (p *Provider) ProvideDb(ctx context.Context) (*database.Database, error) {
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

	return database.NewDatabase(awsClients.NewDynamodbClient(cfg, p.infoLog, p.errorLog), p.infoLog), nil
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
