// A common FIX Client required by QuickFIX, should work for most of market connectors.
//
// If a market connector has special tags etc for incoming message, only need its own NewFIXClient() to
//   replace common calls to its own callback routes.
// If a market connector has complete different behaviour, it will need to implements its own FIXClient
package order

import (
	"strconv"
	"time"

	messagebus "github.com/nats-io/nats"
	"github.com/quickfixgo/quickfix"
	fix44er "github.com/quickfixgo/fix44/executionreport"
	fix44nos "github.com/quickfixgo/fix44/newordersingle"
	fix44ocj "github.com/quickfixgo/fix44/ordercancelreject"
	"github.com/quickfixgo/tag"

	"strings"
	log "github.com/cyanly/gotrade/core/logger"
	"github.com/cyanly/gotrade/database"
	proto "github.com/cyanly/gotrade/proto/order"
	util "github.com/cyanly/gotrade/services/marketconnectors"
	"github.com/cyanly/gotrade/services/marketconnectors/common"
)

type FIXClient struct {
	*quickfix.Initiator
	*quickfix.MessageRouter

	Session    quickfix.SessionID
	MessageBus *messagebus.Conn
	OrderStore database.OrderStore

	marketConnectorName string
}

// Create a FIXClient with common routes for market connectors
func NewFIXClient(c common.Config) *FIXClient {
	app := &FIXClient{
		MessageRouter: quickfix.NewMessageRouter(),

		marketConnectorName: c.MarketConnectorName,
	}

	// Initiate message bus listening for requests
	if msgbus, err := messagebus.Connect(c.MessageBusURL); err != nil {
		log.Fatalf("error: Cannot connect to order message bus @ %v", c.MessageBusURL)
	} else {
		app.MessageBus = msgbus
	}

	// Connect to database storage driver
	if storage, err := database.NewOrderStore(c.DatabaseDriver, c.DatabaseUrl, nil); err != nil {
		log.Fatalf("error: Cannot connect to database driver %v @ %v", c.DatabaseDriver, c.DatabaseUrl)
	} else {
		app.OrderStore = storage
	}

	// QuickFIX settings from config
	appSettings := quickfix.NewSettings()
	var settings *quickfix.SessionSettings
	settings = appSettings.GlobalSettings()
	settings.Set("SocketConnectHost", "127.0.0.1")
	settings.Set("SocketConnectPort", "5001")
	settings.Set("HeartBtInt", "30")
	settings.Set("SenderCompID", "CYAN")
	settings.Set("TargetCompID", "CORP")
	settings.Set("ResetOnLogon", "Y")
	settings.Set("FileLogPath", "tmp")

	settings = quickfix.NewSessionSettings()
	settings.Set("BeginString", "FIX.4.4")
	if session, err := appSettings.AddSession(settings); err != nil {
		log.WithError(err).Fatal("FIX Session Error")
	} else {
		app.Session = session
	}

	// FIX routes
	app.AddRoute(fix44er.Route(app.onFIX44ExecutionReport))
	app.AddRoute(fix44ocj.Route(app.onFIX44OrderCancelReject))

	// FIX logging
	logFactory := quickfix.NewNullLogFactory() // quickfix.NewScreenLogFactory()

	// Create initiator
	if initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, logFactory); err != nil {
		log.WithError(err).Fatal("FIX NewInitiator Error")
	} else {
		app.Initiator = initiator
	}

	return app
}

func (app FIXClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

func (app FIXClient) OnLogon(sessionID quickfix.SessionID) {
	return
}

func (app FIXClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

func (app FIXClient) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

func (app FIXClient) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID) {
	return
}

func (app FIXClient) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) (err error) {
	return
}

func (app FIXClient) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Infof("FIX->MC EXEC: \n%v", msg.String())
	return app.Route(msg, sessionID)
}

func stringPtr(s string) *string { return &s }

