// A simulated order book bid offer matching by bloomberg pricefeed quotes
//   this is to give our sell side FIX engine simulate a false feel of reality,
//   so as to support MARKET or LIMIT orders
package orderbook

import (
	orac "github.com/cyanly/gotrade/core/order"
	er "github.com/cyanly/gotrade/core/order/execution"
	ems "github.com/cyanly/gotrade/proto/order"
	marketdata "github.com/cyanly/gotrade/proto/pricing"
	"github.com/nats-io/nats"

	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type OrderBook struct {
	OrdersList map[int32]*ems.Order
	Symbol     string
	Quote      *marketdata.Quote
	Subscriber *nats.Subscription
}

type QuoteBand struct {
	price    float64
	quantity float64
	next     *QuoteBand
}

var (
	//Public Variables
	PXMessageBus    *nats.Conn
	OrderUpdateChan chan *ems.Order

	//Private Variables
	closed             chan bool
	newOrderChan       chan *ems.Order
	removeOrderChan    chan *ems.Order
	rejectListenerChan chan *OrderBook
	listeningSymbols   map[string]*OrderBook
)

func init() {
	closed = make(chan bool)
	newOrderChan = make(chan *ems.Order)
	removeOrderChan = make(chan *ems.Order)
	rejectListenerChan = make(chan *OrderBook)
	listeningSymbols = make(map[string]*OrderBook)

	fmt.Println("BBGPX @ " + fmt.Sprint(os.Getenv("RMQ_URL")) + "/BBGPX")
	PXMessageBus, _ := nats.Connect(fmt.Sprint(os.Getenv("RMQ_URL")) + "/BBGPX")

	OrderUpdateChan = make(chan *ems.Order)

	//Imagine how to get below right in C#
	go func() {
		for {
			select {
			case <-closed:
				return
			case order := <-newOrderChan:

				if listener, exists := listeningSymbols[order.Symbol]; exists == false {
					quote := &marketdata.Quote{}
					listener := &OrderBook{
						OrdersList: make(map[int32]*ems.Order),
						Symbol:     order.Symbol,
						//			Subscriber: subscriber,
						Quote: quote,
					}
					listeningSymbols[order.Symbol] = listener
					listener.OrdersList[order.OrderId] = order

					go func(order *ems.Order, listener *OrderBook) {
						quoteInitialReq := &marketdata.QuoteInitialRequest{
							Symbol: &order.Symbol,
						}
						data, _ := quoteInitialReq.Marshal()
						initMsg, err := PXMessageBus.Request("Schemas.Pricefeed.QuoteInitialRequest", data, 5*time.Second)
						if err != nil {
							quoteInitialRes := &marketdata.QuoteInitialResponse{}
							quoteInitialRes.Unmarshal(initMsg.Data)
							if quoteInitialRes.Quote != nil {
								listener.Quote = quoteInitialRes.Quote
							}
						} else {
							rejectListenerChan <- listener
							return
						}

						if len(listener.OrdersList) == 0 {
							return
						}
						subscriber, _ := PXMessageBus.Subscribe(order.Symbol, func(m *nats.Msg) {
							//Simulate trade filling triggered only by a trade in market
							quote := &marketdata.Quote{}
							quote.Unmarshal(m.Data)
							if err := listener.Quote.Unmarshal(m.Data); err == nil {
								listener.ProcessOrderBookUpdate(quote)
							}

						})
						listener.Subscriber = subscriber

						log.Println("SUB " + order.Symbol)
					}(order, listener)
				} else {
					listener.OrdersList[order.OrderId] = order
				}

			case listener := <-rejectListenerChan:
				first := true
				for _, order := range listener.OrdersList {
					er.NewStatusExecution(order, ems.Execution_REJECTED, "Missing Market Data")
					OrderUpdateChan <- order
					if first {
						delete(listeningSymbols, order.Symbol)
						log.Println("UNSUB " + order.Symbol)
						first = false
					}
				}
				listener.Subscriber.Unsubscribe()

			case order := <-removeOrderChan:
				if listener, exists := listeningSymbols[order.Symbol]; exists == true {
					if order, exists := listener.OrdersList[order.OrderId]; exists == true {
						delete(listener.OrdersList, order.OrderId)
					}
					if len(listener.OrdersList) == 0 {
						listener.Subscriber.Unsubscribe()
						delete(listeningSymbols, order.Symbol)
						log.Println("UNSUB " + order.Symbol)
					}
				}
			}
		}
	}()
}

func Close() {
	closed <- true
	PXMessageBus.Close()
}

func RegisterOrder(order *ems.Order) {
	order.FilledQuantity = float64(0)
	order.FilledAvgPrice = float64(0)

	newOrderChan <- order
}

func UnRegisterOrder(order *ems.Order) {
	removeOrderChan <- order
}

func (m *OrderBook) ProcessOrderBookUpdate(triggerQuote *marketdata.Quote) {
	bidQuotes, askQuotes := m.ConstrucOrderBookLinkedList()

	for _, order := range m.OrdersList {
		if orac.IsCompleted(order) {
			continue
		}

		var quotebands *QuoteBand
		if order.Side == ems.Side_BUY {
			quotebands = askQuotes
		} else {
			quotebands = bidQuotes
		}
		if quotebands == nil {
			continue
		}

		if order.OrderType == ems.OrderType_MARKET {
			for quotebands != nil && order.Quantity > order.FilledQuantity {
				fillQty := math.Min(quotebands.quantity, (order.Quantity - order.FilledQuantity))
				DoTrade(order, fillQty, quotebands.price)
				quotebands.quantity -= fillQty
				if quotebands.quantity == 0 {
					quotebands = quotebands.next
				}
			}
		}
		if triggerQuote.LastPrice == nil {
			continue
		}
		if order.OrderType == ems.OrderType_LIMIT {
			for quotebands != nil && order.Quantity > order.FilledQuantity {
				if order.Side == ems.Side_BUY {
					if order.LimitPrice < quotebands.price {
						break
					}
				} else {
					if order.LimitPrice > quotebands.price {
						break
					}
				}
				fillQty := math.Min(quotebands.quantity, (order.Quantity - order.FilledQuantity))
				DoTrade(order, fillQty, quotebands.price)
				quotebands.quantity -= fillQty
				if quotebands.quantity == 0 {
					quotebands = quotebands.next
				}
			}
		}
	}
}

func DoTrade(order *ems.Order, quantity float64, price float64) {
	log.Println(fmt.Sprintf("fill %v: %v @ %v", order.OrderId, quantity, price))
	er.NewTradeExecution(order, quantity, price, "")

	if orac.IsCompleted(order) {
		UnRegisterOrder(order)
		OrderUpdateChan <- order
	}
}

func (m *OrderBook) ConstrucOrderBookLinkedList() (*QuoteBand, *QuoteBand) {
	var bidQuotes *QuoteBand
	if m.Quote.Bid != nil && m.Quote.BidSize != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid,
			quantity: float64(*m.Quote.BidSize),
		}
		bidQuotes = quo
	}
	if m.Quote.Bid1 != nil && m.Quote.Bid1Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid1,
			quantity: float64(*m.Quote.Bid1Size),
		}
		bidQuotes = quo
	}
	if m.Quote.Bid2 != nil && m.Quote.Bid2Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid2,
			quantity: float64(*m.Quote.Bid2Size),
		}
		if bidQuotes == nil {
			bidQuotes = quo
		} else {
			bidQuotes.next = quo
		}
	}
	if m.Quote.Bid3 != nil && m.Quote.Bid3Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid3,
			quantity: float64(*m.Quote.Bid3Size),
		}
		if bidQuotes == nil {
			bidQuotes = quo
		} else {
			bidQuotes.next = quo
		}
	}
	if m.Quote.Bid4 != nil && m.Quote.Bid4Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid4,
			quantity: float64(*m.Quote.Bid4Size),
		}
		if bidQuotes == nil {
			bidQuotes = quo
		} else {
			bidQuotes.next = quo
		}
	}
	if m.Quote.Bid5 != nil && m.Quote.Bid5Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Bid5,
			quantity: float64(*m.Quote.Bid5Size),
		}
		if bidQuotes == nil {
			bidQuotes = quo
		} else {
			bidQuotes.next = quo
		}
	}

	var askQuotes *QuoteBand
	if m.Quote.Ask != nil && m.Quote.AskSize != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask,
			quantity: float64(*m.Quote.AskSize),
		}
		askQuotes = quo
	}
	if m.Quote.Ask1 != nil && m.Quote.Ask1Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask1,
			quantity: float64(*m.Quote.Ask1Size),
		}
		askQuotes = quo
	}
	if m.Quote.Ask2 != nil && m.Quote.Ask2Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask2,
			quantity: float64(*m.Quote.Ask2Size),
		}
		if askQuotes == nil {
			askQuotes = quo
		} else {
			askQuotes.next = quo
		}
	}
	if m.Quote.Ask3 != nil && m.Quote.Ask3Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask3,
			quantity: float64(*m.Quote.Ask3Size),
		}
		if askQuotes == nil {
			askQuotes = quo
		} else {
			askQuotes.next = quo
		}
	}
	if m.Quote.Ask4 != nil && m.Quote.Ask4Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask4,
			quantity: float64(*m.Quote.Ask4Size),
		}
		if askQuotes == nil {
			askQuotes = quo
		} else {
			askQuotes.next = quo
		}
	}
	if m.Quote.Ask5 != nil && m.Quote.Ask5Size != nil {
		quo := &QuoteBand{
			price:    *m.Quote.Ask5,
			quantity: float64(*m.Quote.Ask5Size),
		}
		if askQuotes == nil {
			askQuotes = quo
		} else {
			askQuotes.next = quo
		}
	}
	return bidQuotes, askQuotes
}
