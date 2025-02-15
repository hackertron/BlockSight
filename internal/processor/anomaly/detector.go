package anomaly

import (
	"math"
	"sync"
	"time"

	"github.com/hackertron/blocksight/internal/models"
	"github.com/hackertron/blocksight/internal/utils"
)

type DetectionRule struct {
	Metric    string
	Threshold float64
	Window    time.Duration
}

type Detector struct {
	rules       []DetectionRule
	metrics     map[string][]float64
	timeWindows map[string][]time.Time
	mu          sync.RWMutex
}

func NewDetector() *Detector {
	return &Detector{
		rules: []DetectionRule{
			{
				Metric:    "gas_price",
				Threshold: 2.0, // Standard deviations
				Window:    time.Minute * 5,
			},
			{
				Metric:    "block_time",
				Threshold: 2.0,
				Window:    time.Minute * 5,
			},
			{
				Metric:    "transaction_count",
				Threshold: 2.0,
				Window:    time.Minute * 5,
			},
		},
		metrics:     make(map[string][]float64),
		timeWindows: make(map[string][]time.Time),
	}
}

func (d *Detector) AddMetric(metric string, value float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	if _, exists := d.metrics[metric]; !exists {
		d.metrics[metric] = make([]float64, 0)
		d.timeWindows[metric] = make([]time.Time, 0)
	}

	d.metrics[metric] = append(d.metrics[metric], value)
	d.timeWindows[metric] = append(d.timeWindows[metric], now)

	// Clean old data
	d.cleanOldData(metric)
}

func (d *Detector) cleanOldData(metric string) {
	rule := d.findRule(metric)
	if rule == nil {
		return
	}

	cutoff := time.Now().Add(-rule.Window)
	var i int
	for i = 0; i < len(d.timeWindows[metric]); i++ {
		if d.timeWindows[metric][i].After(cutoff) {
			break
		}
	}

	if i > 0 {
		d.metrics[metric] = d.metrics[metric][i:]
		d.timeWindows[metric] = d.timeWindows[metric][i:]
	}
}

func (d *Detector) DetectAnomalies(block *models.Block) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var anomalies []string

	avgGasPrice := utils.CalculateAverageGasPrice(block)
	d.AddMetric("gas_price", avgGasPrice)
	if d.isAnomaly("gas_price", avgGasPrice) {
		anomalies = append(anomalies, "gas_price")
	}

	if len(d.timeWindows["block_time"]) > 0 {
		blockTime := time.Since(d.timeWindows["block_time"][len(d.timeWindows["block_time"])-1]).Seconds()
		d.AddMetric("block_time", blockTime)
		if d.isAnomaly("block_time", blockTime) {
			anomalies = append(anomalies, "block_time")
		}
	}

	txCount := float64(len(block.Transactions))
	d.AddMetric("transaction_count", txCount)
	if d.isAnomaly("transaction_count", txCount) {
		anomalies = append(anomalies, "transaction_count")
	}

	return anomalies
}

func (d *Detector) isAnomaly(metric string, value float64) bool {
	rule := d.findRule(metric)
	if rule == nil {
		return false
	}

	values := d.metrics[metric]
	if len(values) < 10 { // Need minimum sample size
		return false
	}

	mean := utils.CalculateMean(values)
	stdDev := utils.CalculateStdDev(values, mean)

	zScore := math.Abs(value-mean) / stdDev
	return zScore > rule.Threshold
}

func (d *Detector) findRule(metric string) *DetectionRule {
	for _, rule := range d.rules {
		if rule.Metric == metric {
			return &rule
		}
	}
	return nil
}
