package common

import "github.com/cyanly/gotrade/core/messagebus"

const (
	// DefaultMessageBusURL is the default MessageBus URL for service communications.
	DefaultMessageBusURL = messagebus.DefaultUrl
)

// Config represents the configuration for base service.
type Config struct {
	MessageBusURL  string `json:"mc_message_bus"`
	DatabaseDriver string `json:"database_driver"` //database storage engine driver
	DatabaseUrl    string `json:"database_url"`    //database connection string

	MarketConnectorName string
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		MessageBusURL: DefaultMessageBusURL,
	}
}
