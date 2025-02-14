package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TraceResult represents the result of a transaction trace
// Note: This is a placeholder for future implementation with a Geth node
type TraceResult struct {
	From    common.Address `json:"from"`
	To      common.Address `json:"to"`
	Value   *big.Int       `json:"value"`
	GasUsed uint64         `json:"gasUsed"`
	Error   string         `json:"error,omitempty"`
}

// TODO: Implement tracing functionality when using a Geth node
// Requirements for tracing:
// 1. Connect to a Geth node with --http.api="eth,debug" flag
// 2. Use custom client implementation that supports debug namespace
// 3. Implement proper error handling for trace requests
