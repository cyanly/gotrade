package messagebus

import (
	"github.com/nats-io/nats"
)

type Msg nats.Msg

type MessageBus struct {
	Config     *Config
	Connection *nats.Conn
}

// Create new message bus wrapper, for this repo NATS is hardcoded.
// This module should serve as message routing across multiple protocols (RMQ, 0MQ, Tibrv etc)
func NewMessageBus(config *Config) (*MessageBus, error) {
	if conn, err := nats.Connect(config.Url); err != nil {
		return nil, err
	} else {
		m := &MessageBus{
			Config:     config,
			Connection: conn,
		}
		return m, nil
	}
}

// Close underlying messaging bus driver
func (m *MessageBus) Close() {
	if m.Connection != nil {
		m.Connection.Close()
	}
}
