package kafka

import (
	"context"
	"errors"
	"log"
	"readmodels/internal/events"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	ConsumerGroup sarama.ConsumerGroup
	eventBus      *events.EventBus
}

func NewKafkaConsumer(brokers []string, eventBus *events.EventBus, infoLog, errorLog *log.Logger) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	groupId := "readmodels-group"

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupId, config)
	if err != nil {
		errorLog.Printf("Error creating consumer group client: %v", err)
		return nil, err
	}

	return &KafkaConsumer{
		infoLog:       infoLog,
		errorLog:      errorLog,
		ConsumerGroup: consumerGroup,
		eventBus:      eventBus,
	}, nil
}

func (k *KafkaConsumer) InitConsumption(ctx context.Context) error {
	consumer := Consumer{
		infoLog:  k.infoLog,
		errorLog: k.errorLog,
		ready:    make(chan bool),
		eventBus: k.eventBus,
	}

	k.infoLog.Println("Initiating Kafka Consumer Group...")

	var wg sync.WaitGroup
	wg.Add(1)
	go k.runConsumerGroup(ctx, &wg, &consumer)

	<-consumer.ready // Await till the consumer has been set up
	k.infoLog.Println("Kafka Consumer up and running...")

	<-ctx.Done()
	k.infoLog.Println("Terminating Kafka Consumer: context cancelled")

	wg.Wait()
	if err := k.ConsumerGroup.Close(); err != nil {
		k.errorLog.Printf("Error closing Kafka Consumer Group: %v\n", err)
		return err
	}

	return nil
}

func (k *KafkaConsumer) runConsumerGroup(ctx context.Context, wg *sync.WaitGroup, consumer *Consumer) {
	defer wg.Done()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.ConsumerGroup.Consume(ctx, getTopics(), consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				k.errorLog.Printf("Consumer Group was closed, error: %v", err)
				return
			}
			k.errorLog.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		err := ctx.Err()
		if err != nil {
			return
		}
		consumer.ready = make(chan bool)
	}
}
