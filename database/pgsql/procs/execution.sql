-- !InsertOne InsertExecution
-- $1: exe *Execution
-- $ret: $1.ExecutionId
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
		$1.OrderId,
		$1.OrderStatus,
		TIMESTAMP WITH TIME ZONE '$1.BrokerExecDatetime',
		$1.Quantity,
		$1.CumQuantity,
		$1.Price,
		$1.AvgPrice,
		$1.BrokerOrderId,
		$1.BrokerExecId,
		$1.PrevBrokerExecId,
		$1.CalcCumQuantity,
		$1.CalcAvgPrice,
		$1.ExecType,
		$1.Lastmkt,
		$1.LastLiquidity,
		$1.Text,
		$1.ExecBroker,
		$1.CancelReplaceByExececutionId
    )
RETURNING execution_id;