package provider

import (
	"context"
	awsClients "readmodels/infrastructure/aws"
	"readmodels/infrastructure/kafka"
	"readmodels/internal/api"
	"readmodels/internal/bus"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	database "readmodels/internal/db"
	"readmodels/internal/follow"
	"readmodels/internal/post"
	post_handler "readmodels/internal/post/handler"
	"readmodels/internal/reaction"
	reaction_handler "readmodels/internal/reaction/handler"
	"readmodels/internal/userprofile"
	userprofile_handler "readmodels/internal/userprofile/handlers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/zerolog/log"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
	}
}

func (p *Provider) ProvideApiEndpoint(database *database.Database) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database))
}

func (p *Provider) ProvideApiControllers(database *database.Database) []api.Controller {
	return []api.Controller{
		userprofile.NewUserProfileController(userprofile.UserProfileRepository(*database)),
		post.NewPostController(post.NewPostService(post.PostRepository(*database))),
		follow.NewFollowController(follow.FollowRepository(*database)),
		comment.NewCommentController(comment.NewCommentRepository(database)),
		reaction.NewReactionController(reaction.NewReactionService(reaction.NewReactionRepository(database))),
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
			Handler:   userprofile_handler.NewUserWasRegisteredEventHandler(userprofile.UserProfileRepository(*database)),
		},
		{
			EventType: "UserProfileUpdatedEvent",
			Handler:   userprofile_handler.NewUserProfileUpdatedEventHandler(userprofile.UserProfileRepository(*database)),
		},
		{
			EventType: "UserAFollowedUserBEvent",
			Handler:   userprofile_handler.NewUserAFollowedUserBEventHandler(userprofile.UserProfileRepository(*database)),
		},
		{
			EventType: "UserAUnfollowedUserBEvent",
			Handler:   userprofile_handler.NewUserAUnfollowedUserBEventHandler(userprofile.UserProfileRepository(*database)),
		},
		{
			EventType: "PostWasCreatedEvent",
			Handler:   post_handler.NewPostWasCreatedEventHandler(post.NewPostService(post.PostRepository(*database))),
		},
		{
			EventType: "PostsWereDeletedEvent",
			Handler:   post_handler.NewPostsWereDeletedEventHandler(post.PostRepository(*database)),
		},
		{
			EventType: "CommentWasCreatedEvent",
			Handler:   comment_handler.NewCommentWasCreatedEventHandler(comment.NewCommentService(comment.NewCommentRepository(database))),
		},
		{
			EventType: "CommentWasUpdatedEvent",
			Handler:   comment_handler.NewCommentWasUpdatedEventHandler(comment.NewCommentRepository(database)),
		},
		{
			EventType: "CommentWasDeletedEvent",
			Handler:   comment_handler.NewCommentWasDeletedEventHandler(comment.NewCommentService(comment.NewCommentRepository(database))),
		},
		{
			EventType: "UserLikedPostEvent",
			Handler:   reaction_handler.NewUserLikedPostEventHandler(reaction.NewReactionService(reaction.NewReactionRepository(database))),
		},
		{
			EventType: "UserUnlikedPostEvent",
			Handler:   reaction_handler.NewUserUnlikedPostEventHandler(reaction.NewReactionService(reaction.NewReactionRepository(database))),
		},
		{
			EventType: "UserSuperlikedPostEvent",
			Handler:   reaction_handler.NewUserSuperlikedPostEventHandler(reaction.NewReactionService(reaction.NewReactionRepository(database))),
		},
		{
			EventType: "UserUnsuperlikedPostEvent",
			Handler:   reaction_handler.NewUserUnsuperlikedPostEventHandler(reaction.NewReactionService(reaction.NewReactionRepository(database))),
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
			"172.31.0.242:9092",
			"172.31.7.110:9092",
		}
	}

	return kafka.NewKafkaConsumer(brokers, eventBus)
}
func (p *Provider) ProvideDb(ctx context.Context) (*database.Database, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" || p.env == "test" {
		cfg, err = provideDevEnvironmentDbConfig(ctx)
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load aws configuration")
		return nil, err
	}

	return database.NewDatabase(awsClients.NewDynamodbClient(cfg)), nil
}

func provideAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-3"))
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
