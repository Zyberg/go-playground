package types

import "math/big"

type TxRecord struct {
	BlockNumber uint64     `json:"block_number"`
	TxHash      string     `json:"tx_hash"`
	TxType      string     `json:"tx_type"`
	From        string     `json:"from"`
	To          string     `json:"to"`
	Value       *big.Float `json:"value"`
	Asset       string     `json:"asset"`
}
