// Serving order New/Cancel/Replace tasks, communicates with various market connectors
package orderrouter

import (
	messagebus "github.com/nats-io/nats"
	log "github.com/cyanly/gotrade/core/logger"
	"github.com/cyanly/gotrade/core/order"
	"github.com/cyanly/gotrade/database"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/proto/service"

	"fmt"
	"strings"
	"time"
)

type OrderRouter struct {
	Config Config

	orderStore    database.OrderStore     // order storage driver
	stopChan      chan bool               // signal stop order router processor channel
	msgbusService *messagebus.Conn        // message bus listening to service heartbeats
	msgbus        *messagebus.Conn        // message bus listenning to order
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
	OrderId     int32
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

	// Connect to database storage driver
	if storage, err := database.NewOrderStore(c.DatabaseDriver, c.DatabaseUrl, nil); err != nil {
		log.Fatalf("error: Cannot connect to database driver %v @ %v", c.DatabaseDriver, c.DatabaseUrl)
	} else {
		or.orderStore = storage
	}

	// Keep a list of active MarketConnectors
	if ncSvc, err := messagebus.Connect(c.ServiceMessageBusURL); err != nil {
		log.Fatalf("error: Cannot connect to service message bus @ %v", c.ServiceMessageBusURL)
	} else {
		or.msgbusService = ncSvc
		ncSvc.Subscribe("service.Heartbeat.MC.>", func(m *messagebus.Msg) {
			hbMsg := &service.Heartbeat{}
			if err := hbMsg.Unmarshal(m.Data); err == nil {
				or.mcChan <- hbMsg
			}
		})
	}

	// Order requests processors subscribing
	if msgbus, err := messagebus.Connect(c.MessageBusURL); err != nil {
		log.Fatalf("error: Cannot connect to order message bus @ %v", c.MessageBusURL)
	} else {
		or.msgbus = msgbus

		//CL->OR order NEW request
		msgbus.Subscribe("order.NewOrderRequest", func(m *messagebus.Msg) {
			request := new(proto.NewOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {
				// validate basic fields
				if request.Order == nil {
					// empty request
					return
				}
				if request.Order.MarketConnector == "" {
					// empty market connect
					return
				}

				orReq := OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_NEW,
					Request:     request,
				}
				oO := order.Order{
					Order: request.Order,
				}
				log.Infof("CL->OR NEW %v", oO.String())
				or.reqChan <- orReq
			}
		})

		//CL->OR order CANCEL request
		msgbus.Subscribe("order.CancelOrderRequest", func(m *messagebus.Msg) {
			request := new(proto.CancelOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {
				// validate basic fields
				if request.OrderId == 0 {
					// request without order id
					return
				}

				orReq := OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_CANCEL,
					Request:     request,
					OrderId:     request.OrderId,
				}

				log.Infof("CL->OR CXL OrderKey := %v OrderId := %v", request.OrderKey, request.OrderId)
				or.reqChan <- orReq
			}
		})

		//CL->OR order REPLACE request
		msgbus.Subscribe("order.ReplaceOrderRequest", func(m *messagebus.Msg) {
			request := new(proto.ReplaceOrderRequest)
			if err := request.Unmarshal(m.Data); err == nil && len(m.Reply) > 0 {
				// validate basic fields
				if request.Order == nil {
					// empty request
					return
				}
				if request.Order.OrderId == 0 {
					// request without order id
					return
				}
				if request.Order.MarketConnector == "" {
					// empty market connect
					return
				}

				orReq := OrderRequest{
					ReplyAddr:   m.Reply,
					RequestType: REQ_REPLACE,
					Request:     request,
					OrderId:     request.Order.OrderId,
				}

				log.Infof("CL->OR RPL OrderKey := %v OrderId := %v", request.Order.OrderKey, request.Order.OrderId)
				or.reqChan <- orReq

			}
		})
	}

	return or
}

