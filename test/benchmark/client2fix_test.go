package benchmark

import (
	logger "github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	order "github.com/cyanly/gotrade/core/order"
	"github.com/cyanly/gotrade/core/service"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/services/marketconnectors/sellsidesim"
	MCSimulator "github.com/cyanly/gotrade/services/marketconnectors/simulator"
	"github.com/cyanly/gotrade/services/orderrouter"
	testOrder "github.com/cyanly/gotrade/test/order"
	gnatsd "github.com/nats-io/gnatsd/test"
	"github.com/nats-io/nats"

	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/erikstmartin/go-testdb"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	msgbus *nats.Conn
)

func TestMain(m *testing.M) {

	// Start a temporary messaging broker
	gnatsd.DefaultTestOptions.Port = 22222
	ts := gnatsd.RunDefaultServer()
	defer ts.Shutdown()

	// Start a test DB
	db, _ := sql.Open("testdb", "")
	order.DB = db
	mockDB()

	// Start OrderRouter service
	orSC := orderrouter.NewConfig()
	orSC.MessageBusURL = "nats://localhost:22222"
	orSC.ServiceMessageBusURL = "nats://localhost:22222"
	orSvc := orderrouter.NewOrderRouter(orSC)
	orSvc.Start()
	defer orSvc.Close()

	// Start Simulated FIX Sell Side
	sellsvc := sellsidesim.NewSellSideSimulator("")
	sellsvc.Start()
	defer sellsvc.Close()

	// Start a MarketConnector Service for simuilated market
	mcSC := MCSimulator.NewConfig()
	mcSC.MessageBusURL = "nats://localhost:22222"
	mcSvc := MCSimulator.NewMarketConnector(mcSC)
	mcSvc.Start()
	defer mcSvc.Close()
	sc := service.NewConfig()
	sc.MessageBusURL = "nats://localhost:22222"
	sc.ServiceName = "MC." + mcSvc.Name()
	svc := service.NewService(sc)
	svc.Start()
	defer svc.Stop()

	// Turn off logging to measure performance
	logger.SetHandler(discard.Default)

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
	mcName := "Simulator"
	request.Order.MarketConnector = &mcName
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

// Simulate Database behaviours
func mockDB() {

	orderKey := int(0)
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (result driver.Rows, err error) {
		columns := []string{"id", "name", "age", "created"}
		rows := "unknown"

		// Orders
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

		// Executions
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

		return testdb.RowsFromCSVString(columns, rows), nil
	})
	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})
}

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
