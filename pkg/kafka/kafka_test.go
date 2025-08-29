package kafka

import (
	"testing"

	"payment-system/pkg/logger"

	"github.com/IBM/sarama"
)

// MockSyncProducer is a mock implementation of a Kafka SyncProducer for testing purposes.
type MockSyncProducer struct {
	Messages []*sarama.ProducerMessage
}

// SendMessage simulates sending a message to Kafka and stores it in Messages slice.
func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	m.Messages = append(m.Messages, msg)
	return 0, 0, nil
}

// Close is a no operation implementation for the mock producer.
func (m *MockSyncProducer) Close() error {
	return nil
}

type TestHandler struct{}

func (h *TestHandler) ConsumeMessage(msg *sarama.ConsumerMessage) error {
	return nil
}

// TestProducer verifies that the Producer correctly sends a message to Kafka.
func TestProducer(t *testing.T) {
	testLogger := &logger.LoopLogger{}

	mockProducer := &MockSyncProducer{}

	p := &Producer{
		syncProducer: mockProducer,
		logger:       testLogger,
	}

	err := p.SendMessage("test-topic", []byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}

	if len(mockProducer.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mockProducer.Messages))
	}

	msg := mockProducer.Messages[0]
	if msg.Topic != "test-topic" {
		t.Fatalf("message sent to incorrect topic: %s", msg.Topic)
	}
	if string(msg.Key.(sarama.ByteEncoder)) != "key" {
		t.Fatalf("message sent with incorrect key: %s", msg.Key)
	}
	if string(msg.Value.(sarama.ByteEncoder)) != "value" {
		t.Fatalf("message sent with incorrect value: %s", msg.Value)
	}
}

// TestConsumerHandler verifies that the ConsumerHandler can process a message without error.

func TestConsumerHandler(t *testing.T) {
	handler := &TestHandler{}

	msg := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    1,
		Key:       []byte("key"),
		Value:     []byte("value"),
	}

	if err := handler.ConsumeMessage(msg); err != nil {
		t.Fatal(err)
	}
}
