package simulator

import (
	messagebus "github.com/nats-io/nats"
	"github.com/cyanly/gotrade/core/messagebus/test"
	_ "github.com/cyanly/gotrade/database/memstore"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/services/marketconnectors/common"
	"github.com/cyanly/gotrade/services/marketconnectors/sellsidesim"
	testOrder "github.com/cyanly/gotrade/test/order"

	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//mock message bus
	ts := test.RunDefaultServer()
	defer ts.Shutdown()

	//simulate a sell side FIX server
	sellsvc := sellsidesim.NewSellSideSimulator("")
	sellsvc.Start()
	defer sellsvc.Close()

	code := m.Run()
	os.Exit(code)
}

func TestNewMarketConnectorStartAndStop(t *testing.T) {
	//config
	sc := common.NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.DatabaseDriver = "memstore"

	//start MC
	svc := NewMarketConnector(sc)
	svc.Start()

	//stop and disconnect MC
	time.Sleep(100 * time.Millisecond)
	svc.Close()

	//expect no panic
}

func TestNewMarketConnectorNewOrderRequest(t *testing.T) {
	//config
	sc := common.NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.DatabaseDriver = "memstore"

	//start MC
	svc := NewMarketConnector(sc)
	svc.Start()

	time.Sleep(100 * time.Millisecond)

	// mock a NewOrderRequest, really should come from OrderRouter
	request := &proto.NewOrderRequest{
		Order: testOrder.MockOrder(),
	}
	request.Order.OrderStatus = proto.OrderStatus_ORDER_RECEIVED
	request.Order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano)
	request.Order.MessageType = proto.Order_NEW
	request.Order.OrderKey = int32(321)
	svc.app.OrderStore.OrderCreate(request.Order)

	// subscribe to check if we can receive a fill execution report from sim broker
	recvExecutionReport := false
	svc.app.MessageBus.Subscribe("order.Execution", func(m *messagebus.Msg) {
		recvExecutionReport = true
		exec := new(proto.Execution)
		if err := exec.Unmarshal(m.Data); err == nil {
			if exec.ClientOrderId != "321.1" {
				t.Fatalf("unexpected execution report ClOrdId %v, expecting 321.1", exec.ClientOrderId)
			}
			if exec.OrderId != 1 {
				t.Fatalf("unexpected execution report OrderId %v, expecting 1", exec.OrderId)
			}
		} else {
			t.Fatalf("unexpected execution report: %v", err)
		}
	})

	// send above mock NewOrderRequest
	data, _ := request.Marshal()
	svc.app.MessageBus.Publish("order.NewOrderRequest.MC."+MarketConnectorName, data)

	time.Sleep(100 * time.Millisecond)

	//stop and disconnect MC
	svc.Close()

	if recvExecutionReport == false {
		t.Fatal("No execution report received from mock order new request")
	}

}
