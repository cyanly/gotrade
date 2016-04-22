package orderrouter

import "github.com/cyanly/gotrade/core/messagebus"

const (
	// DefaultMessageBusURL is the default MessageBus URL for service communications.
	DefaultMessageBusURL        = messagebus.DefaultUrl
	DefaultServiceMessageBusURL = messagebus.DefaultUrl
)

// Config represents the configuration for base service.
type Config struct {
	MessageBusURL        string `json:"or_message_bus"`      //trade and order message bus URL
	ServiceMessageBusURL string `json:"service_message_bus"` //services message bus URL, for listening to MC heartbeats
	DatabaseDriver       string `json:"database_driver"`     //database storage engine driver
	DatabaseUrl          string `json:"database_url"`        //database connection string
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		MessageBusURL:        DefaultMessageBusURL,
		ServiceMessageBusURL: DefaultServiceMessageBusURL,
	}
}
