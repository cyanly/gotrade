package sellsidesim

import (
	"strconv"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"

	"github.com/quickfixgo/quickfix/tag"
	log "github.com/cyanly/gotrade/core/logger"
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
	log.Infof("FIX->SIM: FIX44NewOrderSingle \n%v", msg.String())
	return e.Route(msg, sessionID)
}

func stringPtr(s string) *string { return &s }

func (e *Executor) OnFIX44NewOrderSingle(msg fix44nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	execReport := fix44er.Message{
		ClOrdID:      &msg.ClOrdID,
		Account:      msg.Account,
		OrderID:      e.genOrderID(),
		ExecID:       e.genExecID(),
		ExecType:     enum.ExecType_FILL,
		OrdStatus:    enum.OrdStatus_FILLED,
		Side:         msg.Side,
		Instrument:   msg.Instrument,
		OrderQtyData: &msg.OrderQtyData,
		LeavesQty:    0,
		LastQty:      msg.OrderQtyData.OrderQty,
		CumQty:       *msg.OrderQtyData.OrderQty,
		AvgPx:        *msg.Price,
		LastPx:       msg.Price,
		LastMkt:      stringPtr("SIM"),
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}
