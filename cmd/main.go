package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"readmodels/cmd/provider"
	"readmodels/infrastructure/kafka"
	"readmodels/internal/api"
	"readmodels/internal/bus"
	database "readmodels/internal/db"
	"strings"
	"sync"
	"syscall"
)

type app struct {
	infoLog          *log.Logger
	errorLog         *log.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	configuringTasks sync.WaitGroup
	runningTasks     sync.WaitGroup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &app{
		infoLog:  infoLog,
		errorLog: errorLog,
		ctx:      ctx,
		cancel:   cancel,
	}

	infoLog.Printf("Starting ReadModels service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(infoLog, errorLog, env)
	database, err := provider.ProvideDb(ctx)
	if err != nil {
		os.Exit(1)
	}
	eventBus := provider.ProvideEventBus()
	subscriptions := provider.ProvideSubscriptions(database)
	apiEnpoint := provider.ProvideApiEndpoint(database)
	kafkaConsumer, err := provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}

	app.runConfigurationTasks(database, subscriptions, eventBus)
	app.runServerTasks(kafkaConsumer, apiEnpoint)
}

func (app *app) runConfigurationTasks(database *database.Database, subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	app.configuringTasks.Add(2)
	go app.applyMigrations(database)
	go app.subcribeEvents(subscriptions, eventBus) // Always subscribe event before init Kafka
	app.configuringTasks.Wait()
}

func (app *app) runServerTasks(kafkaConsumer *kafka.KafkaConsumer, apiEnpoint *api.Api) {
	app.runningTasks.Add(2)
	go app.initKafkaConsumption(kafkaConsumer)
	go app.runApiEndpoint(apiEnpoint)

	blockForever()

	app.shutdown()
}

func (app *app) applyMigrations(database *database.Database) {
	defer app.configuringTasks.Done()

	err := database.ApplyMigrations(app.ctx)
	if err != nil {
		app.errorLog.Panicln("Migrations failed")
	}
	app.infoLog.Println("Migrations finished")
}

func (app *app) subcribeEvents(subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	defer app.configuringTasks.Done()

	app.infoLog.Println("Subscribing events...")

	for _, subscription := range *subscriptions {
		eventBus.Subscribe(&subscription, app.ctx)
		app.infoLog.Printf("%s subscribed\n", subscription.EventType)
	}

	app.infoLog.Println("All events subscribed")
}

func (app *app) initKafkaConsumption(kafkaConsumer *kafka.KafkaConsumer) {
	defer app.runningTasks.Done()

	err := kafkaConsumer.InitConsumption(app.ctx)
	if err != nil {
		app.errorLog.Panicln("Kafka Consumption failed")
	}
	app.infoLog.Println("Kafka Consumer Group stopped")
}

func (app *app) runApiEndpoint(apiEnpoint *api.Api) {
	defer app.runningTasks.Done()

	err := apiEnpoint.Run(app.ctx)
	if err != nil {
		app.errorLog.Panicln("Closing Readmodels Api failed")
	}
	app.infoLog.Println("Readmodels Api stopped")
}

func blockForever() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func (app *app) shutdown() {
	app.cancel()
	app.infoLog.Println("Shutting down Readmodels Service...")
	app.runningTasks.Wait()
	app.infoLog.Println("Readmodels Service stopped")
}
