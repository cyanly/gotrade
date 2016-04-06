// MC Simulator communicates with SellSideSim service in FIX protocol as a mean to test trade life cycle
package simulator

import (
	logger "github.com/apex/log"
	orderCore "github.com/cyanly/gotrade/core/order"
	util "github.com/cyanly/gotrade/services/marketconnectors"
	//execCore "github.com/cyanly/gotrade/core/order/execution"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/services/marketconnectors/common"

	"github.com/quickfixgo/quickfix"
	_ "github.com/quickfixgo/quickfix/enum"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"
	_ "github.com/quickfixgo/quickfix/tag"

	"github.com/nats-io/nats"
	"log"
	"strconv"
	"time"
)

const (
	MarketConnectorName string = "Simulator"
)

type MCSimulator struct {
	Config     Config
	OrdersList map[int32]*proto.Order

	msgbus    *nats.Conn // message bus listening for requests
	initiator *quickfix.Initiator
	session   quickfix.SessionID
}

func NewMarketConnector(c Config) *MCSimulator {
	mc := &MCSimulator{
		Config: c,
	}

	// QuickFIX settings
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
		log.Panic(err)
	} else {
		mc.session = session
	}

	app := common.NewFIXClient()
	//TODO: move these into common fixclient
	logFactory := quickfix.NewNullLogFactory() // quickfix.NewScreenLogFactory()
	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
	if err != nil {
		log.Panic(err)
	}
	mc.initiator = initiator

	// initiate message bus listening for requests
	if msgbus, err := nats.Connect(c.MessageBusURL); err != nil {
		log.Fatal("error: Cannot connect to order message bus @ ", c.MessageBusURL)
	} else {
		mc.msgbus = msgbus
		orderCore.MessageBus = msgbus
	}

	return mc
}

func stringPtr(s string) *string { return &s }

func (mc *MCSimulator) Start() {
	if err := mc.initiator.Start(); err != nil {
		log.Panic(err)
	}

	mc.msgbus.Subscribe("order.NewOrderRequest.MC."+MarketConnectorName, func(m *nats.Msg) {
		request := new(proto.NewOrderRequest)
		if err := request.Unmarshal(m.Data); err == nil {
			order := request.Order

			fixmsg := fix44nos.Message{
				ClOrdID: strconv.Itoa(int(order.OrderKey)) + "." + strconv.Itoa(int(order.Version)),
				Side: util.ProtoEnumToFIXEnum(int(order.Side)),
				TransactTime: time.Now(),
				OrdType: util.ProtoEnumToFIXEnum(int(order.OrderType)),
			}

			// Instrument
			//TODO: migrate marketconnectors/common/... for common fields
			fixmsg.Instrument.Symbol = &order.Symbol

			//TODO: migrate common limit checks into common/limit
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

			logger.Info("MC->FIX FIX44NewOrderSingle")
			if err := quickfix.SendToTarget(fixmsg, mc.session); err != nil {
				log.Panic(err)
			}

		}
	})
}

func (m *MCSimulator) Close() {
	m.initiator.Stop()
}

func (m *MCSimulator) Name() string {
	return MarketConnectorName
}
