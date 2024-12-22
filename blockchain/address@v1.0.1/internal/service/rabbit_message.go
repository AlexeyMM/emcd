package service

const reloadAddressesCommand = "reload_addresses"

type AdminMessage struct {
	Command      string `json:"command"`
	TxID         string `json:"tx_id"`
	Hash         string `json:"hash"`
	ToBlockScore uint64 `json:"to_block_score"`
}
