// Command line executable entry package for SellSideSimulator
package main

import (
	"github.com/cyanly/gotrade/core/service"
	"github.com/cyanly/gotrade/services/marketconnectors/sellsidesim"
)

func main() {

	// Load configurations

	// Initialise Service Infrastructure
	sc := service.NewConfig()
	sc.ServiceName = "OrderRouter"
	svc := service.NewService(sc)

	// Initialise Component
	orsvc := sellsidesim.NewSellSideSimulator("console")
	orsvc.Start()
	defer orsvc.Close()

	// Go
	<-svc.Start()
}
