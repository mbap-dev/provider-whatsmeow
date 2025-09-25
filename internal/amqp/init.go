package amqp

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"your.org/provider-whatsmeow/internal/config"
	ilog "your.org/provider-whatsmeow/internal/log"
)

// InitExchange declares the exchange and queue consumed by this service.
// It is safe to call multiple times as declarations are idempotent.
func InitExchange(cfg *config.Config) error {
	if cfg.AMQPURL == "" {
		ilog.Infof("AMQP URL is empty; skipping exchange initialization")
		return nil
	}

	conn, err := amqp.Dial(cfg.AMQPURL)
	if err != nil {
		return fmt.Errorf("failed to dial AMQP: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		cfg.AMQPExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if cfg.AMQPQueue != "" {
		if _, err := ch.QueueDeclare(
			cfg.AMQPQueue,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare queue: %w", err)
		}

		if err := ch.QueueBind(
			cfg.AMQPQueue,
			cfg.AMQPBinding,
			cfg.AMQPExchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	return nil
}
