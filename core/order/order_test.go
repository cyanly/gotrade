package order

import (
	proto "github.com/cyanly/gotrade/proto/order"

	"encoding/json"
	"testing"
	"time"
)

func TestOrderJson_Parse(t *testing.T) {
	// Parse configuration.
	var o proto.Order
	if err := json.Unmarshal([]byte(`{
"order_id": 123,
"order_key": 0,
"version": 0,
"client_guid": "1234-5678-abcd-3210",
"instruction": "NEW",
"side": "BUY",
"quantity": 100,
"symbol": "AAPL",
"order_type": "MARKET",
"timeinforce": "DAY",
"handle_inst": "AUTOMATED_EXECUTION_ORDER_PRIVATE",
"exchange": "CME",
"description": "Apple Inc.",
"account_id": 1,
"market_connector": "Simulator",
"source": "mock test",
"trader": "trader1",
"trader_id": 1,
"machine": "CYAN",
"memo": "Super Secret",
"create_datetime": "`+time.Now().UTC().Format(time.RFC3339Nano)+`",
"-": "###"
}`), &o); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if *o.OrderId != 123 {
		t.Fatalf("unexpected order_id: %v, expecting 123", *o.OrderId)
	}
}
