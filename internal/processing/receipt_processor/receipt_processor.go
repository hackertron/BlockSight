package receipt_processor

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/pkg/models"
)

// ReceiptProcessor represents the receipt processing functionality
type ReceiptProcessor struct {
	client *ethclient.Client
}

// NewReceiptProcessor creates a new instance of ReceiptProcessor
func NewReceiptProcessor(client *ethclient.Client) *ReceiptProcessor {
	return &ReceiptProcessor{client: client}
}

// ProcessReceipt processes the transaction receipt and returns a Receipt model
func (rp *ReceiptProcessor) ProcessReceipt(txHash string) (*models.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := rp.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction receipt: %v", err)
	}

	processedReceipt := &models.Receipt{
		TxHash:          receipt.TxHash.Hex(),
		BlockHash:       receipt.BlockHash.Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
		Status:          receipt.Status,
		Logs:            parseLogs(receipt.Logs),
	}

	return processedReceipt, nil
}

// parseLogs parses the logs and returns a slice of Log models
func parseLogs(logs []*types.Log) []models.Log {
	var parsedLogs []models.Log
	for _, log := range logs {
		parsedLog := models.Log{
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
		parsedLogs = append(parsedLogs, parsedLog)
	}
	return parsedLogs
}

// parseTopics parses the topics and returns a slice of topic strings
func parseTopics(topics []common.Hash) []string {
	var parsedTopics []string
	for _, topic := range topics {
		parsedTopics = append(parsedTopics, topic.Hex())
	}
	return parsedTopics
}
