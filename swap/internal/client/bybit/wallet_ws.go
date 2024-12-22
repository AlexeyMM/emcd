package bybit

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/internal/slack"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

type walletWs struct {
	account              *model.Account
	conn                 *websocket.Conn
	unregister           chan<- *unregisterEvent
	restart              chan<- *restartWalletEvent
	coin                 string
	expectedAmount       decimal.Decimal
	wsExpiredTime        time.Duration
	swapStatusUpdater    service.SwapStatusUpdater
	swapID               uuid.UUID
	swapCreated          time.Time
	depositWaitingPeriod time.Duration
	slack                slack.Slack
}

func newWalletWs(
	account *model.Account,
	conn *websocket.Conn,
	unregister chan<- *unregisterEvent,
	restart chan<- *restartWalletEvent,
	coin string,
	expectedAmount decimal.Decimal,
	wsExpiredTime time.Duration,
	swapStatusUpdater service.SwapStatusUpdater,
	swapID uuid.UUID,
	swapCreated time.Time,
	depositWaitingPeriod time.Duration,
	slack slack.Slack,
) *walletWs {
	return &walletWs{
		account:              account,
		conn:                 conn,
		unregister:           unregister,
		restart:              restart,
		coin:                 coin,
		expectedAmount:       expectedAmount,
		wsExpiredTime:        wsExpiredTime,
		swapStatusUpdater:    swapStatusUpdater,
		swapID:               swapID,
		swapCreated:          swapCreated,
		depositWaitingPeriod: depositWaitingPeriod,
		slack:                slack,
	}
}

type walletData struct {
	Op    string `json:"op"`
	Topic string `json:"topic"`
	Data  []struct {
		Coin []coin `json:"coin"`
	} `json:"data"`
}

type coin struct {
	Coin                string `json:"coin"`
	WalletBalance       string `json:"walletBalance"`
	Free                string `json:"free"`
	Locked              string `json:"locked"`
	AvailableToWithdraw string `json:"availableToWithdraw"`
	AvailableToBorrow   string `json:"availableToBorrow"`
	BorrowAmount        string `json:"borrowAmount"`
}

func (w *walletWs) listen(ctx context.Context) {
	log.Debug(ctx, "start listening wallet: %s", w.swapID.String())
	defer log.Debug(ctx, "stop listening wallet: %s", w.swapID.String())

	go pingPong(ctx, w.conn)

	wsLifeTime := time.NewTimer(w.wsExpiredTime)
	defer wsLifeTime.Stop()

	remainTime := w.swapCreated.Add(w.depositWaitingPeriod).Sub(time.Now().UTC())
	if remainTime <= 0 {
		log.Debug(ctx, "walletWs listen: remainTime is less than 0: swapID: %s", w.swapID.String())
		return
	}
	depositWaitingPeriod := time.NewTimer(remainTime)
	defer depositWaitingPeriod.Stop()

	// Когда время ожидания депозита закончится (<-maxDepositWaitingPeriod.C)
	// conn будет закрыт и handleMessage goroutine остановлена
	// так же будет закрыт context и остановлена listen goroutine
	// handleMessage goroutine нужна, что бы отрабатывал select case, без вечного ожидания в default
	go w.handleMessage(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-wsLifeTime.C:
			log.Debug(ctx, "walletWs expired: %s", w.swapID.String())
			w.restart <- &restartWalletEvent{
				ws:          w,
				swapID:      w.swapID,
				swapCreated: w.swapCreated,
				account:     w.account,
				coin:        w.coin,
				amount:      w.expectedAmount,
			}
			return

		case <-depositWaitingPeriod.C:
			log.Debug(ctx, "walletWs: <-maxDepositWaitingPeriod.C: swap_id: %s", w.swapID.String())

			sw, err := w.swapStatusUpdater.GetSwap(ctx, w.swapID)
			if err != nil {
				log.Error(ctx, "walletWs: getSwap: swapID: %s, err 1: %s", w.swapID.String(), err.Error())
				w.unregister <- &unregisterEvent{
					ws:        w,
					accountID: w.account.ID,
				}
				return
			}

			err = w.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Cancel)
			if err != nil {
				log.Error(ctx, "walletWs: updateAndBroadcast: swapID: %s, err 1: %s", w.swapID.String(), err.Error())
				w.unregister <- &unregisterEvent{
					ws:        w,
					accountID: w.account.ID,
				}
				return
			}

			w.unregister <- &unregisterEvent{
				ws:        w,
				accountID: w.account.ID,
			}
			return

		}
	}
}

func (w *walletWs) handleMessage(ctx context.Context) {
	log.Debug(ctx, "start handleMessage, swap_id: %d", w.swapID.String())
	defer log.Debug(ctx, "stop handleMessage, swap_id: %d", w.swapID.String())

	for {
		select {
		case <-ctx.Done():
			return

		default:
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					if err = ctx.Err(); err != nil {
						return
					}

					w.restart <- &restartWalletEvent{
						ws:          w,
						swapID:      w.swapID,
						swapCreated: w.swapCreated,
						account:     w.account,
						coin:        w.coin,
						amount:      w.expectedAmount,
					}
					return
				} else {
					log.Error(ctx, "walletWs: websocket read: %s", err.Error())
					continue
				}
			}

			log.Debug(ctx, "walletWs: message: %s", string(message))

			var data walletData
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Error(ctx, "walletWs: websocket unmarshal 1: %s", err.Error())
				continue
			}

			if data.Op == "pong" {
				log.Info(ctx, "pong", string(message))
				continue
			}

			if !strings.EqualFold(data.Topic, walletTopic) {
				continue
			}

			log.Debug(ctx, "walletWs: wallet data: swap_id: %s: %+v\n", w.swapID.String(), data)
			for _, dataStruct := range data.Data {
				for _, c := range dataStruct.Coin {
					if c.Coin != w.coin {
						log.Error(ctx, "walletWs: websocket coin mismatch: %s, expected %s", c.Coin, w.coin)
						continue
					}

					// TODO проверка, amount_from и фактическая полученная сумма

					log.Info(ctx, "walletWs: websocket received deposit: accountID: %d, coin: %s", w.account.ID, c.Coin)

					sw, err := w.swapStatusUpdater.GetSwap(ctx, w.swapID)
					if err != nil {
						log.Error(ctx, "walletWs: getSwap: swapID: %s, err 2: %s", w.swapID.String(), err.Error())
						return
					}

					err = w.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.CheckDeposit)
					if err != nil {
						log.Error(ctx, "walletWs: updateAndBroadcast: swapID: %s, err 2: %s", w.swapID.String(), err.Error())
						w.unregister <- &unregisterEvent{
							ws:        w,
							accountID: w.account.ID,
						}
						return
					}

					log.Debug(ctx, "walletWs: send unregister event: swap_id: %s", w.swapID.String())
					w.unregister <- &unregisterEvent{
						ws:        w,
						accountID: w.account.ID,
					}

					return
				}
			}
		}
	}
}
