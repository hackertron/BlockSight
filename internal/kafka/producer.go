package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/hackertron/blocksight/internal/config"
	"github.com/hackertron/blocksight/internal/models"
)

type Producer struct {
	producer sarama.SyncProducer
	topics   config.KafkaTopics
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topics:   cfg.Kafka.Topics,
	}, nil
}

func (p *Producer) PublishBlock(block *models.Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topics.Blocks,
		Value: sarama.StringEncoder(data),
		Key:   sarama.StringEncoder(block.Hash.Hex()),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) PublishBlocks(blocks []*models.Block) error {
	for _, block := range blocks {
		if err := p.PublishBlock(block); err != nil {
			return err
		}
	}
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
