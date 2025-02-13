package client

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/internal/models"
)

type AlchemyClient struct {
	client *ethclient.Client
}

// NewAlchemyClient creates a new Alchemy client
func NewAlchemyClient() (*AlchemyClient, error) {
	apiKey := os.Getenv("ALCHEMY_API_KEY")
	network := os.Getenv("ALCHEMY_NETWORK")

	if apiKey == "" || network == "" {
		return nil, fmt.Errorf("ALCHEMY_API_KEY and ALCHEMY_NETWORK must be set")
	}

	url := fmt.Sprintf("https://%s.g.alchemy.com/v2/%s", network, apiKey)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Alchemy: %v", err)
	}

	return &AlchemyClient{
		client: client,
	}, nil
}

// GetLatestBlock retrieves the latest block
func (c *AlchemyClient) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %v", err)
	}

	return c.GetBlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

// GetBlockByNumber retrieves a specific block by number
func (c *AlchemyClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*models.Block, error) {
	block, err := c.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number %v: %v", number, err)
	}

	transactions := make([]models.Transaction, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		transactions[i] = models.Transaction{
			Hash:     tx.Hash(),
			From:     common.Address{}, // Need to get from receipt
			To:       tx.To(),
			Value:    tx.Value(),
			GasPrice: tx.GasPrice(),
			Gas:      tx.Gas(),
			Input:    tx.Data(),
			Nonce:    tx.Nonce(),
		}
	}

	return &models.Block{
		Number:        block.Number(),
		Hash:          block.Hash(),
		ParentHash:    block.ParentHash(),
		Timestamp:     time.Unix(int64(block.Time()), 0),
		Transactions:  transactions,
		GasUsed:       new(big.Int).SetUint64(block.GasUsed()),
		GasLimit:      new(big.Int).SetUint64(block.GasLimit()),
		BaseFeePerGas: block.BaseFee(),
	}, nil
}

// GetTransactionReceipt retrieves a transaction receipt
func (c *AlchemyClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*models.TransactionReceipt, error) {
	receipt, err := c.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt for %v: %v", txHash, err)
	}

	logs := make([]models.Log, len(receipt.Logs))
	for i, log := range receipt.Logs {
		logs[i] = models.Log{
			Address:     log.Address,
			Topics:      log.Topics,
			Data:        log.Data,
			BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
			TxHash:      log.TxHash,
			LogIndex:    uint(log.Index),
		}
	}

	return &models.TransactionReceipt{
		TransactionHash: receipt.TxHash,
		BlockHash:       receipt.BlockHash,
		BlockNumber:     receipt.BlockNumber,
		GasUsed:         new(big.Int).SetUint64(receipt.GasUsed),
		Status:          receipt.Status,
		Logs:            logs,
	}, nil
}
