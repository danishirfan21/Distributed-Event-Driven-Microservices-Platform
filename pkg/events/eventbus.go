package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type EventBus struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewEventBus(url string) (*EventBus, error) {
	nc, err := nats.Connect(url, nats.RetryOnFailedConnect(true), nats.MaxReconnects(10), nats.ReconnectWait(time.Second*2))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	return &EventBus{nc: nc, js: js}, nil
}

func (eb *EventBus) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = eb.js.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	return nil
}

func (eb *EventBus) Subscribe(subject, queue string, handler func([]byte) error) (*nats.Subscription, error) {
	sub, err := eb.js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			// In production, you might want to handle retries or move to DLQ
			fmt.Printf("Error handling message on %s: %v\n", subject, err)
			msg.Nak()
			return
		}
		msg.Ack()
	}, nats.ManualAck(), nats.DeliverNew())

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	return sub, nil
}

func (eb *EventBus) Close() {
	eb.nc.Close()
}

// Event subjects
const (
	OrderCreated      = "order.created"
	PaymentConfirmed  = "payment.confirmed"
	PaymentFailed     = "payment.failed"
	InventoryUpdated  = "inventory.updated"
)
