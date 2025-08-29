package kafka

import (
	"payment-system/pkg/logger"

	"github.com/IBM/sarama"
)

type SyncProducerInterface interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
	Close() error
}

// Producer wraps a Sarama SyncProducer.
type Producer struct {
	syncProducer SyncProducerInterface
	logger       logger.Logger
}

// NewProducer creates a Kafka producer.
func NewProducer(brokers []string, clientID string,
	logger logger.Logger) (*Producer, error) {
	config := NewSaramaConfig(clientID)
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Producer{syncProducer: producer}, nil
}

// SendMessage sends a message to Kafka topic with retries.
func (p *Producer) SendMessage(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		p.logger.Error("failed to send message to topic %s: %v", topic, err)
	}

	return err
}

// Close closes the producer connection.
func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
