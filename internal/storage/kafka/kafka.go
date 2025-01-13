package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/hackertron/blocksight/pkg/models"
)

// Kafka represents the Kafka producer and consumer functionality
type Kafka struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

// NewKafka creates a new instance of Kafka
func NewKafka(brokers []string) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &Kafka{producer: producer, consumer: consumer}, nil
}

// Close closes the Kafka producer and consumer
func (k *Kafka) Close() {
	if err := k.producer.Close(); err != nil {
		log.Printf("failed to close Kafka producer: %v", err)
	}
	if err := k.consumer.Close(); err != nil {
		log.Printf("failed to close Kafka consumer: %v", err)
	}
}

// PublishMessage publishes a message to a Kafka topic
func (k *Kafka) PublishMessage(topic string, message *models.Message) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message.Value),
	}

	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}

// ConsumeMessages consumes messages from a Kafka topic
func (k *Kafka) ConsumeMessages(topic string, handler func(*models.Message)) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to start Kafka consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			message := &models.Message{
				Value: string(msg.Value),
			}
			handler(message)
		case err := <-partitionConsumer.Errors():
			log.Printf("error consuming Kafka message: %v", err)
		}
	}
}
