package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// Config holds Kafka connection config
type Config struct {
	Brokers  []string
	ClientID string
}

// NewSaramaConfig returns a default Sarama config
func NewSaramaConfig(clientID string) *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = clientID
	config.Version = sarama.V2_8_0_0

	// Producer settings
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	// Consumer settings
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
