package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	database "readmodels/infrastructure/db"
	"readmodels/infrastructure/kafka"
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

	provider := NewProvider(infoLog, errorLog, env)
	database, err := provider.ProvideDb(ctx)
	if err != nil {
		os.Exit(1)
	}
	kafkaConsumer, err := provider.ProvideKafkaConsumer()
	if err != nil {
		os.Exit(1)
	}

	infoLog.Println("Readmodels service started")

	var wg sync.WaitGroup
	wg.Add(2)
	go applyMigrations(database, ctx, &wg, infoLog, errorLog)
	go initKafkaConsumption(kafkaConsumer, ctx, &wg, infoLog, errorLog)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	infoLog.Println("Shutting down Readmodels Service...")

	wg.Wait()

	infoLog.Println("Readmodels Service stopped")
}

func applyMigrations(database *database.Database, ctx context.Context, wg *sync.WaitGroup, infoLog, errorLog *log.Logger) {
	defer wg.Done()

	infoLog.Println("Applying migrations...")
	err := database.ApplyMigrations(ctx)
	if err != nil {
		errorLog.Panicln("Migrations failed")
	}
	infoLog.Println("Migrations finished")
}

func initKafkaConsumption(kafkaConsumer *kafka.KafkaConsumer, ctx context.Context, wg *sync.WaitGroup, infoLog, errorLog *log.Logger) {
	defer wg.Done()

	infoLog.Println("Initiating Kafka Consumer Group...")
	err := kafkaConsumer.InitConsumption(ctx)
	if err != nil {
		errorLog.Panicln("Kafka Consumption failed")
	}
	infoLog.Println("Kafka Consumer Group stopped")
}
