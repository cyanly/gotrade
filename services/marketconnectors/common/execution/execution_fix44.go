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
	logger.Infof("FIX->MC EXEC: \n%v", msg.Message.String())

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

	if clOrdId, err := msg.ClOrdID(); err != nil {
		return nil, err
	} else {
		er.ClientOrderId = string(clOrdId.FIXString)

		//split into order key and version
	}

	if execType, err := msg.ExecType(); err != nil {
		return nil, err
	} else {
		er.ExecType = proto.Execution_ExecType(util.FIXEnumToProtoEnum(string(execType.FIXString)))
	}

	if ordStatus, err := msg.OrdStatus(); err != nil {
		return nil, err
	} else {
		er.OrderStatus = proto.OrderStatus(util.FIXEnumToProtoEnum(string(ordStatus.FIXString)))
	}

	if execId, err := msg.ExecID(); err != nil {
		return nil, err
	} else {
		er.BrokerExecId = string(execId.FIXString)
	}

	if cumQty, err := msg.CumQty(); err != nil {
		return nil, err
	} else {
		er.CumQuantity = float64(cumQty.FIXFloat)
	}

	if avgPx, err := msg.AvgPx(); err != nil {
		return nil, err
	} else {
		er.AvgPrice = float64(avgPx.FIXFloat)
	}

	// optional common tags

	if lastQty, err := msg.LastQty(); err == nil {
		er.Quantity = float64(lastQty.FIXFloat)
	}

	if lastPx, err := msg.LastPx(); err == nil {
		er.Price = float64(lastPx.FIXFloat)
	}

	if lastMkt, err := msg.LastMkt(); err == nil {
		er.Lastmkt = string(lastMkt.FIXString)
	}

	if execTime, err := msg.TransactTime(); err == nil {
		er.BrokerExecDatetime = execTime.FIXUTCTimestamp.Value.UTC().Format(time.RFC3339Nano)
	}

	if brkOrdId, err := msg.OrderID(); err == nil {
		er.BrokerOrderId = string(brkOrdId.FIXString)
	}

	if prevExecId, err := msg.ExecRefID(); err == nil {
		er.PrevBrokerExecId = string(prevExecId.FIXString)
	}

	if text, err := msg.Text(); err == nil {
		er.Text = string(text.FIXString)
	}

	//TODO: Check for Dups - Tags.PossDupFlag

	return er, nil
}
