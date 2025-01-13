package transaction_processor

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/pkg/models"
)

// TransactionProcessor represents the transaction processing functionality
type TransactionProcessor struct {
	client *ethclient.Client
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(client *ethclient.Client) *TransactionProcessor {
	return &TransactionProcessor{client: client}
}

// ProcessTransaction processes the transaction data and returns a Transaction model
func (tp *TransactionProcessor) ProcessTransaction(txHash string) (*models.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, isPending, err := tp.client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	processedTransaction := &models.Transaction{
		Hash:     tx.Hash().Hex(),
		Nonce:    tx.Nonce(),
		To:       tx.To().Hex(),
		Value:    tx.Value().String(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().String(),
		Input:    tx.Data(),
		V:        tx.V().String(),
		R:        tx.R().String(),
		S:        tx.S().String(),
		IsPending: isPending,
	}

	return processedTransaction, nil
}
