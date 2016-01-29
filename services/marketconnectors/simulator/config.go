package simulator

import "github.com/nats-io/nats"

const (
	// DefaultMessageBusURL is the default MessageBus URL for service communications.
	DefaultMessageBusURL = nats.DefaultURL
)

// Config represents the configuration for base service.
type Config struct {
	MessageBusURL string `json:"mc_message_bus"`
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		MessageBusURL: DefaultMessageBusURL,
	}
}
