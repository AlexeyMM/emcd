package controller

import (
	"errors"
	"net/http"

	"code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	_ "code.emcdtech.com/emcd/sdk/error"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SwapWithdraw struct {
	srv service.SwapWithdraw
}

func NewSwapWithdraw(srv service.SwapWithdraw) *SwapWithdraw {
	return &SwapWithdraw{
		srv: srv,
	}
}

type getTransactionLinkResponse struct {
	TransactionLink string `json:"transaction_link"`
}

// GetTransactionLink returns the transaction link in the explorer
// @Summary Get transaction link
// @Description Returns the transaction link for a given SwapID
// @Tags swapWithdraw
// @Produce json
// @Param swap_id path string true "Swap ID"
// @Success 200 {object} getTransactionLinkResponse "Transaction link successfully retrieved"
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/withdraw/tx_link/{swap_id} [get]
func (s *SwapWithdraw) GetTransactionLink(c echo.Context) error {
	swapIDParam := c.Param("swap_id")

	swapID, err := uuid.Parse(swapIDParam)
	if err != nil {
		log.Error(c.Request().Context(), "getTransactionLink: invalid swap_id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	transactionLink, err := s.srv.GetTransactionLink(c.Request().Context(), swapID)
	if err != nil {
		log.Error(c.Request().Context(), "getTransactionLink: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, getTransactionLinkResponse{
		TransactionLink: transactionLink,
	})
}
