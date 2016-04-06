package execution

import (
	logger "github.com/apex/log"
	proto "github.com/cyanly/gotrade/proto/order"
	util "github.com/cyanly/gotrade/services/marketconnectors"

	"github.com/quickfixgo/quickfix"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	"time"
)

// Common route handler for FIX4.4 Execution Report message
//   this function can be sub-classed to extend with special fields if a market connector requires
func OnFIX44ExecutionReport(msg fix44er.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	logger.Infof("FIX->MC EXEC: \n%v", msg)

	// Execution Report common fields
	er, err := NewExecutionFromFIX44Message(msg)
	if err != nil {
		return err
	}

	ProcessExecutionReport(er)

	return
}

// Construct execution report protobuf struct from FIX message for common fields
func NewExecutionFromFIX44Message(msg fix44er.Message) (*proto.Execution, quickfix.MessageRejectError) {
	er := &proto.Execution{}

	// Required fields
	er.ClientOrderId = *msg.ClOrdID
	er.ExecType = proto.Execution_ExecType(util.FIXEnumToProtoEnum(msg.ExecType))
	er.OrderStatus = proto.OrderStatus(util.FIXEnumToProtoEnum(msg.OrdStatus))
	er.BrokerExecId = msg.ExecID
	er.CumQuantity = msg.CumQty
	er.AvgPrice = msg.AvgPx

	// optional common tags
	er.Quantity = *msg.LastQty
	er.Price = *msg.LastPx
	er.Lastmkt = *msg.LastMkt
	if msg.TransactTime != nil {
		er.BrokerExecDatetime = msg.TransactTime.UTC().Format(time.RFC3339Nano)
	}
	er.BrokerOrderId = msg.OrderID
	if msg.ExecRefID != nil {
		er.PrevBrokerExecId = *msg.ExecRefID
	}
	if msg.Text != nil {
		er.Text = *msg.Text
	}

	//TODO: Check for Dups - Tags.PossDupFlag

	return er, nil
}
