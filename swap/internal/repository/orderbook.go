package repository

import (
	"code.emcdtech.com/b2b/swap/model"
)

//go:generate mockery --name=OrderBook
type OrderBook interface {
	Init([]*model.Symbol) error
	AddSnapshot(symbol string, bids, asks [][2]string) error
	AddDelta(symbol string, bids, asks [][2]string) error
	GetBidTopLevels(symbol string, n int) ([][2]float64, error)
	GetAskTopLevels(symbol string, n int) ([][2]float64, error)
	IsExist(symbol string) bool
	Len() int
}
