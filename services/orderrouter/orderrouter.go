// Serving order New/Cancel/Replace tasks, communicates with various market connectors
package orderrouter

import (
	logger "github.com/apex/log"
	orderCore "github.com/cyanly/gotrade/core/order"
	execCore "github.com/cyanly/gotrade/core/order/execution"
	proto "github.com/cyanly/gotrade/proto/order"
	service "github.com/cyanly/gotrade/proto/service"
	"github.com/nats-io/nats"

	"fmt"
	"log"
	"strings"
	"time"
)

type OrderRouter struct {
	Config Config

	stopChan      chan bool               // signal stop order router processor channel
	msgbusService *nats.Conn              // message bus listening to service heartbeats
	msgbus        *nats.Conn              // message bus listenning to order
	mclist        map[string]time.Time    // list of available market connectors
	mcChan        chan *service.Heartbeat // channel updates market connector heartbeats
	reqChan       chan OrderRequest       // channel for order requests
}

type ReqType int

const (
	REQ_NEW ReqType = iota
	REQ_CANCEL
	REQ_REPLACE
)

type OrderRequest struct {
	ReplyAddr   string
	RequestType ReqType
	Request     interface{}
}

// Initialise OrderRouter instance and set up topic subscribers
func NewOrderRouter(c Config) *OrderRouter {

	or := &OrderRouter{
		Config: c,

		// internal variables
		stopChan: make(chan bool),
		mclist:   map[string]time.Time{},
		mcChan:   make(chan *service.Heartbeat),
		reqChan:  make(chan OrderRequest),
	}

	// Keep a list of active MarketConnectors
	if ncSvc, err := nats.Connect(c.ServiceMessageBusURL); err != nil {
		log.Fatal("error: Cannot connect to service message bus @ ", c.ServiceMessageBusURL)
	} else {
		or.msgbusService = ncSvc
		ncSvc.Subscribe("service.Heartbeat.MC.>", func(m *nats.Msg) {
			hbMsg := &service.Heartbeat{}
			if err := hbMsg.Unmarshal(m.Data); err == nil {
				or.mcChan <- hbMsg
			}
		})
	}

	// Order requests processors subscribing
	if msgbus, err := nats.Connect(c.MessageBusURL); err != nil {
		log.Fatal("error: Cannot connect to order message bus @ ", c.MessageBusURL)
	} else {
		or.msgbus = msgbus
		orderCore.MessageBus = msgbus

		//CL->OR order NEW request
		msgbus.Subscribe("order.NewOrderRequest", func(m *nats.Msg) {
			request := new(proto.NewOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {

				logger.Infof("CL->OR NEW %v", orderCore.Stringify(request.Order))
				or.reqChan <- OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_NEW,
					Request:     request,
				}
			}
		})

		//CL->OR order CANCEL request
		msgbus.Subscribe("order.CancelOrderRequest", func(m *nats.Msg) {
			request := new(proto.CancelOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {
				logger.Infof("CL->OR CXL OrderKey := %v OrderId := %v", request.OrderKey, request.OrderId)
				or.reqChan <- OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_CANCEL,
					Request:     request,
				}
			}
		})

		//CL->OR order REPLACE request
		msgbus.Subscribe("order.ReplaceOrderRequest", func(m *nats.Msg) {
			request := new(proto.ReplaceOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {
				logger.Infof("CL->OR RPL %v", orderCore.Stringify(request.Order))
				or.reqChan <- OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_REPLACE,
					Request:     request,
				}
			}
		})
	}

	return or
}

