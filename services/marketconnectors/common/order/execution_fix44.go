package order

import (
	"time"

	"github.com/quickfixgo/quickfix"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"

	proto "github.com/cyanly/gotrade/proto/order"
	util "github.com/cyanly/gotrade/services/marketconnectors"
)

// Common route handler for FIX4.4 Execution Report message
//   this function can be sub-classed to extend with special fields if a market connector requires
func (app FIXClient) onFIX44ExecutionReport(msg fix44er.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	// Construct execution report protobuf struct from FIX message
	er := &proto.Execution{
		// required fields
		ClientOrderId: *msg.ClOrdID,
		BrokerOrderId: msg.OrderID,
		ExecType:      proto.Execution_ExecType(util.FIXEnumToProtoEnum(msg.ExecType)),
		OrderStatus:   proto.OrderStatus(util.FIXEnumToProtoEnum(msg.OrdStatus)),
		BrokerExecId:  msg.ExecID,
		CumQuantity:   msg.CumQty,
		AvgPrice:      msg.AvgPx,

		// common tags
		Quantity: *msg.LastQty,
		Price:    *msg.LastPx,
		Lastmkt:  *msg.LastMkt,
	}

	if msg.TransactTime != nil {
		er.BrokerExecDatetime = msg.TransactTime.UTC().Format(time.RFC3339Nano)
	}
	if msg.ExecRefID != nil {
		er.PrevBrokerExecId = *msg.ExecRefID
	}
	if msg.Text != nil {
		er.Text = *msg.Text
	}

	//TODO: Check for Dups - Tags.PossDupFlag

	return app.processExecutionReport(er)
}
