package sellsidesim

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"

	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"

	logger "github.com/apex/log"
	"strconv"
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

func (e *Executor) OnFIX44NewOrderSingle(msg fix44nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	logger.Infof("FIX->SIM: FIX44NewOrderSingle \n%v", msg.String())
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.String() != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix44er.New(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(float64(orderQty.FIXFloat)),
		field.NewAvgPx(float64(price.FIXFloat)),
	)

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.NewLastQty(float64(orderQty.FIXFloat)))
	execReport.Body.Set(field.NewLastPx(float64(price.FIXFloat)))
	execReport.Body.Set(field.NewLastMkt("SIM"))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.Message, sessionID)

	return
}
