package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"
)

type orderWs struct {
	account                        *model.Account
	conn                           *websocket.Conn
	unregister                     chan<- *unregisterEvent
	restart                        chan<- *restartOrderEvent
	orders                         []*model.Order
	receivedFirstOrder             bool
	orderSrv                       service.Order
	wsExpiredTime                  time.Duration
	byBitReconnectWebsocketCounter *prometheus.CounterVec
}

func newOrderWs(
	account *model.Account,
	conn *websocket.Conn,
	unregister chan<- *unregisterEvent,
	restart chan<- *restartOrderEvent,
	orders []*model.Order,
	receivedFirstOrder bool,
	orderSrv service.Order,
	wsExpiredTime time.Duration,
	ByBitReconnectWebsocketCounter *prometheus.CounterVec,
) *orderWs {
	return &orderWs{
		account:                        account,
		conn:                           conn,
		unregister:                     unregister,
		restart:                        restart,
		orders:                         orders,
		receivedFirstOrder:             receivedFirstOrder,
		orderSrv:                       orderSrv,
		wsExpiredTime:                  wsExpiredTime,
		byBitReconnectWebsocketCounter: ByBitReconnectWebsocketCounter,
	}
}

type orderData struct {
	Op    string  `json:"op"`
	Topic string  `json:"topic"`
	Data  []order `json:"data"`
}

type order struct {
	Symbol       string `json:"symbol"`
	OrderLinkID  string `json:"orderLinkId"`
	CumExecQty   string `json:"cumExecQty"`
	CumExecValue string `json:"cumExecValue"`
	CumExecFee   string `json:"cumExecFee"`
	OrderStatus  string `json:"orderStatus"`
}

func (w *orderWs) listen(ctx context.Context) {
	log.Debug(ctx, "start listening orders")
	defer log.Debug(ctx, "stop listening orders")

	go pingPong(ctx, w.conn)

	wsLifeTime := time.NewTimer(w.wsExpiredTime)
	defer wsLifeTime.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-wsLifeTime.C:
			log.Debug(ctx, "order ws expired")
			w.restart <- &restartOrderEvent{
				ws:                 w,
				account:            w.account,
				orders:             w.orders,
				receivedFirstOrder: w.receivedFirstOrder,
			}
			return

		default:
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					if err = ctx.Err(); err != nil {
						return
					}

					w.byBitReconnectWebsocketCounter.WithLabelValues("order_ws").Inc()

					w.restart <- &restartOrderEvent{
						ws:                 w,
						account:            w.account,
						orders:             w.orders,
						receivedFirstOrder: w.receivedFirstOrder,
					}
					return
				} else {
					log.Error(ctx, "websocket read: %s", err.Error())
					continue
				}
			}

			var data orderData
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Error(ctx, "websocket unmarshal 1: %s", err.Error())
				continue
			}

			if data.Op == "pong" {
				log.Debug(ctx, "pong")
				continue
			}

			if !strings.EqualFold(data.Topic, orderTopic) {
				continue
			}

			log.Debug(ctx, "order data: %+v\n", data)

			if len(w.orders) == 0 {
				log.Error(ctx, "listen order: len(orders) == 0")
				w.unregister <- &unregisterEvent{
					ws:        w,
					accountID: w.account.ID,
				}
				return
			}

			var (
				firstOrder  *model.Order
				secondOrder *model.Order
			)
			firstOrder = w.orders[0]
			if len(w.orders) == 2 {
				secondOrder = w.orders[1]
			}

			for _, c := range data.Data {
				if !w.receivedFirstOrder {
					amountFrom, amountTo, orderStatus, err := w.handleOrder(ctx, &c, firstOrder)
					if err != nil {
						log.Error(ctx, "websocket handle order 1: %s", err.Error())
						continue
					}

					err = w.orderSrv.Update(ctx, firstOrder,
						&model.OrderFilter{
							ID: &firstOrder.ID,
						},
						&model.OrderPartial{
							AmountFrom: &amountFrom,
							AmountTo:   &amountTo,
							Status:     &orderStatus,
						})
					if err != nil {
						log.Error(ctx, "listenOrder: orderRep.Update 1: order: %+v: %s", firstOrder, err.Error())
						w.unregister <- &unregisterEvent{
							ws:        w,
							accountID: w.account.ID,
						}
						return
					}
					log.Debug(ctx, "listenOrder: update order: %+v", firstOrder)
				} else {
					amountFrom, amountTo, orderStatus, err := w.handleOrder(ctx, &c, secondOrder)
					if err != nil {
						log.Error(ctx, "websocket handle order 2: %s", err.Error())
						continue
					}

					err = w.orderSrv.Update(ctx, secondOrder,
						&model.OrderFilter{
							ID: &secondOrder.ID,
						},
						&model.OrderPartial{
							AmountFrom: &amountFrom,
							AmountTo:   &amountTo,
							Status:     &orderStatus,
						})
					if err != nil {
						log.Error(ctx, "listenOrder: orderRep.Update 1: order: %+v: %s", secondOrder, err.Error())
						w.unregister <- &unregisterEvent{
							ws:        w,
							accountID: w.account.ID,
						}
						return
					}
					log.Debug(ctx, "listenOrder: update order: %+v", secondOrder)
				}

				if secondOrder == nil {
					log.Debug(ctx, "listenOrder: finished correctly")
					w.unregister <- &unregisterEvent{
						ws:        w,
						accountID: w.account.ID,
					}
					return
				}
				if w.receivedFirstOrder {
					log.Debug(ctx, "listenOrder: finished correctly")
					w.unregister <- &unregisterEvent{
						ws:        w,
						accountID: w.account.ID,
					}
					return
				}
				w.receivedFirstOrder = true
			}
		}
	}
}

