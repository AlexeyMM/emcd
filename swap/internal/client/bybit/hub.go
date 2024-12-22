package bybit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/internal/slack"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"
)

const chanSize = 10

type registerEvent struct {
	ws     ws
	conn   *websocket.Conn
	cancel context.CancelFunc
}

type unregisterEvent struct {
	ws        ws
	accountID int64
}

type restartWalletEvent struct {
	ws          ws
	swapID      uuid.UUID
	swapCreated time.Time
	account     *model.Account
	coin        string
	amount      decimal.Decimal
}

type restartOrderEvent struct {
	ws                 ws
	account            *model.Account
	orders             []*model.Order
	receivedFirstOrder bool
}

type restartOrderbookEvent struct {
	ws      ws
	symbols []*model.Symbol
}

type client struct {
	conn   *websocket.Conn
	cancel context.CancelFunc
}

type Hub struct {
	mu                             sync.Mutex
	clients                        map[ws]*client
	registerCh                     chan *registerEvent
	unregisterCh                   chan *unregisterEvent
	restartWalletCh                chan *restartWalletEvent
	restartOrderCh                 chan *restartOrderEvent
	restartOrderBookCh             chan *restartOrderbookEvent
	orderSrv                       service.Order
	swapStatusUpdater              service.SwapStatusUpdater
	orderbookUpdateMessageCh       chan<- model.OrderBookUpdateMessage
	orderbookWsRecoveryCh          chan []*model.Symbol
	wsExpiredTime                  time.Duration
	masterKeys                     *model.Secrets
	depositWaitingPeriod           time.Duration
	slack                          slack.Slack
	byBitOrderBookWebsocketGauge   prometheus.Gauge
	byBitReconnectWebsocketCounter *prometheus.CounterVec
}

func NewHub(
	swapStatusUpdater service.SwapStatusUpdater,
	orderSrv service.Order,
	wsExpiredTime time.Duration,
	masterKeys *model.Secrets,
	orderbookUpdateMessageCh chan<- model.OrderBookUpdateMessage,
	orderbookWsRecoveryCh chan []*model.Symbol,
	depositWaitingPeriod time.Duration,
	slack slack.Slack,
	byBitOrderBookWebsocketGauge prometheus.Gauge,
	byBitReconnectWebsocketCounter *prometheus.CounterVec,
) *Hub {
	return &Hub{
		clients:                        make(map[ws]*client),
		registerCh:                     make(chan *registerEvent, chanSize),
		unregisterCh:                   make(chan *unregisterEvent, chanSize),
		restartWalletCh:                make(chan *restartWalletEvent, chanSize),
		restartOrderCh:                 make(chan *restartOrderEvent, chanSize),
		restartOrderBookCh:             make(chan *restartOrderbookEvent, chanSize),
		orderSrv:                       orderSrv,
		swapStatusUpdater:              swapStatusUpdater,
		orderbookUpdateMessageCh:       orderbookUpdateMessageCh,
		orderbookWsRecoveryCh:          orderbookWsRecoveryCh,
		wsExpiredTime:                  wsExpiredTime,
		masterKeys:                     masterKeys,
		depositWaitingPeriod:           depositWaitingPeriod,
		slack:                          slack,
		byBitOrderBookWebsocketGauge:   byBitOrderBookWebsocketGauge,
		byBitReconnectWebsocketCounter: byBitReconnectWebsocketCounter,
	}
}

func (h *Hub) Run(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return

		case <-t.C:
			log.Debug(ctx, "run clients: %d", len(h.clients))
			h.byBitOrderBookWebsocketGauge.Set(float64(len(h.clients)))

		case event := <-h.registerCh:
			if event == nil {
				log.Error(ctx, "register: received nil event")
				continue
			}
			h.register(event)

		case event := <-h.unregisterCh:
			if event == nil {
				log.Error(ctx, "unregister: received nil event")
				continue
			}
			log.Debug(ctx, "unregister: %+v", event)
			h.unregister(ctx, event)

		case event := <-h.restartWalletCh:
			if event == nil {
				log.Error(ctx, "restartWalletCh: received nil event")
				continue
			}

			h.unregister(ctx, &unregisterEvent{
				ws:        event.ws,
				accountID: event.account.ID,
			})

			err := h.SubscribeOnWallet(ctx, event.swapID, event.swapCreated, event.account, event.coin, event.amount)
			if err != nil {
				log.Error(ctx, "restart wallet fail, account_id: %d: %s", event.account.ID, err.Error())
				continue
			}

		case event := <-h.restartOrderCh:
			if event == nil {
				log.Error(ctx, "restartOrderCh: received nil event")
				continue
			}

			h.unregister(ctx, &unregisterEvent{
				ws:        event.ws,
				accountID: event.account.ID,
			})

			err := h.SubscribeOnOrders(ctx, event.account, event.orders, event.receivedFirstOrder)
			if err != nil {
				log.Error(ctx, "restart order fail, account_id: %d: %s", event.account.ID, err.Error())
				continue
			}

		case event := <-h.restartOrderBookCh:
			if event == nil {
				log.Error(ctx, "restartOrderBookCh: received nil event")
				continue
			}
			log.Info(ctx, "restartOrderBookCh: %+v", event)
			h.unregister(ctx, &unregisterEvent{
				ws:        event.ws,
				accountID: 0, // use only for log
			})

			err := h.SubscribeOnOrderbooks(ctx, event.symbols)
			if err != nil {
				log.Error(ctx, "restartOrderBookCh: subscribeOnOrderbooks: %s", err.Error())
			}
		}
	}
}