// start the logic spinning code for OrderRouter, using a single for..select.. pattern so there is no need to lock resources
// hence to avoid synchronisation issues
func (self *OrderRouter) Start() {
	go func() {
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
				func() { // wrap in func for defer reply

					requestError := ""
					var order *proto.Order

					// make sure reply message sent in the end
					defer func() {
						var data []byte
						errCode := int32(0)
						if requestError != "" {
							errCode = int32(-1)
						}
						switch req.RequestType {
						case REQ_NEW:
							resp := &proto.NewOrderResponse{
								ErrorCode:    errCode,
								ErrorMessage: &requestError,
								Order:        order,
							}
							data, _ = resp.Marshal()
						case REQ_CANCEL:
							resp := &proto.CancelOrderResponse{
								ErrorCode:    errCode,
								ErrorMessage: &requestError,
								Order:        order,
							}
							data, _ = resp.Marshal()
						case REQ_REPLACE:
							resp := &proto.ReplaceOrderResponse{
								ErrorCode:    errCode,
								ErrorMessage: &requestError,
								Order:        order,
							}
							data, _ = resp.Marshal()
						}
						self.msgbus.Publish(req.ReplyAddr, data)
					}()

					// Prepare order
					if req.RequestType == REQ_CANCEL || req.RequestType == REQ_REPLACE {
						// retrieve previous order for state check
						if prev_order, err := self.orderStore.OrderGet(req.OrderId); err != nil {
							requestError = "Order does not exists"
							return
						} else {
							// rule.1: working orders
							if prev_order.MessageType == proto.Order_NEW {
								if !(prev_order.OrderStatus == proto.OrderStatus_PARTIALLY_FILLED ||
								prev_order.OrderStatus == proto.OrderStatus_NEW) {
									requestError = fmt.Sprintf("Disallowed due to order status: %s", prev_order.OrderStatus)
									return
								}
							}
							// rule.2: pending cancel rejected
							if prev_order.MessageType == proto.Order_CANCEL {
								if !(prev_order.OrderStatus == proto.OrderStatus_REJECTED) {
									requestError = "Pending cancel on order"
									return
								}
							}
							// rule.2: pending replace acked by broker
							if prev_order.MessageType == proto.Order_REPLACE {
								if !(prev_order.OrderStatus == proto.OrderStatus_NEW ||
								prev_order.OrderStatus == proto.OrderStatus_REPLACED ||
								prev_order.OrderStatus == proto.OrderStatus_REJECTED) {
									requestError = "Pending replace on order"
									return
								}
							}

							// prepare for new order entry
							if req.RequestType == REQ_CANCEL {
								pReq := req.Request.(*proto.CancelOrderRequest)
								order = prev_order
								order.OrderId = 0                                                // wipe order id for new id
								order.Version++                                                  // bump up order key version
								order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano) // timing this cancel
								order.Trader = pReq.Trader                                       // trader who tries to cancel
								order.TraderId = pReq.TraderId                                   // trader who tries to cancel
								order.Source = pReq.Source                                       // source of cancel
								order.MessageType = proto.Order_CANCEL
								order.OrderStatus = proto.OrderStatus_CANCEL_RECEIVED
							}
							if req.RequestType == REQ_REPLACE {
								pReq := req.Request.(*proto.ReplaceOrderRequest)
								pReq.Order.OrderKey = prev_order.OrderKey

								order := pReq.Order
								order.OrderId = 0                                                // wipe order id for new id
								order.OrderKey = prev_order.OrderKey                             // in case client failed to provide correct order key
								order.Version = prev_order.Version + 1                           // bump up order key version
								order.SubmitDatetime = time.Now().UTC().Format(time.RFC3339Nano) // timing this replace
								order.Trader = pReq.Trader                                       // trader who tries to replace
								order.TraderId = pReq.TraderId                                   // trader who tries to replace
								order.MessageType = proto.Order_REPLACE
								order.OrderStatus = proto.OrderStatus_REPLACE_RECEIVED
							}

							//TODO: check ClientGUID consistency
						}
					} else {
						// New order
						pReq := req.Request.(*proto.NewOrderRequest)
						order = pReq.Order
						order.OrderId = 0
						order.Version = 0
					}

					// Persist order entry
					if err := self.orderStore.OrderCreate(order); err != nil {
						requestError = fmt.Sprint(err)
						log.WithError(err).Error("[ OR ] ERROR Create order")
						return
					}

					//TODO: allocations

					// Check target market connector is up
					if _, ok := self.mclist[order.MarketConnector]; ok == false {
						log.Warnf("OR->CL REJECT OrderKey: %v MC: %v : %v", order.OrderKey, order.MarketConnector, "LINK TO BROKER DOWN")

						// REJECT due to market connector down
						requestError = "LINK TO BROKER DOWN"

						// create Reject execution
						self.reject(order, "LINK TO BROKER DOWN")

					} else {
						log.WithField("order", order).Infof("OR->MC OrderKey: %v", order.OrderKey)

						// relay order with idents to its market connector
						// todo: special case on MC.AlgoAggregator
						var data []byte
						switch req.RequestType {
						case REQ_NEW:
							request := req.Request.(*proto.NewOrderRequest)
							data, _ = request.Marshal()
						case REQ_CANCEL:
							request := req.Request.(*proto.CancelOrderRequest)
							data, _ = request.Marshal()
						case REQ_REPLACE:
							request := req.Request.(*proto.ReplaceOrderRequest)
							data, _ = request.Marshal()
						}
						self.msgbus.Publish("order.NewOrderRequest.MC."+order.MarketConnector, data)
					}
				}()
			}
		}
	}()
}

func (self *OrderRouter) Close() {
	self.stopChan <- true
}

func (self *OrderRouter) reject(o *proto.Order, reason string) {

	execution := &proto.Execution{
		OrderId:            o.OrderId,
		OrderKey:           o.OrderKey,
		ClientOrderId:      fmt.Sprintf("%v.%v", o.OrderKey, o.Version),
		ExecType:           proto.Execution_REJECTED,
		OrderStatus:        proto.OrderStatus_REJECTED,
		BrokerExecDatetime: time.Now().UTC().Format(time.RFC3339Nano),
		Text:               reason,
	}

	if err := self.orderStore.ExecutionCreate(execution); err == nil {
		//publish to message bus
		data, _ := execution.Marshal()
		self.msgbus.Publish("order.Execution", data)
	} else {
		log.WithField("execution", execution).WithError(err).Error("[ OR ] ERROR Create Reject Execution")
	}

}
