package execution

import (
	proto "github.com/cyanly/gotrade/proto/order"
	util "github.com/cyanly/gotrade/services/marketconnectors"

	"github.com/quickfixgo/quickfix"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	"log"
	"time"
)

// Common route handler for FIX4.4 Execution Report message
//   this function can be sub-classed to extend with special fields if a market connector requires
func OnFIX44ExecutionReport(msg fix44er.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	log.Println("FIX->MC EXEC: \n", msg.Message.String())

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
		pclordId := string(clOrdId.FIXString)
		er.ClientOrderId = &pclordId

		//split into order key and version
	}

	if execType, err := msg.ExecType(); err != nil {
		return nil, err
	} else {
		pexecType := proto.Execution_ExecType(util.FIXEnumToProtoEnum(string(execType.FIXString)))
		er.ExecType = &pexecType
	}

	if ordStatus, err := msg.OrdStatus(); err != nil {
		return nil, err
	} else {
		pordStatus := proto.OrderStatus(util.FIXEnumToProtoEnum(string(ordStatus.FIXString)))
		er.OrderStatus = &pordStatus
	}

	if execId, err := msg.ExecID(); err != nil {
		return nil, err
	} else {
		pexecId := string(execId.FIXString)
		er.BrokerExecId = &pexecId
	}

	if cumQty, err := msg.CumQty(); err != nil {
		return nil, err
	} else {
		pcumQty := float64(cumQty.FIXFloat)
		er.CumQuantity = &pcumQty
	}

	if avgPx, err := msg.AvgPx(); err != nil {
		return nil, err
	} else {
		pavgPx := float64(avgPx.FIXFloat)
		er.AvgPrice = &pavgPx
	}

	// optional common tags

	if lastQty, err := msg.LastQty(); err == nil {
		plastQty := float64(lastQty.FIXFloat)
		er.Quantity = &plastQty
	}

	if lastPx, err := msg.LastPx(); err == nil {
		plastPx := float64(lastPx.FIXFloat)
		er.Price = &plastPx
	}

	if lastMkt, err := msg.LastMkt(); err == nil {
		plastMkt := string(lastMkt.FIXString)
		er.Lastmkt = &plastMkt
	}

	if execTime, err := msg.TransactTime(); err == nil {
		pexecTime := execTime.FIXUTCTimestamp.Value.UTC().Format(time.RFC3339Nano)
		er.BrokerExecDatetime = &pexecTime
	}

	if brkOrdId, err := msg.OrderID(); err == nil {
		pbrkOrdId := string(brkOrdId.FIXString)
		er.BrokerOrderId = &pbrkOrdId
	}

	if prevExecId, err := msg.ExecRefID(); err == nil {
		pprevExecId := string(prevExecId.FIXString)
		er.PrevBrokerExecId = &pprevExecId
	}

	if text, err := msg.Text(); err == nil {
		ptext := string(text.FIXString)
		er.Text = &ptext
	}

	//TODO: Check for Dups - Tags.PossDupFlag

	return er, nil
}
