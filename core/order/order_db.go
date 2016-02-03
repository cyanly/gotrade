package order

import (
	logger "github.com/apex/log"
	proto "github.com/cyanly/gotrade/proto/order"

	"database/sql"
)

var (
	DB *sql.DB
)

func GetOrderByOrderId(id int32) (*proto.Order, error) {
	logger.Infof("sql: SELECT orders WHERE order_id = %v", id)
	rows, err := DB.Query(`
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
	 from orders WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &proto.Order{}
		var OrderId int32
		var ClientGuid string
		var OrderKey int32
		var OrderKeyVersion int32
		var OrderSubmittedTime string
		var MessageType string
		var MarketConnector string
		var OrderType string
		var TimeInForce string
		var HandlInst string
		var Symbol string
		var Exchange string
		var Side string
		var Qty float64
		var LimitPrice sql.NullFloat64
		var FilledQty sql.NullFloat64
		var FilledAvgPrice sql.NullFloat64
		var OrderStatusId int32
		var IsComplete bool
		var IsBooked bool
		var IsExpired bool
		var TradeBookingId sql.NullInt64
		var TraderId int32
		var Account int32
		var BrokerUserID string
		var BrokerAccount string
		var Description string
		var Source string
		var OpenClose string
		var AlgoOne string

		if err := rows.Scan(
			&OrderId,
			&ClientGuid,
			&OrderKey,
			&OrderKeyVersion,
			&OrderSubmittedTime,
			&MessageType,
			&MarketConnector,
			&OrderType,
			&TimeInForce,
			&HandlInst,
			&Symbol,
			&Exchange,
			&Side,
			&Qty,
			&LimitPrice,
			&FilledQty,
			&FilledAvgPrice,
			&OrderStatusId,
			&IsComplete,
			&IsBooked,
			&IsExpired,
			&TradeBookingId,
			&TraderId,
			&Account,
			&BrokerUserID,
			&BrokerAccount,
			&Description,
			&Source,
			&OpenClose,
			&AlgoOne,
		); err != nil {
			return nil, err
		}

		order.OrderId = &OrderId
		order.ClientGuid = &ClientGuid
		order.OrderKey = &OrderKey
		order.Version = &OrderKeyVersion
		order.SubmitDatetime = &OrderSubmittedTime
		orderinstruct := proto.Order_OrderInstruction(proto.Order_OrderInstruction_value[MessageType])
		order.Instruction = &orderinstruct
		order.MarketConnector = &MarketConnector
		ordType := proto.OrderType(proto.OrderType_value[OrderType])
		order.OrderType = &ordType
		timeinforce := proto.TimeInForce(proto.TimeInForce_value[TimeInForce])
		order.Timeinforce = &timeinforce
		handlinst := proto.HandlInst(proto.HandlInst_value[HandlInst])
		order.HandleInst = &handlinst
		order.Symbol = &Symbol
		order.Exchange = &Exchange
		side := proto.Side(proto.Side_value[Side])
		order.Side = &side
		order.Quantity = &Qty
		if LimitPrice.Valid {
			order.LimitPrice = &LimitPrice.Float64
		}
		if FilledQty.Valid {
			order.FilledQuantity = &FilledQty.Float64
		}
		if FilledAvgPrice.Valid {
			order.FilledAvgPrice = &FilledAvgPrice.Float64
		}
		ordStatus := proto.OrderStatus(OrderStatusId)
		order.OrderStatus = &ordStatus
		order.IsComplete = &IsComplete
		order.IsBooked = &IsBooked
		order.IsExpired = &IsExpired
		if TradeBookingId.Valid {
			tradebookid := int32(TradeBookingId.Int64)
			order.TradeBookingId = &tradebookid
		}
		order.TraderId = &TraderId
		order.AccountId = &Account
		order.BrokerUserid = &BrokerUserID
		order.BrokerAccount = &BrokerAccount
		order.Description = &Description
		order.Source = &Source
		order.OpenClose = &OpenClose
		order.Algo = &AlgoOne

		if err := Validate(order); err != nil {
			return nil, err
		}

		return order, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func GetOrderByOrderKey(key int32) (*proto.Order, error) {
	logger.Infof("sql: SELECT orders WHERE order_key = %v", key)
	rows, err := DB.Query(`
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
	 from orders WHERE order_key = $1`, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &proto.Order{}
		var OrderId int32
		var ClientGuid string
		var OrderKey int32
		var OrderKeyVersion int32
		var OrderSubmittedTime string
		var MessageType string
		var MarketConnector string
		var OrderType string
		var TimeInForce string
		var HandlInst string
		var Symbol string
		var Exchange string
		var Side string
		var Qty float64
		var LimitPrice sql.NullFloat64
		var FilledQty sql.NullFloat64
		var FilledAvgPrice sql.NullFloat64
		var OrderStatusId int32
		var IsComplete bool
		var IsBooked bool
		var IsExpired bool
		var TradeBookingId sql.NullInt64
		var TraderId int32
		var Account int32
		var BrokerUserID string
		var BrokerAccount string
		var Description string
		var Source string
		var OpenClose string
		var AlgoOne string

		if err := rows.Scan(
			&OrderId,
			&ClientGuid,
			&OrderKey,
			&OrderKeyVersion,
			&OrderSubmittedTime,
			&MessageType,
			&MarketConnector,
			&OrderType,
			&TimeInForce,
			&HandlInst,
			&Symbol,
			&Exchange,
			&Side,
			&Qty,
			&LimitPrice,
			&FilledQty,
			&FilledAvgPrice,
			&OrderStatusId,
			&IsComplete,
			&IsBooked,
			&IsExpired,
			&TradeBookingId,
			&TraderId,
			&Account,
			&BrokerUserID,
			&BrokerAccount,
			&Description,
			&Source,
			&OpenClose,
			&AlgoOne,
		); err != nil {
			return nil, err
		}

		order.OrderId = &OrderId
		order.ClientGuid = &ClientGuid
		order.OrderKey = &OrderKey
		order.Version = &OrderKeyVersion
		order.SubmitDatetime = &OrderSubmittedTime
		orderinstruct := proto.Order_OrderInstruction(proto.Order_OrderInstruction_value[MessageType])
		order.Instruction = &orderinstruct
		order.MarketConnector = &MarketConnector
		ordType := proto.OrderType(proto.OrderType_value[OrderType])
		order.OrderType = &ordType
		timeinforce := proto.TimeInForce(proto.TimeInForce_value[TimeInForce])
		order.Timeinforce = &timeinforce
		handlinst := proto.HandlInst(proto.HandlInst_value[HandlInst])
		order.HandleInst = &handlinst
		order.Symbol = &Symbol
		order.Exchange = &Exchange
		side := proto.Side(proto.Side_value[Side])
		order.Side = &side
		order.Quantity = &Qty
		if LimitPrice.Valid {
			order.LimitPrice = &LimitPrice.Float64
		}
		if FilledQty.Valid {
			order.FilledQuantity = &FilledQty.Float64
		}
		if FilledAvgPrice.Valid {
			order.FilledAvgPrice = &FilledAvgPrice.Float64
		}
		ordStatus := proto.OrderStatus(OrderStatusId)
		order.OrderStatus = &ordStatus
		order.IsComplete = &IsComplete
		order.IsBooked = &IsBooked
		order.IsExpired = &IsExpired
		if TradeBookingId.Valid {
			tradebookid := int32(TradeBookingId.Int64)
			order.TradeBookingId = &tradebookid
		}
		order.TraderId = &TraderId
		order.AccountId = &Account
		order.BrokerUserid = &BrokerUserID
		order.BrokerAccount = &BrokerAccount
		order.Description = &Description
		order.Source = &Source
		order.OpenClose = &OpenClose
		order.Algo = &AlgoOne

		if err := Validate(order); err != nil {
			return nil, err
		}

		return order, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil //errors.New("No Result")
}

func GetNextOrderKey() (int32, error) {
	logger.Info("sql: SELECT nextval('orderkeysequence')")

	var nextOrderKey int32
	if err := DB.QueryRow(`SELECT nextval('orderkeysequence')::INT`).Scan(&nextOrderKey); err != nil {
		return 0, err
	}

	logger.Infof("sql ret: ID = %d", nextOrderKey)

	return nextOrderKey, nil
}

func InsertOrder(order *proto.Order) error {
	logger.Info("sql: INSERT INTO orders ...")
	var lastId int

	if err := DB.QueryRow(`
	INSERT INTO orders (
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
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				TIMESTAMP WITH TIME ZONE '$5',
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
				$18,
				$19,
				$20,
				$21,
				$22,
				$23,
				$24,
				$25,
				$26,
				$27,
				$28,
				$29
			)
	RETURNING order_id
`,
		order.ClientGuid,
		order.OrderKey,
		order.Version,
		order.SubmitDatetime,
		order.Instruction,
		order.MarketConnector,
		order.OrderType,
		order.Timeinforce,
		order.HandleInst,
		order.Symbol,
		order.Exchange,
		order.Side,
		order.Quantity,
		order.LimitPrice,
		order.FilledQuantity,
		order.FilledAvgPrice,
		order.OrderStatus,
		order.IsComplete,
		order.IsBooked,
		order.IsExpired,
		order.TradeBookingId,
		order.TraderId,
		order.AccountId,
		order.BrokerUserid,
		order.BrokerAccount,
		order.Description,
		order.Source,
		order.OpenClose,
		order.Algo,
	).Scan(&lastId); err != nil {
		return err
	}

	logger.Infof("sql ret: ID = %d", lastId)

	lastId32 := int32(lastId)
	order.OrderId = &lastId32

	return nil
}

func UpdateOrderStatus(order *proto.Order) error {
	logger.Infof("sql: UPDATE orders WHERE order_id = %v", *order.OrderId)
	if _, err := DB.Exec(`
		UPDATE orders SET
				filled_qty = $1,
				filled_avg_price = $2,
				order_status_id = $3,
				is_complete = $4,
				is_booked = COALESCE($5, 0),
				is_expired = COALESCE($6, 0),
				trade_booking_id = $7
		where order_id = $8
`,
		order.FilledQuantity,
		order.FilledAvgPrice,
		order.OrderStatus,
		order.IsComplete,
		order.IsExpired,
		order.IsBooked,
		order.TradeBookingId,
		order.OrderId); err != nil {
		return err
	}

	return nil
}
