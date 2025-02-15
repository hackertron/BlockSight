package analytics

import (
	"sync"
	"time"

	"github.com/hackertron/blocksight/internal/models"
	"github.com/hackertron/blocksight/internal/utils"
)

type Analytics struct {
	blockTimes []float64
	gasPrices  []float64
	txCounts   []float64
	timeWindow time.Duration
	mu         sync.RWMutex
}

func NewAnalytics(timeWindow time.Duration) *Analytics {
	return &Analytics{
		blockTimes: make([]float64, 0),
		gasPrices:  make([]float64, 0),
		txCounts:   make([]float64, 0),
		timeWindow: timeWindow,
	}
}
func (a *Analytics) GetMetrics() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return map[string]interface{}{
		"average_gas_price": utils.CalculateMean(a.gasPrices),
		"average_tx_count":  utils.CalculateMean(a.txCounts),
		"gas_price_trend":   calculateTrend(a.gasPrices),
		"tx_count_trend":    calculateTrend(a.txCounts),
	}
}

func (a *Analytics) AddBlock(block *models.Block) {
	a.mu.Lock()
	defer a.mu.Unlock()

	avgGasPrice := utils.CalculateAverageGasPrice(block)
	txCount := float64(len(block.Transactions))

	a.gasPrices = append(a.gasPrices, avgGasPrice)
	a.txCounts = append(a.txCounts, txCount)

	a.cleanOldData()
}

func calculateTrend(values []float64) string {
	if len(values) < 2 {
		return "stable"
	}

	first := values[0]
	last := values[len(values)-1]
	change := ((last - first) / first) * 100

	switch {
	case change > 10:
		return "increasing"
	case change < -10:
		return "decreasing"
	default:
		return "stable"
	}
}
func (a *Analytics) cleanOldData() {
	cutoff := time.Now().Add(-a.timeWindow)

	// Remove old gas prices
	var i int
	for i = 0; i < len(a.gasPrices); i++ {
		if time.Now().Add(-a.timeWindow).Before(cutoff) {
			break
		}
	}
	a.gasPrices = a.gasPrices[i:]

	// Remove old transaction counts
	for i = 0; i < len(a.txCounts); i++ {
		if time.Now().Add(-a.timeWindow).Before(cutoff) {
			break
		}
	}
	a.txCounts = a.txCounts[i:]
}
