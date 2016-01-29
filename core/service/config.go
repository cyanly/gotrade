package service

import "github.com/nats-io/nats"

const (
	// DefaultMessageBusURL is the default MessageBus URL for service communications.
	DefaultMessageBusURL = nats.DefaultURL
	// DefaultHeartbeatFreq is the period between service heartbeats in seconds.
	DefaultHeartbeatFreq = int(3)
)

// Config represents the configuration for base service.
type Config struct {
	MessageBusURL string `json:"service_message_bus"`
	ServiceName   string `json:"service_name"`
	HeartbeatFreq int    `json:"service_heartbeat_freq"`
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		ServiceName:   "Unamed Service",
		MessageBusURL: DefaultMessageBusURL,
		HeartbeatFreq: DefaultHeartbeatFreq,
	}
}
