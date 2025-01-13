package models

// Transaction represents the structure for transaction data
type Transaction struct {
	Hash     string
	Nonce    uint64
	To       string
	Value    string
	Gas      uint64
	GasPrice string
	Input    []byte
	V        string
	R        string
	S        string
	IsPending bool
}
