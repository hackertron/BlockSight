package models

import "math/big"

// Block represents the structure for block data
type Block struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Nonce        uint64
	Transactions []Transaction
}

// Transaction represents the structure for transaction data
type Transaction struct {
	Hash  string
	Nonce uint64
	// Add more fields as needed
}
