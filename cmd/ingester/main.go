package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hackertron/blocksight/internal/cache"
	"github.com/hackertron/blocksight/internal/client"
	"github.com/hackertron/blocksight/internal/config"
	"github.com/hackertron/blocksight/internal/kafka"
	"github.com/hackertron/blocksight/internal/logger"
	"github.com/hackertron/blocksight/internal/metrics"
	"github.com/hackertron/blocksight/internal/models"
	"github.com/hackertron/blocksight/internal/processor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Init(true); err != nil {
		panic(err)
	}
	log := logger.Get()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize components
	blockCache := cache.NewBlockCache(time.Hour)
	metricsCollector := metrics.NewMetrics()

	// Create Alchemy client
	alchemyClient, err := client.NewAlchemyClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Alchemy client", zap.Error(err))
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	// Create batch processor
	batchProcessor := processor.NewBatchProcessor(
		100,
		time.Second*5,
		func(blocks []*models.Block) error {
			return producer.PublishBlocks(blocks)
		},
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start metrics server
	if cfg.Metrics.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Info("Starting metrics server",
				zap.Int("port", cfg.Metrics.Port))

			if err := http.ListenAndServe(
				fmt.Sprintf(":%d", cfg.Metrics.Port),
				nil,
			); err != nil {
				log.Error("Metrics server error", zap.Error(err))
			}
		}()
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start block processing
	go func() {
		if err := processBlocks(ctx, alchemyClient, batchProcessor, blockCache, metricsCollector, log); err != nil {
			log.Error("Block processing error", zap.Error(err))
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Info("Shutting down...")

	// Allow some time for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Wait for processing to complete or timeout
	select {
	case <-ctx.Done():
		log.Info("Graceful shutdown completed")
	case <-shutdownCtx.Done():
		log.Warn("Shutdown timeout exceeded")
	}
}

func processBlocks(
	ctx context.Context,
	client *client.AlchemyClient,
	processor *processor.BatchProcessor,
	cache *cache.BlockCache,
	metrics *metrics.Metrics,
	log *zap.Logger,
) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			start := time.Now()

			// Get latest block
			block, err := client.GetLatestBlock(ctx)
			if err != nil {
				metrics.ProcessingErrors.Inc()
				log.Error("Failed to get latest block", zap.Error(err))
				continue
			}

			// Check if block was already processed
			if cache.Exists(block.Hash) {
				continue
			}

			// Process block
			if err := processor.Add(block); err != nil {
				metrics.ProcessingErrors.Inc()
				log.Error("Failed to process block",
					zap.String("hash", block.Hash.Hex()),
					zap.Error(err))
				continue
			}

			// Update metrics
			metrics.BlocksProcessed.Inc()
			metrics.TransactionsProcessed.Add(float64(len(block.Transactions)))
			metrics.ProcessingLatency.Observe(time.Since(start).Seconds())

			// Add to cache
			cache.Add(block.Hash)

			log.Info("Processed block",
				zap.String("hash", block.Hash.Hex()),
				zap.Int("transactions", len(block.Transactions)),
				zap.Duration("duration", time.Since(start)))
		}
	}
}
