// Core order APIs
package order

import (
	"errors"
	"fmt"

	proto "github.com/cyanly/gotrade/proto/order"
	messagebus "github.com/nats-io/nats.go"
)

var (
	MessageBus *messagebus.Conn
)

// Order struct extends protobuf Order /proto/order/order.proto
type Order struct {
	*proto.Order
}

func (m *Order) IsCompleted() bool {
	switch m.OrderStatus {
	case proto.OrderStatus_CANCELLED,
		proto.OrderStatus_REJECTED,
		proto.OrderStatus_FILLED,
		proto.OrderStatus_DONE_FOR_DAY,
		proto.OrderStatus_EXPIRED:
		return true

	default:
		return false
	}
}

// Status that Order can be cancelled
func (m *Order) CanCancel() bool {
	switch m.OrderStatus {
	case proto.OrderStatus_NEW,
		proto.OrderStatus_PARTIALLY_FILLED,
		proto.OrderStatus_REPLACED:
		return true
	case proto.OrderStatus_REJECTED:
		switch m.MessageType {
		case proto.Order_NEW:
			return false
		default:
			return true
		}

	default:
		return false
	}
}

// Status that Order can be replaced
func (m *Order) CanReplace() bool {
	return m.CanCancel()
}

// Basic validation of Order, usually unmarshaled off message bus
func (m *Order) Validate() error {
	if m.OrderId <= 0 {
		return errors.New("Invalid OrderId")
	}

	return nil
}

// Pretty print order, usually for logging
func (m *Order) String() string {
	return fmt.Sprintf("(%s) %s %v %s -> %s", m.Trader, m.Side, m.Quantity, m.Symbol, m.MarketConnector)
}