// Common Order FIX client routines serving requests from order router
//   if a market connector has special case it will need to implement its own start routine like below
func (app FIXClient) Start() error {
	err := app.Initiator.Start()

	// Subscribe to order flow topics
	// NEW
	app.MessageBus.Subscribe("order.NewOrderRequest.MC."+app.marketConnectorName, func(m *messagebus.Msg) {
		request := new(proto.NewOrderRequest)
		if err := request.Unmarshal(m.Data); err == nil {
			order := request.Order

			//TODO: this is only prototype, migrate common tasks: instruments / limits processing

			fixmsg := fix44nos.Message{
				ClOrdID:      strconv.Itoa(int(order.OrderKey)) + "." + strconv.Itoa(int(order.Version)),
				Side:         util.ProtoEnumToFIXEnum(int(order.Side)),
				TransactTime: time.Now(),
				OrdType:      util.ProtoEnumToFIXEnum(int(order.OrderType)),
			}

			// Instrument specific
			fixmsg.Instrument.Symbol = &order.Symbol

			fixmsg.OrderQty = &order.Quantity
			if order.OrderType == proto.OrderType_LIMIT || order.OrderType == proto.OrderType_LIMIT_ON_CLOSE {
				fixmsg.Price = &order.LimitPrice
			}

			// Broker specific
			fixmsg.Account = &order.BrokerAccount
			fixmsg.HandlInst = stringPtr(util.ProtoEnumToFIXEnum(int(order.HandleInst)))

			// 142 SenderLocationID
			//     Mandatory for CME exchanges. It contains a 2-character country. For the US and Canada, the state/province is included.
			fixmsg.SenderLocationID = stringPtr("UK")

			log.Info("MC->FIX FIX44NewOrderSingle")
			if err := quickfix.SendToTarget(fixmsg, app.Session); err != nil {
				log.WithError(err).Fatal("FIX quickfix.SendToTarget Error")
			}

		}
	})

	// TODO: CANCEL

	// TODO: REPLACE

	return err
}

// Common behaviours to persist and publish populated Execution entity into our data layer and message bus
func (app FIXClient) processExecutionReport(er *proto.Execution) quickfix.MessageRejectError {

	// Find the Order by using the OrderKey without the Version
	// Need to try to LOAD the order, if that fails then alert stating unsolicited order received
	keyVersion := strings.Split(er.ClientOrderId, ".")
	if len(keyVersion) != 2 {
		return quickfix.ValueIsIncorrect(tag.ClOrdID)
	}
	if oKey, err := strconv.Atoi(keyVersion[0]); err != nil {
		return quickfix.ValueIsIncorrect(tag.ClOrdID)
	} else {
		if ord, err := app.OrderStore.OrderGetByOrderKey(int32(oKey)); err != nil {
			log.WithError(err).Warnf("Received ExecutionReport with OrderKey(%v) does not exist.", oKey)
			return nil // TODO: ignore or quickfix.ValueIsIncorrect(tag.ClOrdID)
		} else {
			er.OrderId = ord.OrderId
		}
	}

	// Persist the execution report
	if err := app.OrderStore.ExecutionCreate(er); err == nil {
		//publish to message bus
		data, _ := er.Marshal()
		app.MessageBus.Publish("order.Execution", data)
	} else {
		log.WithField("execution", er).WithError(err).Error("[ MC ] ERROR Create Reject Execution")
	}

	// Do any validating and if we are good to continue

	// Check if we need to cancel any execution
	// Check if we need to correct any execution
	switch er.ExecType {
	case proto.Execution_TRADE_CANCEL:
	case proto.Execution_TRADE_CORRECT:
	//TODO: TryGetAmendedExecution
	case proto.Execution_RESTATED:
	case proto.Execution_REJECTED:
	}

	// Check if the trade is fully filled or is done for day or canceled
	// .. Send message to TradePump to say trade filled and already to book.

	return nil
}
