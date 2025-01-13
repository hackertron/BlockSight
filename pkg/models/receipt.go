package models

// Receipt represents the structure for transaction receipt data
type Receipt struct {
	TxHash          string
	BlockHash       string
	BlockNumber     uint64
	TransactionIndex uint
	Status          uint64
	Logs            []Log
}