func (h *Hub) SubscribeOnWallet(ctx context.Context, swapID uuid.UUID, swapCreated time.Time, account *model.Account,
	coin string, amount decimal.Decimal) error {

	// TODO перед использованием включить проверку на полученную сумму. Если сумма депозита отличается от amount_from,
	// нужно пересчитать и обновить amount_from, amount_to, а так же проверить на лимиты.
	// Смотреть реализацию в service.Swap.CheckDeposit.updateAmountsAfterDeposit
	// реализовать можно в listen().handleMessage()
	if true {
		return fmt.Errorf("unimplemented")
	}

	conn, err := getPrivateWSConn(account.Keys, h.wsExpiredTime)
	if err != nil {
		return fmt.Errorf("getPrivateWSConn: %w", err)
	}

	w := newWalletWs(account, conn, h.unregisterCh, h.restartWalletCh, coin, amount, h.wsExpiredTime, h.swapStatusUpdater, swapID, swapCreated, h.depositWaitingPeriod, h.slack)

	newCtx, cancel := context.WithCancel(ctx)

	h.register(&registerEvent{
		ws:     w,
		conn:   conn,
		cancel: cancel,
	})

	err = subscribe(conn, []string{walletTopic})
	if err != nil {
		log.Debug(ctx, "subscribe on wallet err, unregister event sent: acc id: %d", account.ID)
		h.unregister(ctx, &unregisterEvent{
			ws:        w,
			accountID: w.account.ID,
		})
		return fmt.Errorf("subscribe: %w", err)
	}

	go w.listen(newCtx)

	return nil
}

func (h *Hub) SubscribeOnOrders(ctx context.Context, account *model.Account, orders []*model.Order, receivedFirstOrder bool) error {
	conn, err := getPrivateWSConn(account.Keys, h.wsExpiredTime)
	if err != nil {
		return fmt.Errorf("getPrivateWSConn: %w", err)
	}

	w := newOrderWs(account, conn, h.unregisterCh, h.restartOrderCh, orders, receivedFirstOrder, h.orderSrv, h.wsExpiredTime, h.byBitReconnectWebsocketCounter)

	newCtx, cancel := context.WithCancel(ctx)

	h.register(&registerEvent{
		ws:     w,
		conn:   conn,
		cancel: cancel,
	})

	err = subscribe(conn, []string{orderTopic})
	if err != nil {
		h.unregister(ctx, &unregisterEvent{
			ws:        w,
			accountID: w.account.ID,
		})
		return fmt.Errorf("subscribe: %w", err)
	}

	go w.listen(newCtx)

	return nil
}

func (h *Hub) SubscribeOnOrderbooks(ctx context.Context, symbols []*model.Symbol) error {
	var (
		conn *websocket.Conn
		err  error
	)
	defer func() {
		// Отправляем задачу на повторную подписку в случае провала.
		if err != nil {
			h.orderbookWsRecoveryCh <- symbols
			log.Debug(ctx, "hub: sent in orderbookWsRecovery channel: %+v", symbols)
		}
	}()

	conn, err = getPublicWSConn()
	if err != nil {
		return fmt.Errorf("getPublicWSConn: %w", err)
	}

	w := newOrderbookWs(conn, h.wsExpiredTime, symbols, h.restartOrderBookCh, h.orderbookUpdateMessageCh, h.byBitReconnectWebsocketCounter)

	newCtx, cancel := context.WithCancel(ctx)

	h.register(&registerEvent{
		ws:     w,
		conn:   conn,
		cancel: cancel,
	})

	topics := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		topics = append(topics, fmt.Sprintf(orderbookTopic, symbol.Title))
	}

	err = subscribe(conn, topics)
	if err != nil {
		h.unregister(ctx, &unregisterEvent{
			ws:        w,
			accountID: 0,
		})
		return fmt.Errorf("subscribe: %w", err)
	}

	go w.listen(newCtx)

	return nil
}

func (h *Hub) register(event *registerEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[event.ws] = &client{
		conn:   event.conn,
		cancel: event.cancel,
	}
}

func (h *Hub) unregister(ctx context.Context, event *unregisterEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	cli, ok := h.clients[event.ws]
	if !ok {
		log.Error(ctx, "unregister: client not found: account_id: %d", event.accountID)
		return
	}

	cli.cancel()

	err := cli.conn.Close()
	if err != nil {
		log.Error(ctx, "unregister: close conn: %s", err.Error())
	}

	delete(h.clients, event.ws)

	log.Debug(ctx, "unregister: client successfully removed, account_id: %d", event.accountID)
}

func (h *Hub) shutdown() {
	log.Debug(context.Background(), "hub shutdown")

	for w := range h.clients {
		h.unregister(context.Background(), &unregisterEvent{
			ws: w,
		})
	}
}
