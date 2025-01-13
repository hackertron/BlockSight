package block_parser

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/pkg/models"
)

// BlockParser represents the block parsing functionality
type BlockParser struct {
	client *ethclient.Client
}

// NewBlockParser creates a new instance of BlockParser
func NewBlockParser(client *ethclient.Client) *BlockParser {
	return &BlockParser{client: client}
}

// ParseBlock parses the block data and returns a Block model
func (bp *BlockParser) ParseBlock(blockNumber *big.Int) (*models.Block, error) {
	block, err := bp.client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}

	parsedBlock := &models.Block{
		Number:       block.Number().Uint64(),
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Nonce:        block.Nonce(),
		Transactions: parseTransactions(block.Transactions()),
	}

	return parsedBlock, nil
}

// parseTransactions parses the transactions and returns a slice of Transaction models
func parseTransactions(transactions types.Transactions) []models.Transaction {
	var parsedTransactions []models.Transaction
	for _, tx := range transactions {
		parsedTx := models.Transaction{
			Hash:  tx.Hash().Hex(),
			Nonce: tx.Nonce(),
			// Add more fields as needed
		}
		parsedTransactions = append(parsedTransactions, parsedTx)
	}
	return parsedTransactions
}
