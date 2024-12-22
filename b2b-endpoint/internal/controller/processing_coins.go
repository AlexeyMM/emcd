package controller

import (
	"errors"
	"net/http"

	businessError "code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/processing/protocol/coinpb"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProcessingCoinsController struct {
	processingCoinCli coinpb.CoinsServiceClient
}

func NewProcessingCoinsController(processingCoinCli coinpb.CoinsServiceClient) *ProcessingCoinsController {
	return &ProcessingCoinsController{
		processingCoinCli: processingCoinCli,
	}
}

type getCoinsResponseNetwork struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type getCoinsResponseCoin struct {
	ID          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	MediaURL    string                     `json:"media_url"`
	Networks    []*getCoinsResponseNetwork `json:"networks"`
}

type getCoinsResponse struct {
	Coins []*getCoinsResponseCoin `json:"coins"`
}

// GetCoins returns list of available coins with networks
// @Summary Get coins list
// @Description Get coins list with networks
// @Tags Processing
// @Produce json
// @Success 200 {object} getCoinsResponse
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /processing/coins [get]
func (p *ProcessingCoinsController) GetCoins(c echo.Context) error {
	resp, err := p.processingCoinCli.GetCoins(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		log.Error(c.Request().Context(), "getCoins: %s", err.Error())
		// TODO: use sdk/error in processing instead of custom proto message (though keep custom business error struct)
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}

		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	coins := make([]*getCoinsResponseCoin, 0, len(resp.Coins))
	for _, coin := range resp.Coins {
		networks := make([]*getCoinsResponseNetwork, 0, len(coin.Networks))
		for _, network := range coin.Networks {
			networks = append(networks, &getCoinsResponseNetwork{
				ID:    network.Id,
				Title: network.Title,
			})
		}
		coins = append(coins, &getCoinsResponseCoin{
			ID:          coin.Id,
			Title:       coin.Title,
			Description: coin.Description,
			MediaURL:    coin.MediaUrl,
			Networks:    networks,
		})
	}

	return c.JSON(http.StatusOK, getCoinsResponse{
		Coins: coins,
	})
}
