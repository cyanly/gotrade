package sellsidesim

import (
	"github.com/quickfixgo/quickfix"

	"log"
)

type SellSideSimulator struct {
	acceptor *quickfix.Acceptor
}

func NewSellSideSimulator(logger string) *SellSideSimulator {
	mc := &SellSideSimulator{}

	// QuickFIX settings
	appSettings := quickfix.NewSettings()
	var settings *quickfix.SessionSettings
	settings = appSettings.GlobalSettings()
	settings.Set("SocketAcceptPort", "5001")
	settings.Set("SenderCompID", "CORP")
	settings.Set("TargetCompID", "CYAN")
	settings.Set("ResetOnLogon", "Y")
	settings.Set("FileLogPath", "tmp")
	settings = quickfix.NewSessionSettings()
	settings.Set("BeginString", "FIX.4.4")
	appSettings.AddSession(settings)

	app := NewExecutor()

	var logFactory quickfix.LogFactory
	switch logger {
	case "console":
		logFactory = quickfix.NewScreenLogFactory()
	default:
		logFactory = quickfix.NewNullLogFactory()

	}

	acceptor, err := quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
	if err != nil {
		log.Panic(err)
	}

	mc.acceptor = acceptor

	return mc
}

func (m *SellSideSimulator) Start() {
	if err := m.acceptor.Start(); err != nil {
		log.Panic(err)
	}
}

func (m *SellSideSimulator) Close() {
	m.acceptor.Stop()
}
