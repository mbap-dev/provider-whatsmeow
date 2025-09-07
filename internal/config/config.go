package config

import "os"

// Config holds all configurable settings for the adapter.  Each field
// corresponds to an environment variable.  Defaults are applied where
// reasonable so the service can run locally with minimal setup.
type Config struct {
	// AMQPURL is the connection string used to connect to the RabbitMQ
	// broker.  Example: amqp://guest:guest@localhost:5672/.
	AMQPURL string
	// AMQPExchange is the name of the exchange to bind to for outgoing
	// messages.  By default this is set to "unoapi.outgoing".
	AMQPExchange string
	// AMQPBinding is the routing key pattern used when binding the
	// queue.  It is set to "provider.whatsmeow.*" to receive all
	// messages destined for this provider.
	AMQPBinding string
	// WebhookBase is the base URL to which webhook payloads should be
	// delivered.  If left empty the service will not emit webhooks.
	WebhookBase string
	// SessionStore is the directory on disk where session files are
	// persisted.  A separate subdirectory will be created for each
	// session identifier.
	SessionStore string
	// HTTPAddr is the host:port on which to expose the HTTP API and
	// health checks.  The default is ":8080" which listens on all
	// interfaces.
	HTTPAddr string
}

// NewConfig reads configuration from the environment and returns a
// populated Config instance.  Missing variables fall back to sensible
// defaults as documented on the struct fields.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.AMQPURL = getEnv("AMQP_URL", "amqp://user_test:123456@localhost:5672/EnvolveNEXT")
	cfg.AMQPExchange = getEnv("AMQP_EXCHANGE", "unoapi.outgoing")
	cfg.AMQPBinding = getEnv("AMQP_BINDING", "provider.whatsmeow.*")
	cfg.WebhookBase = getEnv("WEBHOOK_BASE", "https://localhost/webhooks/whatsapp")
	cfg.SessionStore = getEnv("SESSION_STORE", "./state/whatsmeow")
	cfg.HTTPAddr = getEnv("HTTP_ADDR", ":8080")
	return cfg
}

// getEnv returns the value of the environment variable named by key.  If
// the variable is not present or empty then defaultVal is returned.
func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
