package orderrouter

import (
	"os"
	"testing"
	"time"

	"github.com/cyanly/gotrade/core/messagebus/test"
	_ "github.com/cyanly/gotrade/database/memstore"
	proto "github.com/cyanly/gotrade/proto/order"
	pService "github.com/cyanly/gotrade/proto/service"
	testOrder "github.com/cyanly/gotrade/test/order"
)

func TestMain(m *testing.M) {
	//mock message bus
	ts := test.RunDefaultServer()
	defer ts.Shutdown()

	code := m.Run()
	os.Exit(code)
}

func TestOrderRouterStartAndStop(t *testing.T) {

	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceMessageBusURL = "nats://localhost:22222"
	sc.DatabaseDriver = "memstore"

	svc := NewOrderRouter(sc)
	svc.Start()

	svc.Close()
	time.Sleep(200 * time.Millisecond)

	if svc.msgbus.IsClosed() != true {
		t.Fatal("order router failed to shut down")
	}
}

func TestOrderRouterNewOrderRequest(t *testing.T) {

	svc := mockOrderRouter()
	svc.Start()
	defer svc.Close()

	// mock new order request
	request := &proto.NewOrderRequest{
		Order: testOrder.MockOrder(),
	}
	data, merr := request.Marshal()
	if merr != nil {
		t.Fatal(merr)
	}

	// send mock order, expecting reject due to market connector not up
	req := new(proto.NewOrderRequest)
	if err := req.Unmarshal(data); err != nil {
		t.Fatal(err)
	}
	if reply, err := svc.msgbus.Request("order.NewOrderRequest", data, 200*time.Millisecond); err != nil {
		t.Fatalf("an error '%s' was not expected when sending NewOrderRequest", err)
	} else {
		resp := new(proto.NewOrderResponse)
		if err := resp.Unmarshal(reply.Data); err != nil {
			t.Fatalf("an error '%s' was not expected when decoding NewOrderResponse", err)
		}

		if *resp.ErrorMessage != "LINK TO BROKER DOWN" {
			t.Fatalf("unexpected NewOrderResponse %s, expecting LINK TO BROKER DOWN", *resp.ErrorMessage)
		}
	}

	// mock a test market connector
	mcname := "MC.TestCase"
	hbMsg := pService.Heartbeat{
		Name:   mcname,
		Status: pService.RUNNING,
	}
	if hbMsgData, err := hbMsg.Marshal(); err != nil {
		t.Fatal(err)
	} else {
		svc.msgbusService.Publish("service.Heartbeat.MC.TestCase", hbMsgData)

		time.Sleep(100 * time.Millisecond)

		if reply, err := svc.msgbus.Request("order.NewOrderRequest", data, 200*time.Millisecond); err != nil {
			t.Fatalf("an error '%s' was not expected when sending NewOrderRequest", err)
		} else {
			resp := new(proto.NewOrderResponse)
			if err := resp.Unmarshal(reply.Data); err != nil {
				t.Fatalf("an error '%s' was not expected when decoding NewOrderResponse", err)
			}

			if resp.ErrorMessage != nil && len(*resp.ErrorMessage) > 0 {
				t.Fatalf("unexpected NewOrderResponse error message: %s", *resp.ErrorMessage)
			}

			if resp.Order == nil || resp.Order.OrderId != 2 {
				t.Fatalf("unexpected NewOrderResponse id not matching mock OrderId(%v), expecting 2")
			}

			if resp.Order == nil || resp.Order.OrderKey != 2 {
				t.Fatalf("unexpected NewOrderResponse OrderKey %v, expecting 2", resp.Order.OrderKey)
			}
		}
	}
}

func mockOrderRouter() *OrderRouter {

	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceMessageBusURL = "nats://localhost:22222"
	sc.DatabaseDriver = "memstore"

	svc := NewOrderRouter(sc)

	return svc
}
