package database

import (
	"errors"

	proto "github.com/cyanly/gotrade/proto/order"
)

// FOR CLIENT =================================================================

// Defines the OrderStore interface.
// Every backing store must implement at least this interface.
type OrderStore interface {
	// Returns an order row for OrderId,  excludes Execution/Allocations etc
	OrderGet(int32) (*proto.Order, error)
	// Returns an order row for OrderKey,  excludes Execution/Allocations etc
	OrderGetByOrderKey(int32) (*proto.Order, error)
	// Save order as new entry
	OrderCreate(*proto.Order) error

	// Create an Execution from entity
	ExecutionCreate(*proto.Execution) error

	// Close the store and clean up. (Flush to disk, cleanly sever connections, etc)
	Close()
}

// Datastore options to be passed by client
type Options map[string]interface{}

func NewOrderStore(name, dbpath string, opts Options) (OrderStore, error) {
	r, registered := storeRegistry[name]
	if !registered {
		return nil, errors.New("OrderStore: name '" + name + "' is not registered")
	}
	return r.DialFunc(dbpath, opts)
}

// FOR Storage Implementations =====================================================

type OrderStoreRegistration struct {
	DialFunc func(string, Options) (OrderStore, error)
}

var storeRegistry = make(map[string]OrderStoreRegistration)

func RegisterOrderStore(name string, register OrderStoreRegistration) {
	if _, found := storeRegistry[name]; found {
		panic("already registered OrderStore " + name)
	}
	storeRegistry[name] = register
}
