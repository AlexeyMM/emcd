package worker

import (
	"context"
	"sync/atomic"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
)

const (
	OrderbookChanSize = 10000

	workerPoolMinSize           = 3
	workerPoolMaxSize           = 120
	workerPoolStartSize         = 10
	workerAlivePeriodWithoutJob = 10 * time.Second
	workersCheckPeriod          = 10 * time.Second
	sizeUpdateChanToStartWorker = 10
	numberWorkersToStartATime   = 5
)

type OrderBookManager struct {
	srv         service.OrderBook
	updateCh    <-chan model.OrderBookUpdateMessage // Изменения в orderbook
	newSymbols  <-chan []*model.Symbol              // Новые символы, на которые нужно подписаться
	busyWorkers int32
	allWorkers  int32
}

func NewOrderBookManager(srv service.OrderBook, updateCh <-chan model.OrderBookUpdateMessage, newSymbols <-chan []*model.Symbol) *OrderBookManager {
	return &OrderBookManager{
		srv:        srv,
		updateCh:   updateCh,
		newSymbols: newSymbols,
	}
}

func (m *OrderBookManager) Run(ctx context.Context) error {
	for i := 0; i < workerPoolStartSize; i++ {
		atomic.AddInt32(&m.allWorkers, 1)
		go m.listenUpdates(ctx)
	}
	go m.listenNewSymbols(ctx)
	go m.workersManagement(ctx)
	return nil
}

func (m *OrderBookManager) listenUpdates(ctx context.Context) {
	t := time.NewTimer(workerAlivePeriodWithoutJob)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if m.allWorkers > workerPoolMinSize {
				atomic.AddInt32(&m.allWorkers, -1)
				return
			} else {
				t.Reset(workerAlivePeriodWithoutJob)
			}
		case update := <-m.updateCh:
			atomic.AddInt32(&m.busyWorkers, 1)

			err := m.srv.Update(ctx, update.Symbol, update.Bids, update.Asks, update.IsSnapshot)
			if err != nil {
				log.Error(ctx, "orderbook manager: update: %s", err.Error())
				atomic.AddInt32(&m.busyWorkers, -1)
				continue
			}

			if !t.Stop() {
				<-t.C
			}
			t.Reset(workerAlivePeriodWithoutJob)
			atomic.AddInt32(&m.busyWorkers, -1)
		}
	}
}

func (m *OrderBookManager) listenNewSymbols(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newSymbols := <-m.newSymbols:
			log.Debug(ctx, "orderbook manager: new symbols: %d", len(newSymbols))

			err := m.srv.Subscribe(ctx, newSymbols)
			if err != nil {
				log.Error(ctx, "orderbook manager: subscribe: %s", err.Error())
				continue
			}
		}
	}
}

func (m *OrderBookManager) workersManagement(ctx context.Context) {
	t := time.NewTicker(workersCheckPeriod)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			//fmt.Printf("size chan: %d, busyWorkers: %d, allWorkers: %d \n", len(m.updateCh), m.busyWorkers, m.allWorkers)
			if len(m.updateCh) > sizeUpdateChanToStartWorker {

				if m.allWorkers < workerPoolMaxSize {
					for i := 0; i < numberWorkersToStartATime; i++ {
						atomic.AddInt32(&m.allWorkers, 1)
						go m.listenUpdates(ctx)
					}
				} else {
					log.Warn(ctx, "the limit of workers has been reached: %d", m.allWorkers)
				}
			}
		}
	}
}
