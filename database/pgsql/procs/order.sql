-- !Select GetOrderByOrderId
-- $1: orderId int32
-- $ret: *Order()
SELECT
    order_id,
    client_guid,
    order_key,
    order_key_version,
    to_char(order_submitted_time at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
    instruction,
    market_connector,
    order_type,
    time_in_force,
    handl_inst,
    symbol,
    exchange,
    side,
    qty,
    limit_price,
    filled_qty,
    filled_avg_price,
    order_status_id,
    is_complete,
    is_booked,
    is_expired,
    trade_booking_id,
    trader_id,
    account,
    broker_user_id,
    broker_account,
    description,
    source,
    open_close,
    algo
from orders WHERE order_id = $1;


-- !Execute UpdateOrderStatus
-- $1: order *Order
UPDATE orders SET
        filled_qty = $1.FilledQuantity,
        filled_avg_price = $1.FilledAvgPrice,
        order_status_id = $1.OrderStatus,
        is_complete = $1.IsComplete,
        is_booked = COALESCE($1.IsBooked, 0),
        is_expired = COALESCE($1.IsExpired, 0),
        trade_booking_id = $1.TradeBookingId
where order_id = $1.OrderId;

-- !SelectOne GetNextOrderKey
-- $ret: nextOrderKey int32
SELECT nextval('orderkeysequence')::INT;

-- !InsertOne InsertOrder
-- $1: order *Order
-- $ret: $1.OrderId
INSERT INTO orders (
    client_guid,
    order_key,
    order_key_version,
    order_submitted_time,
    instruction,
    market_connector,
    order_type,
    time_in_force,
    handl_inst,
    symbol,
    exchange,
    side,
    qty,
    limit_price,
    filled_qty,
    filled_avg_price,
    order_status_id,
    is_complete,
    is_booked,
    is_expired,
    trade_booking_id,
    trader_id,
    account,
    broker_user_id,
    broker_account,
    description,
    source,
    open_close,
    algo )
        VALUES (
            $1.ClientGuid,
            $1.OrderKey,
            $1.Version,
            TIMESTAMP WITH TIME ZONE '$1.SubmitDatetime',
            $1.Instruction,
            $1.MarketConnector,
            $1.OrderType,
            $1.Timeinforce,
            $1.HandleInst,
            $1.Symbol,
            $1.Exchange,
            $1.Side,
            $1.Quantity,
            $1.LimitPrice,
            $1.FilledQuantity,
            $1.FilledAvgPrice,
            $1.OrderStatus,
            $1.IsComplete,
            $1.IsBooked,
            $1.IsExpired,
            $1.TradeBookingId,
            $1.TraderId,
            $1.AccountId,
            $1.BrokerUserid,
            $1.BrokerAccount,
            $1.Description,
            $1.Source,
            $1.OpenClose,
            $1.Algo
        )
RETURNING order_id;

-- !SelectOne GetOrderIdentsByClientOrdId
-- $1: clOrdId string
-- $ret: ordKey int32, ordId int32
SELECT order_key, order_id FROM orders
WHERE order_key = $1 AND order_key_version = $2;