// start the logic spinning code for OrderRouter, using a single for..select.. pattern so there is no need to lock resources
// hence to avoid thread synchronisation issues
func (self *OrderRouter) Start() {

	//lock free request processing channel
	go func(self *OrderRouter) {
		for {
			select {
			case <-self.stopChan:
				self.msgbus.Close()
				self.msgbusService.Close()
				return

			// update known market connector list based on their heartbeat status
			case hbMsg := <-self.mcChan:
				svcName := strings.Replace(hbMsg.Name, "MC.", "", 1)
				if hbMsg.Status == service.RUNNING {
					self.mclist[svcName] = time.Now()
				} else {
					delete(self.mclist, svcName)
				}
			// remove market connector if we have not heard from it for a while
			case <-time.After(6 * time.Second):
				for name, lastHb := range self.mclist {
					if time.Since(lastHb).Seconds() > 6 {
						delete(self.mclist, name)
					}
				}
			case req := <-self.reqChan:

				switch req.RequestType {

				case REQ_NEW:
					request := req.Request.(*proto.NewOrderRequest)

					func() {
						// begin with order received status
						request.Order.OrderStatus = proto.OrderStatus_ORDER_RECEIVED
						request.Order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano)
						request.Order.Instruction = proto.Order_NEW

						// persist new order
						if orderKey, err := orderCore.GetNextOrderKey(); err != nil {
							log.Panic("sql error: ", err)
						} else {
							request.Order.OrderKey = orderKey
						}
						if err := orderCore.InsertOrder(request.Order); err != nil {
							log.Panic("sql error: ", err)
						}

						// construct response msg
						resp := &proto.NewOrderResponse{
							Order: request.Order,
						}
						resp.ErrorCode = int32(0)
						// make sure reply message sent in the end
						defer func() {
							if data, err := resp.Marshal(); err == nil {
								self.msgbus.Publish(req.ReplyAddr, data)
							}
						}()

						// check target market connector is up
						if _, ok := self.mclist[request.Order.MarketConnector]; ok == false {
							logger.Warnf("OR->CL REJECT:%v : %v", resp.Order.OrderKey, "LINK TO BROKER DOWN")

							// REJECT due to market connector down
							resp.ErrorCode = int32(-1)
							respErrMsg := "LINK TO BROKER DOWN"
							resp.ErrorMessage = &respErrMsg

							// insert Reject execution
							execCore.NewStatusExecution(resp.Order, proto.Execution_REJECTED, "LINK TO BROKER DOWN")
						} else {
							logger.Info("OR->MC NewOrderRequest")

							// relay order with idents to its market connector
							data, _ := request.Marshal()
							self.msgbus.Publish("order.NewOrderRequest.MC."+resp.Order.MarketConnector, data)
						}
					}()

				case REQ_CANCEL:
					request := req.Request.(*proto.CancelOrderRequest)

					func() {
						// construct response msg
						resp := &proto.CancelOrderResponse{}
						resp.ErrorCode = int32(0)
						// make sure reply message sent in the end
						defer func() {
							if data, err := resp.Marshal(); err == nil {
								self.msgbus.Publish(req.ReplyAddr, data)
							}
						}()

						// Retrieve previous order
						if order, err := orderCore.GetOrderByOrderKey(request.OrderKey); err != nil {
							log.Panic("sql error: ", err)
						} else if order == nil {
							resp.ErrorCode = int32(-1)
							respErrMsg := "Order does not exists"
							resp.ErrorMessage = &respErrMsg
						} else {

							// Check if this order is allowed to cancel
							if respErrMsg := isOrderCanCancelReplace(order); len(respErrMsg) > 0 {
								resp.ErrorCode = int32(-1)
								resp.ErrorMessage = &respErrMsg
								return
							}

							// Persist as new CANCEL order
							order.OrderId = 0                                                // wipe order id for new id
							order.Version++                                                  // bump up order key version
							order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano) // timing this cancel
							order.Trader = request.Trader                                    // trader who tries to cancel
							order.TraderId = request.TraderId                                // trader who tries to cancel
							order.Source = request.Source                                    // source of cancel
							order.Instruction = proto.Order_CANCEL
							order.OrderStatus = proto.OrderStatus_ORDER_RECEIVED

							//TODO: allocations

							if err := orderCore.InsertOrder(order); err != nil {
								log.Panic("sql error: ", err)
							}

							// check target market connector is up
							if _, ok := self.mclist[order.MarketConnector]; ok == false {
								logger.Warnf("OR->CL REJECT CXL:%v : %v", resp.Order.OrderKey, "LINK TO BROKER DOWN")

								// REJECT due to market connector down
								resp.ErrorCode = int32(-1)
								respErrMsg := "LINK TO BROKER DOWN"
								resp.ErrorMessage = &respErrMsg

								// insert Reject execution
								execCore.NewStatusExecution(resp.Order, proto.Execution_REJECTED, "LINK TO BROKER DOWN")
							} else {
								logger.Info("OR->MC CancelOrderRequest")

								// relay order with idents to its market connector
								data, _ := request.Marshal()
								self.msgbus.Publish("order.CancelOrderRequest.MC."+resp.Order.MarketConnector, data)
							}
						}

					}()

				case REQ_REPLACE:
					request := req.Request.(*proto.ReplaceOrderRequest)

					func() {
						// construct response msg
						resp := &proto.ReplaceOrderResponse{}
						resp.ErrorCode = int32(0)
						// make sure reply message sent in the end
						defer func() {
							if data, err := resp.Marshal(); err == nil {
								self.msgbus.Publish(req.ReplyAddr, data)
							}
						}()

						// Retrieve previous order
						if prev_order, err := orderCore.GetOrderByOrderKey(request.Order.OrderKey); err != nil {
							log.Panic("sql error: ", err)
						} else if prev_order == nil {
							resp.ErrorCode = int32(-1)
							respErrMsg := "Order does not exists"
							resp.ErrorMessage = &respErrMsg
						} else {

							// Check if this order is allowed to cancel
							if respErrMsg := isOrderCanCancelReplace(prev_order); len(respErrMsg) > 0 {
								resp.ErrorCode = int32(-1)
								resp.ErrorMessage = &respErrMsg
								return
							}

							// Persist new REPLACE order
							order := request.Order
							order.OrderId = 0                                                // wipe order id for new id
							order.Version = prev_order.Version + 1                           // bump up order key version
							order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano) // timing this replace
							order.Instruction = proto.Order_REPLACE
							order.OrderStatus = proto.OrderStatus_ORDER_RECEIVED

							//TODO: allocations

							if err := orderCore.InsertOrder(order); err != nil {
								log.Panic("sql error: ", err)
							}

							// check target market connector is up
							if _, ok := self.mclist[order.MarketConnector]; ok == false {
								logger.Warnf("OR->CL REJECT RPL:%v : %v", resp.Order.OrderKey, "LINK TO BROKER DOWN")

								// REJECT due to market connector down
								resp.ErrorCode = int32(-1)
								respErrMsg := "LINK TO BROKER DOWN"
								resp.ErrorMessage = &respErrMsg

								// insert Reject execution
								execCore.NewStatusExecution(resp.Order, proto.Execution_REJECTED, "LINK TO BROKER DOWN")
							} else {
								logger.Info("OR->MC ReplaceOrderRequest")

								// relay order with idents to its market connector
								data, _ := request.Marshal()
								self.msgbus.Publish("order.ReplaceOrderRequest.MC."+resp.Order.MarketConnector, data)
							}
						}

					}()
				}

			}
		}
	}(self)
}

func (self *OrderRouter) Close() {
	self.stopChan <- true
}

func isOrderCanCancelReplace(order *proto.Order) string {
	if order.Instruction == proto.Order_NEW {
		if !(order.OrderStatus == proto.OrderStatus_PARTIALLY_FILLED ||
			order.OrderStatus == proto.OrderStatus_NEW) {
			return fmt.Sprintf("Order can not be cancelled, current status:%s", order.OrderStatus)
		}
	}
	if order.Instruction == proto.Order_CANCEL {
		if !(order.OrderStatus == proto.OrderStatus_REJECTED) {
			return fmt.Sprintf("Pending cancel on order")
		}
	}
	if order.Instruction == proto.Order_REPLACE {
		if !(order.OrderStatus == proto.OrderStatus_NEW ||
			order.OrderStatus == proto.OrderStatus_REPLACED ||
			order.OrderStatus == proto.OrderStatus_REJECTED) {
			return fmt.Sprintf("Pending replace on order")
		}
	}

	return ""
}
