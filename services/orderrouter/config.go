package orderrouter

import "github.com/nats-io/nats"

const (
	// DefaultMessageBusURL is the default MessageBus URL for service communications.
	DefaultMessageBusURL        = nats.DefaultURL
	DefaultServiceMessageBusURL = nats.DefaultURL
)

// Config represents the configuration for base service.
type Config struct {
	MessageBusURL        string `json:"or_message_bus"`      //trade and order message bus URL
	ServiceMessageBusURL string `json:"service_message_bus"` //services message bus URL, for listening to MC heartbeats
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		MessageBusURL:        DefaultMessageBusURL,
		ServiceMessageBusURL: DefaultServiceMessageBusURL,
	}
}
