package models

// Log represents the structure for log data
type Log struct {
	Address     string
	Topics      []string
	Data        []byte
	BlockNumber uint64
	TxHash      string
	TxIndex     uint
	BlockHash   string
	Index       uint
	Removed     bool
}
