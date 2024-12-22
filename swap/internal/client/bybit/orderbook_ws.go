package bybit

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
)

type orderbookWs struct {
	conn                           *websocket.Conn
	wsExpiredTime                  time.Duration
	symbols                        []*model.Symbol
	restart                        chan *restartOrderbookEvent
	orderbookUpdateMessageCh       chan<- model.OrderBookUpdateMessage
	byBitReconnectWebsocketCounter *prometheus.CounterVec
}

func newOrderbookWs(
	conn *websocket.Conn,
	wsExpiredTime time.Duration,
	symbols []*model.Symbol,
	restart chan *restartOrderbookEvent,
	orderbookUpdateMessageCh chan<- model.OrderBookUpdateMessage,
	byBitReconnectWebsocketCounter *prometheus.CounterVec,
) *orderbookWs {
	return &orderbookWs{
		conn:                           conn,
		wsExpiredTime:                  wsExpiredTime,
		symbols:                        symbols,
		restart:                        restart,
		orderbookUpdateMessageCh:       orderbookUpdateMessageCh,
		byBitReconnectWebsocketCounter: byBitReconnectWebsocketCounter,
	}
}

type orderbookData struct {
	Op    string    `json:"op"`
	Topic string    `json:"topic"`
	Type  string    `json:"type"`
	Ts    int64     `json:"ts"`
	Data  orderbook `json:"data"`
	Cts   int64     `json:"cts"`
}

type orderbook struct {
	S    string      `json:"s"`
	Bids [][2]string `json:"b"`
	Asks [][2]string `json:"a"`
	U    int64       `json:"u"`
	Seq  int64       `json:"seq"`
}

func (w *orderbookWs) listen(ctx context.Context) {
	log.Debug(ctx, "start listening orderbook")
	defer log.Debug(ctx, "stop listening orderbook")

	go pingPong(ctx, w.conn)

	wsLifeTime := time.NewTimer(w.wsExpiredTime)
	defer wsLifeTime.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-wsLifeTime.C:
			w.restart <- &restartOrderbookEvent{
				ws:      w,
				symbols: w.symbols,
			}
			return

		default:
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					if err = ctx.Err(); err != nil {
						return
					}

					w.byBitReconnectWebsocketCounter.WithLabelValues("orderbook_ws").Inc()

					w.restart <- &restartOrderbookEvent{
						ws:      w,
						symbols: w.symbols,
					}
					return
				} else {
					log.Error(ctx, "websocket read: %s", err.Error())
					continue
				}
			}

			var data orderbookData
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Error(ctx, "websocket unmarshal 1: %s", err.Error())
				continue
			}

			if data.Op == "pong" {
				log.Info(ctx, "pong orderbook: %s: %+v", string(message), w.symbols)
				continue
			}

			if !strings.Contains(data.Topic, "orderbook") {
				continue
			}

			var symbol string
			for i := range w.symbols {
				if strings.Contains(strings.ToUpper(data.Topic), w.symbols[i].Title) {
					symbol = w.symbols[i].Title
					break
				}
			}

			var isSnapshot bool
			if data.Type == "snapshot" {
				isSnapshot = true
			}

			w.orderbookUpdateMessageCh <- model.OrderBookUpdateMessage{
				Symbol:     symbol,
				Bids:       data.Data.Bids,
				Asks:       data.Data.Asks,
				IsSnapshot: isSnapshot,
			}
		}
	}
}
