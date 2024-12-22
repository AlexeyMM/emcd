package bybit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/gorilla/websocket"
)

type ws interface {
	listen(ctx context.Context)
}

const (
	pingPongInterval = 20 * time.Second

	walletTopic    = "wallet"
	orderTopic     = "order.spot"
	orderbookTopic = "orderbook.200.%s"
)

type subscribeMessage struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

func subscribe(conn *websocket.Conn, topics []string) error {
	m := subscribeMessage{
		Op:   "subscribe",
		Args: topics,
	}

	mJson, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshall: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, mJson)
	if err != nil {
		return fmt.Errorf("write: %s", err.Error())
	}

	return nil
}

type pingPongRequest struct {
	Op string `json:"op"`
}

func pingPong(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(pingPongInterval):
			pingJSON, err := json.Marshal(pingPongRequest{
				Op: "ping",
			})
			if err != nil {
				log.Error(ctx, "pingPong: marshal: %s", err.Error())
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, pingJSON)
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return
				}
				log.Error(ctx, "pingPong: write: %s", err.Error())
				continue
			}
		}
	}
}

type authRequest struct {
	ReqID string   `json:"req_id,omitempty"`
	Op    string   `json:"op"`
	Args  []string `json:"args"`
}

type authResponse struct {
	Success bool   `json:"success"`
	RetMsg  string `json:"ret_msg"`
	Op      string `json:"op"`
	ConnID  string `json:"conn_id"`
}

func getPrivateWSConn(keys *model.Secrets, wsExpiredTime time.Duration) (*websocket.Conn, error) {
	expires := time.Now().Add(wsExpiredTime).UnixMilli()

	signature := generateSignature(keys.ApiSecret, expires)

	u := url.URL{
		Scheme: "wss",
		Host:   "stream.bybit.com",
		Path:   "/v5/private"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	authMsg := authRequest{
		Op: "auth",
		Args: []string{
			keys.ApiKey,
			strconv.Itoa(int(expires)),
			signature,
		},
	}
	authJSON, err := json.Marshal(authMsg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, authJSON)
	if err != nil {
		return nil, fmt.Errorf("websocket write 3: %w", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	var response authResponse
	err = json.Unmarshal(message, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("auth failed: %s", response.RetMsg)
	}

	return conn, nil
}

func getPublicWSConn() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "stream.bybit.com",
		Path:   "/v5/public/spot",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return conn, nil
}

func generateSignature(apiSecret string, expires int64) string {
	message := fmt.Sprintf("GET/realtime%d", expires)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
