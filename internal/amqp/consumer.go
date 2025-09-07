// consumer.go
package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"your.org/provider-whatsmeow/internal/config"
	"your.org/provider-whatsmeow/internal/provider"
)

// Consumer listens on a RabbitMQ queue for outbound messages destined
// for WhatsApp.  It connects to the broker, declares the exchange
// specified in the configuration and binds a fresh queue to the
// routing key pattern configured via AMQPBinding.  The Consumer
// unmarshals JSON messages into provider.OutgoingMessage and
// dispatches them via the provider.ClientManager.
type Consumer struct {
	cfg      *config.Config
	provider *provider.ClientManager
	conn     *amqp.Connection
	ch       *amqp.Channel
	queue    amqp.Queue
}

// NewConsumer constructs a new Consumer.  It does not connect to
// RabbitMQ until Start is invoked.
func NewConsumer(cfg *config.Config, provider *provider.ClientManager) (*Consumer, error) {
	return &Consumer{cfg: cfg, provider: provider}, nil
}

// suffixFromRoutingKey extracts the wildcard suffix (e.g., phone_number_id)
// from an AMQP routing key given the binding pattern.
// Example:
//
//	binding: "provider.whatsmeow.*"
//	rk:      "provider.whatsmeow.553121159080"
//	return:  "553121159080"
func suffixFromRoutingKey(binding, rk string) string {
	// Take the prefix before the wildcard ("*" or "#"), trimming a trailing dot.
	prefix := binding
	if i := strings.IndexAny(binding, "*#"); i >= 0 {
		prefix = strings.TrimSuffix(binding[:i], ".")
	}
	if prefix != "" && strings.HasPrefix(rk, prefix+".") {
		return strings.TrimPrefix(rk, prefix+".")
	}
	// Fallback: last token.
	parts := strings.Split(rk, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return rk
}

// Start opens the connection to RabbitMQ and begins consuming
// messages.  It blocks until the provided context is cancelled.  Any
// unrecoverable error is returned.  On shutdown Start closes the
// channel and connection gracefully.  Messages with a type of
// "text" are forwarded to the provider while all other types are
// ignored.
func (c *Consumer) Start(ctx context.Context) error {
	if c.cfg.AMQPURL == "" {
		// If no broker URL is configured then there is nothing to do.
		log.Println("AMQP URL is empty; skipping consumer startup")
		<-ctx.Done()
		return nil
	}
	// Establish connection to the AMQP server.
	conn, err := amqp.Dial(c.cfg.AMQPURL)
	if err != nil {
		return fmt.Errorf("failed to dial AMQP: %w", err)
	}
	c.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	c.ch = ch
	// Declare the exchange; it is assumed to be a durable topic exchange.
	if err := ch.ExchangeDeclare(
		c.cfg.AMQPExchange,
		"topic", // type
		true,    // durable
		false,   // autoDelete
		false,   // internal
		false,   // noWait
		nil,     // args
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	// Declare an exclusive queue with a generated name.
	queue, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	c.queue = queue
	// Bind the queue to the exchange using the configured routing key.
	if err := ch.QueueBind(
		queue.Name,
		c.cfg.AMQPBinding,
		c.cfg.AMQPExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}
	// Begin consuming messages with automatic acknowledgements.
	deliveries, err := ch.Consume(
		queue.Name,
		"",    // consumer tag
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume from queue: %w", err)
	}
	// Notify readiness via log.  In a real readiness probe this would
	// integrate with a health system.
	log.Printf("AMQP consumer connected, waiting for messages on %s", c.cfg.AMQPBinding)
	// Listen for deliveries or context cancellation.
	for {
		select {
		case <-ctx.Done():
			// Close channel and connection when context is cancelled.
			if err := ch.Close(); err != nil {
				log.Printf("failed to close AMQP channel: %v", err)
			}
			if err := conn.Close(); err != nil {
				log.Printf("failed to close AMQP connection: %v", err)
			}
			return nil
		case d, ok := <-deliveries:
			if !ok {
				// Channel closed by broker or error; wait briefly and exit.
				time.Sleep(500 * time.Millisecond)
				return fmt.Errorf("AMQP deliveries channel closed")
			}
			var outgoing provider.OutgoingMessage
			if err := json.Unmarshal(d.Body, &outgoing); err != nil {
				log.Printf("failed to decode message: %v", err)
				continue
			}

			// Extract wildcard suffix (e.g., phone_number_id) from routing key
			// and pass it to the provider via context.
			phoneID := suffixFromRoutingKey(c.cfg.AMQPBinding, d.RoutingKey)

			// Dispatch the message to the provider with phoneID in context.
			go func(msg provider.OutgoingMessage, phone string) {
				ctxSend := context.WithValue(context.Background(), provider.CtxKeyPhoneNumberID, phone)
				if err := c.provider.Send(ctxSend, msg); err != nil {
					log.Printf("failed to send message (phoneID=%s): %v", phone, err)
				}
			}(outgoing, phoneID)
		}
	}
}
