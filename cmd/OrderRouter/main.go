// Command line executable entry package for Order Router
package main

import (
	"github.com/cyanly/gotrade/core/service"
	"github.com/cyanly/gotrade/services/orderrouter"
)

func main() {

	// Load configurations

	// Initialise Service Infrastructure
	sc := service.NewConfig()
	sc.ServiceName = "OrderRouter"
	svc := service.NewService(sc)

	// Initialise Database Connection

	// Initialise OrderRouter
	orc := orderrouter.NewConfig()
	orc.ServiceMessageBusURL = sc.MessageBusURL

	orsvc := orderrouter.NewOrderRouter(orc)
	orsvc.Start()
	defer orsvc.Close()

	// Go
	<-svc.Start()
}
