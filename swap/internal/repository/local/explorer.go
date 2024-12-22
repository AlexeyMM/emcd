package local

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
)

var (
	txBTC = "https://mempool.space/ru/tx/%s"
	txETH = "https://etherscan.io/tx/%s"
	txSOL = "https://solscan.io/tx/%s"
	txBNB = "https://bscscan.com/tx/%s"
	txXRP = "https://xrpscan.com/tx/%s"
	txTON = "https://tonscan.org/tx/%s"
	txINJ = "https://explorer.injective.network/transaction/%s"
)

type Explorer struct{}

func NewExplorer() *Explorer {
	return &Explorer{}
}

func (e *Explorer) GetTransactionLink(ctx context.Context, coin, hashID string) (string, error) {
	switch coin {
	case model.CoinBTC:
		return fmt.Sprintf(txBTC, hashID), nil
	case "ETH":
		return fmt.Sprintf(txETH, hashID), nil
	case "SOL":
		return fmt.Sprintf(txSOL, hashID), nil
	case "BNB":
		return fmt.Sprintf(txBNB, hashID), nil
	case "XRP":
		return fmt.Sprintf(txXRP, hashID), nil
	case "TON":
		return fmt.Sprintf(txTON, hashID), nil
	case "INJ":
		return fmt.Sprintf(txINJ, hashID), nil
	default:
		return "", nil
	}
}
