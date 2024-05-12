package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"readmodels/cmd/provider"
	database "readmodels/infrastructure/db"
	"readmodels/infrastructure/kafka"
	"readmodels/internal/events"
	"strings"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Starting ReadModels service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(infoLog, errorLog, env)
	database, err := provider.ProvideDb(ctx)
	if err != nil {
		os.Exit(1)
	}
	eventBus := provider.ProvideEventBus()
	subscriptions := provider.ProvideSubscriptions(database)
	kafkaConsumer, err := provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}

	var configuringTasks sync.WaitGroup
	var runningTasks sync.WaitGroup
	configuringTasks.Add(2)
	go applyMigrations(database, ctx, &configuringTasks, infoLog, errorLog)
	go subcribeEvents(subscriptions, eventBus, ctx, &configuringTasks, infoLog) // Always subscribe event before init Kafka
	configuringTasks.Wait()
	configuringTasks.Add(1)
	runningTasks.Add(1)
	go initKafkaConsumption(kafkaConsumer, ctx, &configuringTasks, &runningTasks, infoLog, errorLog)
	configuringTasks.Wait()

	run(infoLog, &configuringTasks, &runningTasks, cancel)
}

func applyMigrations(database *database.Database, ctx context.Context, configurationTasks *sync.WaitGroup, infoLog, errorLog *log.Logger) {
	defer configurationTasks.Done()

	infoLog.Println("Applying migrations...")
	err := database.ApplyMigrations(ctx)
	if err != nil {
		errorLog.Panicln("Migrations failed")
	}
	infoLog.Println("Migrations finished")
}

func subcribeEvents(subscriptions []events.EventSubscription, eventBus *events.EventBus, ctx context.Context, configurationTasks *sync.WaitGroup, infoLog *log.Logger) {
	defer configurationTasks.Done()

	infoLog.Println("Subscribing events...")

	for _, subscription := range subscriptions {
		eventBus.Subscribe(subscription, ctx)
		infoLog.Printf("%s subscribed\n", subscription.EventType)
	}

	infoLog.Println("All events subscribed")
}

func initKafkaConsumption(kafkaConsumer *kafka.KafkaConsumer, ctx context.Context, configurationTasks, runningTasks *sync.WaitGroup, infoLog, errorLog *log.Logger) {
	defer runningTasks.Done()

	infoLog.Println("Initiating Kafka Consumer Group...")
	err := kafkaConsumer.InitConsumption(configurationTasks, ctx)
	if err != nil {
		errorLog.Panicln("Kafka Consumption failed")
	}
	infoLog.Println("Kafka Consumer Group stopped")
}

func run(infoLog *log.Logger, configurationTasks, runningTasks *sync.WaitGroup, cancel context.CancelFunc) {
	configurationTasks.Wait()
	infoLog.Println("Readmodels service started")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()

	infoLog.Println("Shutting down Readmodels Service...")

	runningTasks.Wait()

	infoLog.Println("Readmodels Service stopped")
}
