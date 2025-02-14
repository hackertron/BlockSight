package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	BlocksProcessed       prometheus.Counter
	TransactionsProcessed prometheus.Counter
	ProcessingErrors      prometheus.Counter
	ProcessingLatency     prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		BlocksProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "blockchain_blocks_processed_total",
			Help: "The total number of processed blocks",
		}),
		TransactionsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "blockchain_transactions_processed_total",
			Help: "The total number of processed transactions",
		}),
		ProcessingErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "blockchain_processing_errors_total",
			Help: "The total number of processing errors",
		}),
		ProcessingLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "blockchain_processing_latency_seconds",
			Help:    "Time taken to process blocks",
			Buckets: prometheus.DefBuckets,
		}),
	}
}
