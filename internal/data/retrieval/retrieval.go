package retrieval

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Retrieval represents the data retrieval functionality
type Retrieval struct {
	client *ethclient.Client
}

// NewRetrieval creates a new instance of Retrieval
func NewRetrieval(alchemyAPIKey string) (*Retrieval, error) {
	client, err := ethclient.Dial(fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", alchemyAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum network: %v", err)
	}

	return &Retrieval{client: client}, nil
}

// Close closes the Ethereum client connection
func (r *Retrieval) Close() {
	r.client.Close()
}

// RetrieveBlockByNumber retrieves a block by its block number
func (r *Retrieval) RetrieveBlockByNumber(blockNumber *big.Int) (*types.Block, error) {
	block, err := r.client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}

	return block, nil
}

// RetrieveTransactionReceipt retrieves the transaction receipt for a given transaction hash
func (r *Retrieval) RetrieveTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	receipt, err := r.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction receipt: %v", err)
	}

	return receipt, nil
}
