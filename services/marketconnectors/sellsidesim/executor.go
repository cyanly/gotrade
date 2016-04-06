package sellsidesim

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"

	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"

	logger "github.com/apex/log"
	"strconv"
	"github.com/quickfixgo/quickfix/fix44/orderqtydata"
)

type Executor struct {
	orderID int
	execID  int
	*quickfix.MessageRouter
}

func NewExecutor() *Executor {
	e := &Executor{MessageRouter: quickfix.NewMessageRouter()}
	//e.AddRoute(fix40nos.Route(e.OnFIX40NewOrderSingle))
	//e.AddRoute(fix41nos.Route(e.OnFIX41NewOrderSingle))
	//e.AddRoute(fix42nos.Route(e.OnFIX42NewOrderSingle))
	//e.AddRoute(fix43nos.Route(e.OnFIX43NewOrderSingle))
	e.AddRoute(fix44nos.Route(e.OnFIX44NewOrderSingle))
	//e.AddRoute(fix50nos.Route(e.OnFIX50NewOrderSingle))

	return e
}

func (e *Executor) genOrderID() string {
	e.orderID++
	return strconv.Itoa(e.orderID)
}

func (e *Executor) genExecID() string {
	e.execID++
	return strconv.Itoa(e.execID)
}

//quickfix.Application interface
func (e Executor) OnCreate(sessionID quickfix.SessionID)                          { return }
func (e Executor) OnLogon(sessionID quickfix.SessionID)                           { return }
func (e Executor) OnLogout(sessionID quickfix.SessionID)                          { return }
func (e Executor) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID)     { return }
func (e Executor) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e Executor) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *Executor) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func stringPtr(s string) *string { return &s }

func (e *Executor) OnFIX44NewOrderSingle(msg fix44nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	logger.Infof("FIX->SIM: FIX44NewOrderSingle \n%v", msg)

	symbol := msg.Symbol
	side := msg.Side
	orderQty := msg.OrderQty
	ordType := msg.OrdType

	if ordType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(quickfix.Tag(40))
		return
	}
	price := msg.Price
	clOrdID := msg.ClOrdID

	execReport := fix44er.Message{
		ClOrdID: stringPtr(e.genOrderID()),
		ExecID: e.genExecID(),
		ExecType: enum.ExecType_FILL,
		OrdStatus: enum.OrdStatus_FILLED,
		Side: side,
		LeavesQty: 0,
		CumQty: *orderQty,
		AvgPx: *price,
	}

	execReport.ClOrdID = &clOrdID
	execReport.Instrument.Symbol = symbol
	execReport.OrderQtyData = orderqtydata.New()
	execReport.OrderQtyData.OrderQty = orderQty
	execReport.LastQty = orderQty
	execReport.LastPx = price
	execReport.LastMkt = stringPtr("SIM")
	execReport.Account = msg.Account

	quickfix.SendToTarget(execReport, sessionID)

	return
}
