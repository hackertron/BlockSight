package log_analyzer

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/pkg/models"
)

// LogAnalyzer represents the log analysis functionality
type LogAnalyzer struct {
	client *ethclient.Client
}

// NewLogAnalyzer creates a new instance of LogAnalyzer
func NewLogAnalyzer(client *ethclient.Client) *LogAnalyzer {
	return &LogAnalyzer{client: client}
}

// AnalyzeLogs analyzes the logs and returns a slice of Log models
func (la *LogAnalyzer) AnalyzeLogs(blockNumber uint64) ([]models.Log, error) {
	block, err := la.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}

	var analyzedLogs []models.Log
	for _, tx := range block.Transactions() {
		receipt, err := la.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Printf("failed to retrieve transaction receipt: %v", err)
			continue
		}

		for _, log := range receipt.Logs {
			analyzedLog := models.Log{
				Address:     log.Address.Hex(),
				Topics:      parseTopics(log.Topics),
				Data:        log.Data,
				BlockNumber: log.BlockNumber,
				TxHash:      log.TxHash.Hex(),
				TxIndex:     log.TxIndex,
				BlockHash:   log.BlockHash.Hex(),
				Index:       log.Index,
				Removed:     log.Removed,
			}
			analyzedLogs = append(analyzedLogs, analyzedLog)
		}
	}

	return analyzedLogs, nil
}

// parseTopics parses the topics and returns a slice of topic strings
func parseTopics(topics []common.Hash) []string {
	var parsedTopics []string
	for _, topic := range topics {
		parsedTopics = append(parsedTopics, topic.Hex())
	}
	return parsedTopics
}
