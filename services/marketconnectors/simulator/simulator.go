// MC Simulator communicates with SellSideSim service in FIX protocol as a mean to test trade life cycle
package simulator

import (
	log "github.com/cyanly/gotrade/core/logger"
	proto "github.com/cyanly/gotrade/proto/order"
	"github.com/cyanly/gotrade/services/marketconnectors/common"
	"github.com/cyanly/gotrade/services/marketconnectors/common/order"
)

const (
	MarketConnectorName string = "Simulator"
)

type MCSimulator struct {
	Config     common.Config
	OrdersList map[int32]*proto.Order

	app *order.FIXClient
}

func NewMarketConnector(c common.Config) *MCSimulator {
	c.MarketConnectorName = MarketConnectorName

	mc := &MCSimulator{
		Config: c,
	}

	// Simulator is an standard Order only FIX app
	mc.app = order.NewFIXClient(c)

	return mc
}

func (m *MCSimulator) Start() {
	if err := m.app.Start(); err != nil {
		log.WithError(err).Fatal("FIX Client Start Error")
	}
}

func (m *MCSimulator) Close() {
	m.app.Stop()
}

func (m *MCSimulator) Name() string {
	return MarketConnectorName
}
