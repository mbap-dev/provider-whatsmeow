package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"your.org/provider-whatsmeow/internal/config"
	ilog "your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/provider"
)

// Envelope é o wrapper esperado no RB: { payload: {...}, id: "...", options: {...} }
type Envelope struct {
	Payload provider.OutgoingMessage `json:"payload"`
	ID      string                   `json:"id,omitempty"`
	Options *struct {
		Endpoint string `json:"endpoint,omitempty"`
		Priority int    `json:"priority,omitempty"`
	} `json:"options,omitempty"`
}

type Consumer struct {
	cfg      *config.Config
	provider *provider.ClientManager
	conn     *amqp.Connection
	ch       *amqp.Channel
	queue    amqp.Queue
}

func NewConsumer(cfg *config.Config, provider *provider.ClientManager) (*Consumer, error) {
	return &Consumer{cfg: cfg, provider: provider}, nil
}

func suffixFromRoutingKey(binding, rk string) string {
	prefix := binding
	if i := strings.IndexAny(binding, "*#"); i >= 0 {
		prefix = strings.TrimSuffix(binding[:i], ".")
	}
	if prefix != "" && strings.HasPrefix(rk, prefix+".") {
		return strings.TrimPrefix(rk, prefix+".")
	}
	parts := strings.Split(rk, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return rk
}

func (c *Consumer) Start(ctx context.Context) error {
	if c.cfg.AMQPURL == "" {
		ilog.Infof("AMQP URL is empty; skipping consumer startup")
		<-ctx.Done()
		return nil
	}
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

	if err := ch.ExchangeDeclare(
		c.cfg.AMQPExchange,
		"topic",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := ch.QueueDeclare(
		c.cfg.AMQPQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	c.queue = queue

	if err := ch.QueueBind(
		c.cfg.AMQPQueue,
		c.cfg.AMQPBinding,
		c.cfg.AMQPExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	deliveries, err := ch.Consume(
		c.cfg.AMQPQueue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume from queue: %w", err)
	}

	ilog.Infof("AMQP consumer connected, waiting for messages on %s", c.cfg.AMQPBinding)

	for {
		select {
		case <-ctx.Done():
			if err := ch.Close(); err != nil {
				ilog.Errorf("failed to close AMQP channel: %v", err)
			}
			if err := conn.Close(); err != nil {
				ilog.Errorf("failed to close AMQP connection: %v", err)
			}
			return nil

		case d, ok := <-deliveries:
			if !ok {
				time.Sleep(500 * time.Millisecond)
				return fmt.Errorf("AMQP deliveries channel closed")
			}

			// Rapidamente ignore payloads de status (não são mensagens a enviar)
			var probe struct {
				Payload map[string]any `json:"payload"`
			}
			if err := json.Unmarshal(d.Body, &probe); err == nil {
				if s, ok := probe.Payload["status"]; ok {
					mid, _ := probe.Payload["message_id"].(string)
					rid, _ := probe.Payload["recipient_id"].(string)
					ilog.Debugf("skipping non-send payload status=%v message_id=%q recipient_id=%q", s, mid, rid)
					continue
				}
			}

			// 1) Tenta decodificar como Envelope (novo formato)
			var env Envelope
			if err := json.Unmarshal(d.Body, &env); err != nil || (env.Payload.Type == "" && env.Payload.Text == nil && env.Payload.Image == nil && env.Payload.Document == nil && env.Payload.Audio == nil && env.Payload.Sticker == nil) {
				// 2) Retrocompat: body antigo direto como OutgoingMessage
				var legacy provider.OutgoingMessage
				if err2 := json.Unmarshal(d.Body, &legacy); err2 != nil {
					ilog.Errorf("failed to decode message: %v (legacy err: %v)", err, err2)
					continue
				}
				env.Payload = legacy
			}

			// Propaga o ID do envelope para MessageID se não vier no payload
			if env.ID != "" && env.Payload.MessageID == "" {
				env.Payload.MessageID = env.ID
			}

			// Resolve a sessão: sufixo do routing key (ex.: provider.whatsmeow.<sess>)
			phoneID := suffixFromRoutingKey(c.cfg.AMQPBinding, d.RoutingKey)

			// Log de entrada (sanitizado)
			toPreview := strings.TrimSpace(env.Payload.To)
			typ := strings.TrimSpace(env.Payload.Type)
			hasText := env.Payload.Text != nil && strings.TrimSpace(env.Payload.Text.Body) != ""
			hasImage := env.Payload.Image != nil && strings.TrimSpace(env.Payload.Image.Link) != ""
			hasDoc := env.Payload.Document != nil && strings.TrimSpace(env.Payload.Document.Link) != ""
			hasAudio := env.Payload.Audio != nil && strings.TrimSpace(env.Payload.Audio.Link) != ""
			mediaURL := strings.TrimSpace(env.Payload.MediaURL)
			ilog.Debugf("amqp message decoded phoneID=%s id=%s type=%q to=%q media_url=%q text=%v image=%v doc=%v audio=%v", phoneID, env.Payload.MessageID, typ, toPreview, mediaURL, hasText, hasImage, hasDoc, hasAudio)

			go func(msg provider.OutgoingMessage, phone string) {
				ctxSend := context.WithValue(context.Background(), provider.CtxKeyPhoneNumberID, phone)
				if err := c.provider.Send(ctxSend, msg); err != nil {
					to := strings.TrimSpace(msg.To)
					media := strings.TrimSpace(msg.MediaURL)
					ilog.Errorf("failed to send message (phoneID=%s id=%s type=%q to=%q media_url=%q): %v", phone, msg.MessageID, strings.TrimSpace(msg.Type), to, media, err)
					if to == "" {
						raw := string(d.Body)
						if len(raw) > 1500 {
							raw = raw[:1500] + "...(truncated)"
						}
						ilog.Debugf("payload_raw: %s", raw)
					}
				}
			}(env.Payload, phoneID)
		}
	}
}
