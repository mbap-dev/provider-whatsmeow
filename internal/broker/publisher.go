package broker

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"your.org/provider-whatsmeow/internal/config"
	ilog "your.org/provider-whatsmeow/internal/log"
)

// WebhookPublisher publishes Cloud webhook payloads to UnoAPI via RabbitMQ.
// It maintains a single shared connection/channel and declares the
// destination exchange on first use.
type WebhookPublisher struct {
	cfg      *config.Config
	exchange string

	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

var (
	webhookPub     *WebhookPublisher
	webhookPubOnce sync.Once
)

// InitWebhookPublisher initialises a global publisher instance. It is safe
// to call multiple times; the first invocation wins.
func InitWebhookPublisher(cfg *config.Config) error {
	var initErr error
	webhookPubOnce.Do(func() {
		webhookPub = &WebhookPublisher{cfg: cfg, exchange: cfg.AMQPWebhookExchange}
		initErr = webhookPub.ensure()
	})
	return initErr
}

func (p *WebhookPublisher) ensure() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cfg == nil || p.cfg.AMQPURL == "" {
		return fmt.Errorf("AMQP URL not configured")
	}
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil {
		return nil
	}
	// (Re)connect
	if p.conn != nil {
		_ = p.conn.Close()
	}
	conn, err := amqp.Dial(p.cfg.AMQPURL)
	if err != nil {
		return fmt.Errorf("dial amqp: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}
	// Ensure exchange exists
	// UnoAPI expects webhooks on the bridge exchange (direct type),
	// routing key is the session phone (digits only).
	if err := ch.ExchangeDeclare(
		p.exchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("declare exchange: %w", err)
	}
	p.conn = conn
	p.ch = ch
	ilog.Infof("AMQP webhook publisher connected: exchange=%s", p.exchange)
	return nil
}

// PublishWebhook publishes the given Cloud webhook payload to UnoAPI.
// routingKey must be the session phone number (digits only).
func PublishWebhook(routingKey string, cloudPayload any) error {
	if webhookPub == nil {
		return fmt.Errorf("webhook publisher not initialised")
	}
	if err := webhookPub.ensure(); err != nil {
		return err
	}
	body, err := json.Marshal(map[string]any{"payload": cloudPayload})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	pub := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	}
	if err := webhookPub.ch.Publish(
		webhookPub.exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		pub,
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	ilog.Debugf("amqp webhook published rk=%s bytes=%d", routingKey, len(body))
	ilog.Infof("amqp webhook published successfully rk=%s bytes=%d", routingKey, len(body))
	return nil
}
