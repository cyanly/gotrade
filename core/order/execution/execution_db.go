package execution

import (
	logger "github.com/apex/log"
	order "github.com/cyanly/gotrade/core/order"
	proto "github.com/cyanly/gotrade/proto/order"

	"strconv"
	"strings"
)

func InsertExecution(exe *proto.Execution) error {
	logger.Info("sql: INSERT INTO execution ...")
	var lastId int

	if err := order.DB.QueryRow(`
			INSERT INTO execution (
				order_id,
				order_status_id,
				broker_execution_time,
				qty,
				cum_qty,
				price,
				avg_price,
				broker_order_id,
				broker_exec_id,
				previous_broker_exec_id,
				calc_cum_qty,
				calc_avg_price,
				exec_type_id,
				last_mkt,
				last_liquidity_ind_id,
				text,
				exec_broker,
				cancel_replace_by_exececution_id
			)
			VALUES (
				$1,
				$2,
				TIMESTAMP WITH TIME ZONE '$3',
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				$14,
				$15,
				$16,
				$17,
				$18
			)
			RETURNING execution_id
`,
		exe.OrderId,
		exe.OrderStatus,
		exe.BrokerExecDatetime,
		exe.Quantity,
		exe.CumQuantity,
		exe.Price,
		exe.AvgPrice,
		exe.BrokerOrderId,
		exe.BrokerExecId,
		exe.PrevBrokerExecId,
		exe.CalcCumQuantity,
		exe.CalcAvgPrice,
		exe.ExecType,
		exe.Lastmkt,
		exe.LastLiquidity,
		exe.Text,
		exe.ExecBroker,
		exe.CancelReplaceByExececutionId,
	).Scan(&lastId); err != nil {
		return err
	}

	logger.Infof("sql ret: ID = %d", lastId)

	lastId32 := int32(lastId)
	exe.ExecutionId = lastId32

	return nil
}

func GetOrderIdentsByClientOrdId(clOrdId string) (ordKey int32, ordId int32) {
	keyVer := strings.Split(clOrdId, ".")
	if len(keyVer) != 2 {
		return
	}

	if oKey, err := strconv.Atoi(keyVer[0]); err == nil {
		ordKey = int32(oKey)
	}
	ordVer, _ := strconv.Atoi(keyVer[1])

	logger.Info("sql: SELECT order_id FROM orders")

	if err := order.DB.QueryRow(`
	SELECT order_id FROM orders
	WHERE order_key = $1
	AND order_key_version = $2`, ordKey, ordVer).Scan(&ordId); err != nil {
		return
	}

	logger.Infof("sql ret: ID = %d", ordId)

	return
}
