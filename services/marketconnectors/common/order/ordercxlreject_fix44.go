package order

import (
	"github.com/quickfixgo/quickfix"
	fix44ocj "github.com/quickfixgo/fix44/ordercancelreject"
)

func (app FIXClient) onFIX44OrderCancelReject(msg fix44ocj.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {

	// TODO: not finished for prototype

	return nil
}
