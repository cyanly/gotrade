// Core order execution report APIs
package execution

import (
	order "github.com/cyanly/gotrade/core/order"
	proto "github.com/cyanly/gotrade/proto/order"

	"fmt"
	"log"
	"time"
)

func NewStatusExecution(o *proto.Order, execType proto.Execution_ExecType, text string) {

	execution := &proto.Execution{
		OrderId:  o.OrderId,
		OrderKey: o.OrderKey,
	}
	clientOrderId := fmt.Sprintf("%v.%v", *o.OrderKey, *o.Version)
	execution.ClientOrderId = &clientOrderId

	//BrokerOrderId := fmt.Sprintf("SIM-%v", *o.OrderId)
	execution.BrokerOrderId = &clientOrderId

	BrokerExecId := fmt.Sprintf("%v", len(o.Executions)+1)
	execution.BrokerExecId = &BrokerExecId

	execution.ExecBroker = o.MarketConnector

	execution.ExecType = &execType

	var OrderStatus proto.OrderStatus
	switch execType {
	case proto.Execution_ORDER_STATUS:
		OrderStatus = proto.OrderStatus_ORDER_SENT
		if *o.Instruction == proto.Order_CANCEL {
			OrderStatus = proto.OrderStatus_MC_SENT_CANCEL
		}
		if *o.Instruction == proto.Order_REPLACE {
			OrderStatus = proto.OrderStatus_MC_SENT_REPLACE
		}
	case proto.Execution_NEW:
		OrderStatus = proto.OrderStatus_NEW
	case proto.Execution_PENDING_CANCEL:
		OrderStatus = proto.OrderStatus_PENDING_CANCEL
	case proto.Execution_PENDING_REPLACE:
		OrderStatus = proto.OrderStatus_PENDING_REPLACE
	case proto.Execution_CANCELED:
		OrderStatus = proto.OrderStatus_CANCELLED
	case proto.Execution_REPLACE:
		OrderStatus = proto.OrderStatus_REPLACED
	case proto.Execution_REJECTED:
		OrderStatus = proto.OrderStatus_REJECTED
	}
	execution.OrderStatus = &OrderStatus

	BrokerExecDatetime := time.Now().UTC().Format(time.RFC3339Nano)
	execution.BrokerExecDatetime = &BrokerExecDatetime

	if text != "" {
		execution.Text = &text
	}

	if err := InsertExecution(execution); err == nil {
		o.Executions = append(o.Executions, execution)
		o.OrderStatus = execution.OrderStatus
		isComplete := order.IsCompleted(o)
		o.IsComplete = &isComplete
		if err = order.UpdateOrderStatus(o); err == nil {
			//Publis Execution
			data, _ := execution.Marshal()
			order.MessageBus.Publish("order.Execution", data)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

}

func NewTradeExecution(o *proto.Order, fillQty float64, fillPrice float64, text string) {
	execution := &proto.Execution{
		OrderId:  o.OrderId,
		OrderKey: o.OrderKey,
	}
	clientOrderId := fmt.Sprintf("%v.%v", *o.OrderKey, *o.Version)
	execution.ClientOrderId = &clientOrderId

	//BrokerOrderId := fmt.Sprintf("SIM-%v", *o.OrderId)
	execution.BrokerOrderId = &clientOrderId

	BrokerExecId := fmt.Sprintf("%v", len(o.Executions)+1)
	execution.BrokerExecId = &BrokerExecId

	execution.ExecBroker = o.MarketConnector

	execType := proto.Execution_TRADE
	execution.ExecType = &execType
	OrderStatus := proto.OrderStatus_PARTIALLY_FILLED
	*o.FilledQuantity += fillQty
	if *o.FilledQuantity >= *o.Quantity {
		OrderStatus = proto.OrderStatus_FILLED
	}
	execution.OrderStatus = &OrderStatus

	execution.Quantity = &fillQty
	execution.Price = &fillPrice

	cumQty := fillQty
	cumAvgPrice := fillPrice * fillQty
	for i := 0; i < len(o.Executions); i++ {
		if o.Executions[i].Quantity != nil {
			cumQty += *o.Executions[i].Quantity

			if o.Executions[i].Price != nil {
				cumAvgPrice += *o.Executions[i].Price * *o.Executions[i].Quantity
			}
		}
	}
	cumAvgPrice /= cumQty
	execution.CumQuantity = &cumQty
	execution.AvgPrice = &cumAvgPrice
	execution.CalcCumQuantity = &cumQty
	execution.CalcAvgPrice = &cumAvgPrice

	BrokerExecDatetime := time.Now().UTC().Format(time.RFC3339Nano)
	execution.BrokerExecDatetime = &BrokerExecDatetime

	if text != "" {
		execution.Text = &text
	}

	if err := InsertExecution(execution); err == nil {
		o.Executions = append(o.Executions, execution)
		o.OrderStatus = execution.OrderStatus
		isComplete := order.IsCompleted(o)
		o.IsComplete = &isComplete
		if err = order.UpdateOrderStatus(o); err == nil {
			//Publish Execution
			data, _ := execution.Marshal()
			order.MessageBus.Publish("order.Execution", data)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}
