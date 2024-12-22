package controller

import (
	"errors"
	"net/http"

	"code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/internal/controller/mapping"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/labstack/echo/v4"
)

type SwapCoin struct {
	srv service.SwapCoin
}

func NewSwapCoin(srv service.SwapCoin) *SwapCoin {
	return &SwapCoin{
		srv: srv,
	}
}

// GetAllV1 returns all coins and their networks
// Deprecated
// @Summary Get all coins
// @Description Retrieve all coins with their networks
// @Tags swapCoins
// @Produce json
// @Success 200 {object} map[string]response.Coin "List of coins"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/coins [get]
func (s *SwapCoin) GetAllV1(c echo.Context) error {
	coins, err := s.srv.GetAll(c.Request().Context())
	if err != nil {
		log.Error(c.Request().Context(), "getAll: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, mapping.MapToCoinsResponseDeprecated(coins))
}

// GetAllV2 returns all coins and their networks
// @Summary Get all coins
// @Description Retrieve all coins with their networks
// @Tags swapCoins
// @Produce json
// @Success 200 {object} []response.Coin "List of coins"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /v2/swap/coins [get]
func (s *SwapCoin) GetAllV2(c echo.Context) error {
	coins, err := s.srv.GetAll(c.Request().Context())
	if err != nil {
		log.Error(c.Request().Context(), "getAll: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, mapping.MapToCoinsResponse(coins))
}
