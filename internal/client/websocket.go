package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hackertron/blocksight/internal/models"
)

type WebSocketClient struct {
	client  *ethclient.Client
	blockCh chan *models.Block
}

func NewWebSocketClient(wsURL string) (*WebSocketClient, error) {
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, fmt.Errorf("websocket connection failed: %v", err)
	}

	return &WebSocketClient{
		client:  client,
		blockCh: make(chan *models.Block),
	}, nil
}

func (w *WebSocketClient) SubscribeToNewBlocks(ctx context.Context) (<-chan *models.Block, error) {
	headers := make(chan *types.Header)
	sub, err := w.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, fmt.Errorf("subscription failed: %v", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				block, err := w.client.BlockByHash(ctx, header.Hash())
				if err != nil {
					log.Printf("Failed to get block: %v", err)
					continue
				}
				w.blockCh <- convertToModelBlock(block)
			}
		}
	}()

	return w.blockCh, nil
}

func convertToModelBlock(block *types.Block) *models.Block {
	transactions := make([]models.Transaction, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		transactions[i] = models.Transaction{
			Hash:     tx.Hash(),
			From:     common.Address{}, // Will need receipt to get this
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
	}
}
