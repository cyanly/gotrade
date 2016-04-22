package memstore

import (
	"errors"

	"fmt"
	"github.com/cyanly/gotrade/database"
	proto "github.com/cyanly/gotrade/proto/order"
)

const StoreType = "memstore"

func init() {
	database.RegisterOrderStore(StoreType, database.OrderStoreRegistration{
		DialFunc: newOrderStore,
	})
}

type OrderStore struct {
	nextOrderID     int32
	nextOrderKey    int32
	nextExecutionID int32
	orderMap        map[int32]*proto.Order
	executionMap    map[int32]*proto.Execution
}

func newOrderStore(url string, options database.Options) (database.OrderStore, error) {
	return &OrderStore{
		orderMap:        make(map[int32]*proto.Order),
		executionMap:    make(map[int32]*proto.Execution),
		nextOrderID:     1,
		nextOrderKey:    1,
		nextExecutionID: 1,
	}, nil
}

//  Interface Implementations ========================================================

// Returns an order for OrderId
func (store *OrderStore) OrderGet(orderId int32) (*proto.Order, error) {
	if order, ok := store.orderMap[orderId]; ok {
		orderClone := &proto.Order{}
		Copy(orderClone, order)
		return orderClone, nil
	} else {
		return nil, errors.New(fmt.Sprintf("OrderStore: OrderID  %v not found", orderId))
	}
}

// Returns an order for OrderId
func (store *OrderStore) OrderGetByOrderKey(orderKey int32) (*proto.Order, error) {
	for _, o := range store.orderMap {
		if o.OrderKey == orderKey {
			orderClone := &proto.Order{}
			Copy(orderClone, o)
			return orderClone, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("OrderStore: OrderKey %v not found", orderKey))
}

// Save order as new entry
func (store *OrderStore) OrderCreate(order *proto.Order) error {
	clone := order
	//&proto.Order{} Copy(clone, order)
	clone.OrderId = store.nextOrderID
	store.nextOrderID++
	if clone.OrderKey == 0 {
		clone.OrderKey = store.nextOrderKey
		store.nextOrderKey++
	} else {
		clone.Version++
	}
	store.orderMap[clone.OrderId] = clone

	return nil
}

// Create an Execution from entity
func (store *OrderStore) ExecutionCreate(er *proto.Execution) error {
	if o, ok := store.orderMap[er.OrderId]; ok == false {
		return errors.New(fmt.Sprintf("OrderStore: OrderID %v does not exists", er.OrderId))
	} else {
		clone := er
		//&proto.Execution{} Copy(clone, er)
		clone.ExecutionId = store.nextExecutionID
		store.nextExecutionID++
		store.executionMap[clone.ExecutionId] = clone

		// Update Order status, this should be in single transaction
		o.FilledQuantity = er.CalcCumQuantity
		o.FilledAvgPrice = er.CalcAvgPrice
		o.OrderStatus = er.OrderStatus
		switch o.OrderStatus {
		case proto.OrderStatus_CANCELLED,
			proto.OrderStatus_REJECTED,
			proto.OrderStatus_FILLED,
			proto.OrderStatus_DONE_FOR_DAY,
			proto.OrderStatus_EXPIRED:
			o.IsComplete = true

		default:
			o.IsComplete = false
		}

	}

	return nil
}

func (store *OrderStore) Close() {}
