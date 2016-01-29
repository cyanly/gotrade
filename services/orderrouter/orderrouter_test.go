package orderrouter

import (
	order "github.com/cyanly/gotrade/core/order"
	proto "github.com/cyanly/gotrade/proto/order"
	pService "github.com/cyanly/gotrade/proto/service"
	testOrder "github.com/cyanly/gotrade/test/order"
	gnatsd "github.com/nats-io/gnatsd/test"

	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/erikstmartin/go-testdb"
	"log"
	"strings"
	"testing"
	"time"
)

type testResult struct {
	lastId       int64
	affectedRows int64
}

func (r testResult) LastInsertId() (int64, error) {
	return r.lastId, nil
}

func (r testResult) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}

func TestOrderRouterStartAndStop(t *testing.T) {
	gnatsd.DefaultTestOptions.Port = 22222
	ts := gnatsd.RunDefaultServer()
	defer ts.Shutdown()

	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceMessageBusURL = "nats://localhost:22222"

	svc := NewOrderRouter(sc)
	svc.Start()

	svc.Close()
	time.Sleep(200 * time.Millisecond)

	if svc.msgbus.IsClosed() != true {
		t.Fatal("order router failed to shut down")
	}
}

func TestOrderRouterNewOrderRequest(t *testing.T) {
	gnatsd.DefaultTestOptions.Port = 22222
	ts := gnatsd.RunDefaultServer()
	defer ts.Shutdown()

	svc := mockOrderRouter()

	orderKey := int(0)
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (result driver.Rows, err error) {
		columns := []string{"id", "name", "age", "created"}
		rows := "unknown"

		if strings.Contains(query, "INSERT INTO orders") {
			columns = []string{"order_id"}
			rows = "123"
		}

		if strings.Contains(query, "INSERT INTO execution") {
			columns = []string{"execution_id"}
			rows = "111"
		}

		if strings.Contains(query, "SELECT nextval('orderkeysequence')::INT") {
			columns = []string{"orderkeysequence"}
			orderKey++
			rows = fmt.Sprint(orderKey)
		}

		if rows == "unknown" {
			log.Println(query)
		}

		//if args[0] == "joe" {
		//	rows = "2,joe,25,2012-10-02 02:00:02"
		//}
		return testdb.RowsFromCSVString(columns, rows), nil
	})
	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

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
		Name:   &mcname,
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

			if resp.Order == nil || resp.Order.OrderId == nil || *resp.Order.OrderId != 123 {
				t.Fatalf("unexpected NewOrderResponse id not matching mock OrderId(123)")
			}

			if resp.Order == nil || resp.Order.OrderKey == nil || *resp.Order.OrderKey != 2 {
				t.Fatalf("unexpected NewOrderResponse OrderKey %v, expecting 2", *resp.Order.OrderKey)
			}
		}
	}
}

func mockOrderRouter() *OrderRouter {
	db, _ := sql.Open("testdb", "")
	order.DB = db

	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceMessageBusURL = "nats://localhost:22222"

	svc := NewOrderRouter(sc)

	return svc
}
