package delegate

import (
	"context"
	"time"

	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

const (
	apiKeyHeader = "X-API-KEY"
)

type whitelabelEventClients struct {
	client whitelabel.WhitelabelServiceClient
}

func NewWhiteLabelEventClients(
	client whitelabel.WhitelabelServiceClient,
) *whitelabelEventClients {
	return &whitelabelEventClients{
		client: client,
	}
}

func (w *whitelabelEventClients) SendWLEvent(ctx context.Context, whiteLabelID uuid.UUID, req *model.WLEventRequest) {
	resp, err := w.client.GetByID(ctx, &whitelabel.GetByIDRequest{Id: whiteLabelID.String()})
	if err != nil {
		log.Error(ctx, "get white label from id: %s. %s", whiteLabelID, err.Error())
		return
	}
	wlClient := newWLClient(
		resp.GetWhiteLabel().GetUrl(),
		resp.GetWhiteLabel().GetApiKey(),
	)
	wlClient.sendWLEvent(ctx, req)
}

type wlClient struct {
	client *resty.Client
}

func (wm *wlClient) sendWLEvent(ctx context.Context, req *model.WLEventRequest) {
	res, err := wm.client.R().SetBody(req).Post("")
	if err != nil {
		log.Error(ctx, "whitelabelClients.sendWlEvent: %v", err)
		return
	}

	if res.IsError() {
		log.Error(ctx, "whitelabelClients.sendWlEvent: Response: code: %d, body: %v", res.StatusCode(), res.String())
	}
}

func newWLClient(url string, apiKey string) *wlClient {
	return &wlClient{
		client: resty.New().SetTimeout(300*time.Second).SetHeader(apiKeyHeader, apiKey).SetBaseURL(url),
	}
}
