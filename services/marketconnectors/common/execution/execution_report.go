// Common FIX ExecutionReport message processor
package execution

import (
	logger "github.com/apex/log"
	orderCore "github.com/cyanly/gotrade/core/order"
	exeCore "github.com/cyanly/gotrade/core/order/execution"
	proto "github.com/cyanly/gotrade/proto/order"

	"log"
)

// Common behaviours to persist and publish populated Execution entity into our data layer and message bus
func ProcessExecutionReport(er *proto.Execution) {

	er.OrderKey, er.OrderId = exeCore.GetOrderIdentsByClientOrdId(er.ClientOrderId)
	if er.OrderId <= 0 || er.OrderKey <= 0 {
		log.Panic("Unrecognised Order Ident: ", er.ClientOrderId)
	}

	switch er.ExecType {
	case proto.Execution_TRADE_CANCEL:
	case proto.Execution_TRADE_CORRECT:
		//TODO: TryGetAmendedExecution
	case proto.Execution_RESTATED:
	case proto.Execution_REJECTED:
	}

	if err := exeCore.InsertExecution(er); err == nil {

		//TODO: update order status

		//Publish Execution
		logger.Infof("Execution Report: %v", er)
		if data, err := er.Marshal(); err != nil {
			log.Panic(err)
		} else {
			orderCore.MessageBus.Publish("order.Execution", data)
		}

	} else {
		log.Panic(err)
	}
}
