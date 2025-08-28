package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// ConsumerHandler defines interface to process messages
type ConsumerHandler interface {
	ConsumeMessage(msg *sarama.ConsumerMessage) error
}

// Consumer wraps a Sarama ConsumerGroup
type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler ConsumerHandler
}

// NewConsumer creates a Kafka consumer
func NewConsumer(brokers []string, groupID string, topics []string, handler ConsumerHandler) (*Consumer, error) {
	config := NewSaramaConfig("consumer-" + groupID)
	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{group: group, topics: topics, handler: handler}, nil
}

// Start consumes messages (blocking)
func (c *Consumer) Start(ctx context.Context) error {
	for {
		if err := c.group.Consume(ctx, c.topics, c); err != nil {
			log.Printf("error consuming messages: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes the consumer group
func (c *Consumer) Close() error {
	return c.group.Close()
}

// Implement sarama.ConsumerGroupHandler
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.handler.ConsumeMessage(msg); err != nil {
			log.Printf("message processing error: %v", err)
		} else {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}
