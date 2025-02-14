package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/internal/config"
	"github.com/hackertron/blocksight/internal/models"
	"github.com/hackertron/blocksight/internal/utils"
)

type AlchemyClient struct {
	client      *ethclient.Client
	rateLimiter *utils.RateLimiter
	retryConfig utils.RetryConfig
}

// NewAlchemyClient creates a new Alchemy client
func NewAlchemyClient(cfg *config.Config) (*AlchemyClient, error) {
	//apiKey := os.Getenv("ALCHEMY_API_KEY")
	//network := os.Getenv("ALCHEMY_NETWORK")

	//if apiKey == "" || network == "" {
	//	return nil, fmt.Errorf("ALCHEMY_API_KEY and ALCHEMY_NETWORK must be set")
	//}

	url := fmt.Sprintf("https://%s.g.alchemy.com/v2/%s", cfg.Alchemy.Network, cfg.Alchemy.APIKey)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Alchemy: %v", err)
	}

	return &AlchemyClient{
		client:      client,
		rateLimiter: utils.NewRateLimiter(cfg.Alchemy.RateLimit),
		retryConfig: utils.RetryConfig{
			MaxRetries: cfg.Alchemy.RetryCount,
			RetryDelay: cfg.Alchemy.RetryDelay,
		},
	}, nil
}

// GetLatestBlock retrieves the latest block
func (c *AlchemyClient) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	c.rateLimiter.Wait()
	var blockNumber uint64
	err := utils.WithRetry(ctx, func() error {
		var err error
		blockNumber, err = c.client.BlockNumber(ctx)
		return err
	}, c.retryConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %v", err)
	}

	return c.GetBlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

// GetBlockByNumber retrieves a specific block by number
func (c *AlchemyClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*models.Block, error) {
	c.rateLimiter.Wait()
	var block *models.Block

	err := utils.WithRetry(ctx, func() error {
		var ethBlock *types.Block
		var err error
		ethBlock, err = c.client.BlockByNumber(ctx, number)
		if err != nil {
			return fmt.Errorf("failed to get block by number %v: %v", number, err)
		}

		transactions := make([]models.Transaction, len(ethBlock.Transactions()))
		for i, tx := range ethBlock.Transactions() {
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

		block = &models.Block{
			Number:        ethBlock.Number(),
			Hash:          ethBlock.Hash(),
			ParentHash:    ethBlock.ParentHash(),
			Timestamp:     time.Unix(int64(ethBlock.Time()), 0),
			Transactions:  transactions,
			GasUsed:       new(big.Int).SetUint64(ethBlock.GasUsed()),
			GasLimit:      new(big.Int).SetUint64(ethBlock.GasLimit()),
			BaseFeePerGas: ethBlock.BaseFee(),
		}
		return nil
	}, c.retryConfig)

	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetTransactionReceipt retrieves a transaction receipt
func (c *AlchemyClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*models.TransactionReceipt, error) {
	c.rateLimiter.Wait()
	var receipt *models.TransactionReceipt

	err := utils.WithRetry(ctx, func() error {
		ethReceipt, err := c.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			return fmt.Errorf("failed to get transaction receipt for %v: %v", txHash, err)
		}

		logs := make([]models.Log, len(ethReceipt.Logs))
		for i, log := range ethReceipt.Logs {
			logs[i] = models.Log{
				Address:     log.Address,
				Topics:      log.Topics,
				Data:        log.Data,
				BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
				TxHash:      log.TxHash,
				LogIndex:    uint(log.Index),
			}
		}

		receipt = &models.TransactionReceipt{
			TransactionHash: ethReceipt.TxHash,
			BlockHash:       ethReceipt.BlockHash,
			BlockNumber:     ethReceipt.BlockNumber,
			GasUsed:         new(big.Int).SetUint64(ethReceipt.GasUsed),
			Status:          ethReceipt.Status,
			Logs:            logs,
		}
		return nil
	}, c.retryConfig)

	if err != nil {
		return nil, err
	}
	return receipt, nil
}
