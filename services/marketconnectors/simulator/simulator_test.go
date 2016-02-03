package simulator

import (
	logger "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	order "github.com/cyanly/gotrade/core/order"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/services/marketconnectors/sellsidesim"
	testOrder "github.com/cyanly/gotrade/test/order"
	gnatsd "github.com/nats-io/gnatsd/test"
	"github.com/nats-io/nats"

	"database/sql"
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logger.SetHandler(cli.Default)

	//mock message bus
	gnatsd.DefaultTestOptions.Port = 22222
	ts := gnatsd.RunDefaultServer()
	defer ts.Shutdown()

	//simulate a sell side FIX server
	sellsvc := sellsidesim.NewSellSideSimulator("")
	sellsvc.Start()
	defer sellsvc.Close()

	//mock db
	db, _ := sql.Open("testdb", "")
	order.DB = db
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (result driver.Rows, err error) {
		columns := []string{"id", "name", "age", "created"}
		rows := "unknown"

		if strings.Contains(query, "INSERT INTO execution") {
			columns = []string{"execution_id"}
			rows = "111"
		}

		if strings.Contains(query, "SELECT order_id FROM orders") {
			columns = []string{"order_id"}
			rows = "123"
		}

		if rows == "unknown" {
			log.Println(query)
		}

		//if args[0] == "joe" {
		//	rows = "2,joe,25,2012-10-02 02:00:02"
		//}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	code := m.Run()
	os.Exit(code)
}

func TestNewMarketConnectorStartAndStop(t *testing.T) {
	//config
	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"

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
	sc := NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"

	//start MC
	svc := NewMarketConnector(sc)
	svc.Start()

	time.Sleep(100 * time.Millisecond)

	// mock a NewOrderRequest, really should come from OrderRouter
	request := &proto.NewOrderRequest{
		Order: testOrder.MockOrder(),
	}
	os_OR := proto.OrderStatus_ORDER_RECEIVED
	request.Order.OrderStatus = &os_OR
	o_stime := time.Now().UTC().Format(time.RFC3339Nano)
	request.Order.SubmitDatetime = &o_stime
	o_oi := proto.Order_NEW
	request.Order.Instruction = &o_oi
	oid := int32(123)
	request.Order.OrderId = &oid
	okey := int32(321)
	request.Order.OrderKey = &okey
	oVer := int32(1)
	request.Order.Version = &oVer

	// subscribe to check if we can receive a fill execution report from sim broker
	recvExecutionReport := false
	svc.msgbus.Subscribe("order.Execution", func(m *nats.Msg) {
		recvExecutionReport = true
		exec := new(proto.Execution)
		if err := exec.Unmarshal(m.Data); err == nil {
			if *exec.ClientOrderId != "321.1" {
				t.Fatalf("unexpected execution report ClOrdId %v, expecting 321.1", *exec.ClientOrderId)
			}
			if *exec.OrderId != 123 {
				t.Fatalf("unexpected execution report OrderId %v, expecting 123", *exec.OrderId)
			}
		} else {
			t.Fatalf("unexpected execution report: %v", err)
		}
	})

	// send above mock NewOrderRequest
	data, _ := request.Marshal()
	svc.msgbus.Publish("order.NewOrderRequest.MC."+MarketConnectorName, data)

	time.Sleep(100 * time.Millisecond)

	//stop and disconnect MC
	svc.Close()

	if recvExecutionReport == false {
		t.Fatal("No execution report received from mock order new request")
	}

}
