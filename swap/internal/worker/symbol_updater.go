package worker

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
)

const (
	symbolUpdateInterval      = 2 * time.Hour
	ctxTimeoutForSymbolUpdate = time.Minute
)

type SymbolUpdater struct {
	symbolSrv        service.Symbol
	orderBookManager chan<- []*model.Symbol
}

func NewSymbolUpdater(symbolSrv service.Symbol, orderBookManager chan<- []*model.Symbol) *SymbolUpdater {
	return &SymbolUpdater{
		symbolSrv:        symbolSrv,
		orderBookManager: orderBookManager,
	}
}

func (s *SymbolUpdater) Run(ctx context.Context) error {
	t := time.NewTicker(symbolUpdateInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			log.Debug(ctx, "symbolUpdater: execute")

			newCtx, cancel := context.WithTimeout(ctx, ctxTimeoutForSymbolUpdate)

			// Получаем локальный список символов
			oldSymbols, err := s.symbolSrv.GetAll(ctx)
			if err != nil {
				log.Error(ctx, "symbolUpdater: getAll: %s", err.Error())
				cancel()
				continue
			}
			oldSymbolsMap := make(map[string]struct{}, len(oldSymbols))
			for _, oldSymbol := range oldSymbols {
				oldSymbolsMap[oldSymbol.Title] = struct{}{}
			}

			// Получаем актуальный список символов по API, обновляем базу
			actualSymbols, err := s.symbolSrv.SyncWithAPI(newCtx)
			if err != nil {
				log.Error(ctx, "symbolUpdater: uploadSymbols: %s", err.Error())
				cancel()
				continue
			}

			cancel()

			// Находим новые символы, которых нет локально, но есть в API
			var newSymbols []*model.Symbol
			for title, actualSymbol := range actualSymbols {
				if _, ok := oldSymbolsMap[title]; !ok {
					newSymbols = append(newSymbols, actualSymbol)
				}
			}

			// Отправляем новые символы ордербук менеджеру, для подписки на них
			s.orderBookManager <- newSymbols
		}
	}
}
