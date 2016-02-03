// Command line executable entry package for FIX buy side simulator
package main

import (
	"github.com/cyanly/gotrade/core/service"
	"github.com/cyanly/gotrade/services/marketconnectors/simulator"
)

func main() {

	// Load configurations

	// Initialise Service Infrastructure
	sc := service.NewConfig()
	sc.ServiceName = "Simulator"
	svc := service.NewService(sc)

	// Initialise Database Connection

	// Initialise Component
	orc := simulator.NewConfig()
	orc.MessageBusURL = sc.MessageBusURL

	orsvc := simulator.NewMarketConnector(orc)
	orsvc.Start()
	defer orsvc.Close()

	// Go
	<-svc.Start()
}
