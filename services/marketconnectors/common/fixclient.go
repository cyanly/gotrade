// A common FIX Client required by QuickFIX, should work for most of market connectors.
//
// If a market connector has special tags etc for incoming message, only need its own NewFIXClient() to
//   replace common calls to its own callback routes.
// If a market connector has complete different behaviour, it will need to implements its own FIXClient
package common

import (
	exeFIX "github.com/cyanly/gotrade/services/marketconnectors/common/execution"

	"github.com/quickfixgo/quickfix"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
)

type FIXClient struct {
	*quickfix.MessageRouter
}

// Create a FIXClient with common routes for market connectors
func NewFIXClient() *FIXClient {
	e := &FIXClient{
		MessageRouter: quickfix.NewMessageRouter(),
	}
	e.AddRoute(fix44er.Route(exeFIX.OnFIX44ExecutionReport))

	return e
}

func (e FIXClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

func (e FIXClient) OnLogon(sessionID quickfix.SessionID) {
	return
}

func (e FIXClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

func (e FIXClient) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

func (e FIXClient) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID) {
	return
}

func (e FIXClient) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) (err error) {
	return
}

func (e FIXClient) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}
