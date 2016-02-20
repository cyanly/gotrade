// Core order APIs
package order

import (
	proto "github.com/cyanly/gotrade/proto/order"

	"errors"
	"fmt"
	"github.com/nats-io/nats"
)

var (
	MessageBus *nats.Conn
)

func IsCompleted(m *proto.Order) bool {
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

func CanCancel(m *proto.Order) bool {
	switch m.OrderStatus {
	case proto.OrderStatus_NEW,
		proto.OrderStatus_PARTIALLY_FILLED,
		proto.OrderStatus_REPLACED:
		return true
	case proto.OrderStatus_REJECTED:
		switch m.Instruction {
		case proto.Order_NEW:
			return false
		default:
			return true
		}

	default:
		return false
	}
}

func CanReplace(m *proto.Order) bool {
	return CanCancel(m)
}

func Validate(m *proto.Order) error {
	if m.OrderId <= 0 {
		return errors.New("Invalid OrderId")
	}

	return nil
}

func Stringify(m *proto.Order) string {
	return fmt.Sprintf("(%s) %s %v %s -> %s", m.Trader, m.Side, m.Quantity, m.Symbol, m.MarketConnector)
}
