package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hackertron/blocksight/internal/client"
	"github.com/hackertron/blocksight/internal/config"
	"github.com/hackertron/blocksight/internal/kafka"
	"github.com/hackertron/blocksight/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize metrics
	m := metrics.NewMetrics()

	// Create Alchemy client
	alchemyClient, err := client.NewAlchemyClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Alchemy client: %v", err)
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Start metrics server
	if cfg.Metrics.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Printf("Starting metrics server on :%d", cfg.Metrics.Port)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Metrics.Port), nil); err != nil {
				log.Printf("Metrics server error: %v", err)
			}
		}()
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start processing blocks
	go processBlocks(ctx, alchemyClient, producer, m)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
}

func processBlocks(ctx context.Context, client *client.AlchemyClient, producer *kafka.Producer, m *metrics.Metrics) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			block, err := client.GetLatestBlock(ctx)
			if err != nil {
				log.Printf("Error getting latest block: %v", err)
				m.ProcessingErrors.Inc()
				continue
			}

			if err := producer.PublishBlock(block); err != nil {
				log.Printf("Error publishing block: %v", err)
				m.ProcessingErrors.Inc()
				continue
			}

			m.BlocksProcessed.Inc()
			m.TransactionsProcessed.Add(float64(len(block.Transactions)))
			m.ProcessingLatency.Observe(time.Since(start).Seconds())
		}
	}
}
