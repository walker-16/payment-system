package kafka

import (
	"context"
	"fmt"
	"payment-system/pkg/logger"

	"github.com/IBM/sarama"
)

const consumerPrefix = "consumer-"

// ConsumerHandler defines interface to process messages.
type ConsumerHandler interface {
	ConsumeMessage(msg *sarama.ConsumerMessage) error
}

// Consumer wraps a Sarama ConsumerGroup.
type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler ConsumerHandler
	logger  logger.Logger
}

// NewConsumer creates a Kafka consumer.
func NewConsumer(brokers []string, groupID string, topics []string,
	handler ConsumerHandler, log logger.Logger) (*Consumer, error) {
	config := NewSaramaConfig(consumerPrefix + groupID)
	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Error("failed to create consumer group",
			logger.String("groupID", groupID),
			logger.Error(err))
		return nil, err
	}
	return &Consumer{
		group:   group,
		topics:  topics,
		handler: handler,
		logger:  log,
	}, nil
}

// Start consumes messages (blocking).
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("starting consumer",
		logger.String("topics", fmt.Sprintf("%v", c.topics)))

	for {
		if err := c.group.Consume(ctx, c.topics, c); err != nil {
			c.logger.Error("consumer error", logger.Error(err))
			return err
		}
		if ctx.Err() != nil {
			c.logger.Warn("context cancelled, stopping consumer")
			return ctx.Err()
		}
	}
}

// Close closes the consumer group.
func (c *Consumer) Close() error {
	return c.group.Close()
}

// Implement sarama.ConsumerGroupHandler.
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.handler.ConsumeMessage(msg); err != nil {
			c.logger.Error("failed to process message",
				logger.String("topic", msg.Topic),
				logger.Int("partition", int(msg.Partition)),
				logger.Int("offset", int(msg.Offset)),
				logger.Error(err))
		} else {
			c.logger.Debug("message processed",
				logger.String("topic", msg.Topic),
				logger.Int("partition", int(msg.Partition)),
				logger.Int("offset", int(msg.Offset)),
			)
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}
