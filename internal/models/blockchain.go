package models

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Block represents an Ethereum block
type Block struct {
	Number        *big.Int      `json:"number"`
	Hash          common.Hash   `json:"hash"`
	ParentHash    common.Hash   `json:"parentHash"`
	Timestamp     time.Time     `json:"timestamp"`
	Transactions  []Transaction `json:"transactions"`
	GasUsed       *big.Int      `json:"gasUsed"`
	GasLimit      *big.Int      `json:"gasLimit"`
	BaseFeePerGas *big.Int      `json:"baseFeePerGas"`
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash     common.Hash     `json:"hash"`
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Value    *big.Int        `json:"value"`
	GasPrice *big.Int        `json:"gasPrice"`
	Gas      uint64          `json:"gas"`
	Input    []byte          `json:"input"`
	Nonce    uint64          `json:"nonce"`
}

// TransactionReceipt represents an Ethereum transaction receipt
type TransactionReceipt struct {
	TransactionHash common.Hash `json:"transactionHash"`
	BlockHash       common.Hash `json:"blockHash"`
	BlockNumber     *big.Int    `json:"blockNumber"`
	GasUsed         *big.Int    `json:"gasUsed"`
	Status          uint64      `json:"status"`
	Logs            []Log       `json:"logs"`
}

// Log represents an Ethereum event log
type Log struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        []byte         `json:"data"`
	BlockNumber *big.Int       `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	LogIndex    uint           `json:"logIndex"`
}