func (w *orderWs) handleOrder(ctx context.Context, event *order, order *model.Order) (
	decimal.Decimal, decimal.Decimal, model.OrderStatus, error) {
	if event.Symbol != order.Symbol {
		log.Debug(ctx, "unexpected symbol: current order: %+v\n", order)
		return decimal.Zero, decimal.Zero, model.OrderUnknown, fmt.Errorf(
			"unexpected symbol for order received: account_id: %d, symbol: %s, current symbol: %s, order_id: %s",
			w.account.ID, event.Symbol, order.Symbol, event.OrderLinkID)
	}

	var (
		amountFrom, amountTo decimal.Decimal
		status               model.OrderStatus
	)

	switch event.OrderStatus {
	case "Rejected", "Cancelled":
		log.Error(ctx, "order failed, bybit status: %s", event.OrderStatus)
		return decimal.Zero, decimal.Zero, model.OrderFailed, nil

	case "PartiallyFilledCanceled", "Filled":
		log.Debug(ctx, "order filled: %+v", event)
		qty, val, fee, err := getQtyAndValAndFee(event.CumExecQty, event.CumExecValue, event.CumExecFee)
		if err != nil {
			return decimal.Zero, decimal.Zero, model.OrderUnknown, fmt.Errorf("getQtyAndVal: %w", err)
		}

		amountFrom, amountTo, err = getAmountFromAndAmountTo(qty, val, fee, order.Direction)
		if err != nil {
			return decimal.Zero, decimal.Zero, model.OrderUnknown, fmt.Errorf("getAmountFromAndAmountTo: %w", err)
		}

		if event.OrderStatus == "PartiallyFilledCanceled" {
			status = model.OrderPartiallyFilled
		} else if event.OrderStatus == "Filled" {
			status = model.OrderFilled
		}

	default:
		return decimal.Zero, decimal.Zero, model.OrderUnknown, fmt.Errorf(
			"order unexpected status: %s, status: %s", event.OrderLinkID, event.OrderStatus)
	}

	return amountFrom, amountTo, status, nil
}

func getQtyAndValAndFee(cumExecQty, cumExecValue, cumExecFee string) (decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	qty, err := decimal.NewFromString(cumExecQty)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("newFromString 1: %s", err.Error())

	}
	val, err := decimal.NewFromString(cumExecValue)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("newFromString 2: %s", err.Error())
	}
	fee, err := decimal.NewFromString(cumExecFee)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("newFromString 3: %s", err.Error())
	}
	return qty, val, fee, nil
}

func getAmountFromAndAmountTo(qty, val, fee decimal.Decimal, dir model.Direction) (decimal.Decimal, decimal.Decimal, error) {
	var amountFrom, amountTo decimal.Decimal

	if dir == model.Sell {
		amountFrom = qty
		amountTo = val.Sub(fee)
	} else {
		amountFrom = val
		amountTo = qty.Sub(fee)
	}

	return amountFrom, amountTo, nil
}
