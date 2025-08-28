package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Producer wraps a Sarama SyncProducer
type Producer struct {
	syncProducer sarama.SyncProducer
}

// NewProducer creates a Kafka producer
func NewProducer(brokers []string, clientID string) (*Producer, error) {
	config := NewSaramaConfig(clientID)
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Producer{syncProducer: producer}, nil
}

// SendMessage sends a message to Kafka topic with retries
func (p *Producer) SendMessage(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return err
}

// Close closes the producer connection
func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
