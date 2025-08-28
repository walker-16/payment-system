package kafka

import (
	"context"
	"testing"

	"github.com/IBM/sarama"
)

type TestHandler struct{}

func (h *TestHandler) ConsumeMessage(msg *sarama.ConsumerMessage) error {
	if msg.Value == nil {
		return nil
	}
	return nil
}

func TestProducer(t *testing.T) {
	p, err := NewProducer([]string{"localhost:9092"}, "test-producer")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err = p.SendMessage("test-topic", []byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestConsumer(t *testing.T) {
	h := &TestHandler{}
	c, err := NewConsumer([]string{"localhost:9092"}, "test-group", []string{"test-topic"}, h)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := c.Start(ctx); err != nil {
			t.Logf("consumer stopped: %v", err)
		}
	}()

	// Simulate short wait
	<-ctx.Done()
}
