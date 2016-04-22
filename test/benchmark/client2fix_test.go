package benchmark

import (
	"github.com/nats-io/nats"
	logger "github.com/cyanly/gotrade/core/logger"
	"github.com/cyanly/gotrade/core/messagebus/test"
	"github.com/cyanly/gotrade/core/service"
	_ "github.com/cyanly/gotrade/database/memstore"
	proto "github.com/cyanly/gotrade/proto/order"
	MCCommon "github.com/cyanly/gotrade/services/marketconnectors/common"
	"github.com/cyanly/gotrade/services/marketconnectors/sellsidesim"
	MCSimulator "github.com/cyanly/gotrade/services/marketconnectors/simulator"
	"github.com/cyanly/gotrade/services/orderrouter"
	testOrder "github.com/cyanly/gotrade/test/order"

	"log"
	"os"
	"testing"
	"time"
)

var (
	msgbus *nats.Conn
)

func TestMain(m *testing.M) {
	// Start a temporary messaging broker
	ts := test.RunDefaultServer()
	defer ts.Shutdown()

	// Start OrderRouter service
	orSC := orderrouter.NewConfig()
	orSC.MessageBusURL = "nats://localhost:22222"
	orSC.ServiceMessageBusURL = "nats://localhost:22222"
	orSC.DatabaseDriver = "memstore"
	orSvc := orderrouter.NewOrderRouter(orSC)
	orSvc.Start()
	defer orSvc.Close()

	// Start Simulated FIX Sell Side
	sellsvc := sellsidesim.NewSellSideSimulator("")
	sellsvc.Start()
	defer sellsvc.Close()

	// Start a MarketConnector Service for simuilated market
	// 1. MC (implementation)
	mcSC := MCCommon.NewConfig()
	mcSC.MessageBusURL = "nats://localhost:22222"
	mcSC.DatabaseDriver = "memstore"
	mcSvc := MCSimulator.NewMarketConnector(mcSC)
	mcSvc.Start()
	defer mcSvc.Close()
	// 2. Service (heartbeating etc)
	sc := service.NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceName = "MC." + mcSvc.Name()
	svc := service.NewService(sc)
	svc.Start()
	defer svc.Stop()

	// Turn off logging to measure performance
	logger.Discard()

	time.Sleep(100 * time.Millisecond) // async

	// Connect messaging bus for this mock client
	if nc, err := nats.Connect("nats://localhost:22222"); err != nil {
		log.Fatal("error: Cannot connect to message bus")
	} else {
		msgbus = nc
	}

	code := m.Run()
	os.Exit(code)
}

func BenchmarkTradeflow(b *testing.B) {
	b.StopTimer()

	// Mock new order request
	request := &proto.NewOrderRequest{
		Order: testOrder.MockOrder(),
	}
	request.Order.MarketConnector = "Simulator"
	data, merr := request.Marshal()
	if merr != nil {
		b.Fatal(merr)
	}

	////////
	// benchmark the whole journey of:
	// CL->OR:   Client send order protobuf to OrderRouter(OR)
	// OR->MC:   OrderRouter process order and dispatch persisted order entity or target MarketConnector
	// MC->FIX:  MarketConnector translate into NewOrderSingle FIX message based on the session with its counterparty
	// FIX->MC:  MarketConnector received FIX message on its order, here Simulator sending a fully FILL execution
	// EXEC->CL: MarketConnector publish processed and persisted Execution onto messaging bus, here our Client will listen to
	///////
	log.Println("benchmark count: ", b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if reply, err := msgbus.Request("order.NewOrderRequest", data, 2000*time.Millisecond); err != nil {
			b.Fatalf("an error '%s' was not expected when sending NewOrderRequest", err)
		} else {
			resp := new(proto.NewOrderResponse)
			if err := resp.Unmarshal(reply.Data); err != nil {
				b.Fatalf("an error '%s' was not expected when decoding NewOrderResponse", err)
			}

			if resp.ErrorMessage != nil && len(*resp.ErrorMessage) > 0 {
				b.Fatalf("unexpected NewOrderResponse error message: %s", *resp.ErrorMessage)
			}
		}
	}
}
