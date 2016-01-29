package order

import (
	proto "github.com/cyanly/gotrade/proto/order"

	"encoding/json"
	"log"
	"time"
)

func MockOrder() *proto.Order {
	var o proto.Order
	if err := json.Unmarshal([]byte(`{
"order_id": 0,
"order_key": 0,
"version": 0,
"client_guid": "1234-5678-abcd-3210",
"instruction": "NEW",
"side": "BUY",
"quantity": 100,
"symbol": "AAPL",
"order_type": "LIMIT",
"limit_price": 100,
"timeinforce": "DAY",
"handle_inst": "AUTOMATED_EXECUTION_ORDER_PRIVATE",
"exchange": "CME",
"description": "Apple Inc.",
"account_id": 1,
"market_connector": "TestCase",
"source": "mock test",
"trader": "trader1",
"trader_id": 1,
"machine": "CYAN",
"memo": "Super Secret",
"broker_account": "billions_club",
"create_datetime": "`+time.Now().UTC().Format(time.RFC3339Nano)+`",
"-": "###"
}`), &o); err != nil {
		log.Fatal("MockOrder: ", err)
	}

	return &o
}
