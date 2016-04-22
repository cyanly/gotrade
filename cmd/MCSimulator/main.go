// Command line executable entry package for FIX buy side simulator
package main

import (
	"github.com/cyanly/gotrade/core/service"
	_ "github.com/cyanly/gotrade/database/memstore"
	"github.com/cyanly/gotrade/services/marketconnectors/common"
	"github.com/cyanly/gotrade/services/marketconnectors/simulator"
)

func main() {

	// Load configurations

	// Initialise Service Infrastructure
	sc := service.NewConfig()
	sc.ServiceName = "Simulator"
	svc := service.NewService(sc)

	// Initialise Component
	orc := common.NewConfig()
	orc.MessageBusURL = sc.MessageBusURL

	// Initialise Database Connection
	orc.DatabaseDriver = "memstore"

	orsvc := simulator.NewMarketConnector(orc)
	orsvc.Start()
	defer orsvc.Close()

	// Go
	<-svc.Start()
}
