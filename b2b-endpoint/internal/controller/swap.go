package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type Swap struct {
	srv service.Swap
}

func NewSwap(srv service.Swap) *Swap {
	return &Swap{
		srv: srv,
	}
}

type estimateResponse struct {
	AmountFrom string `json:"amount_from"`
	AmountTo   string `json:"amount_to"`
	MinFrom    string `json:"min_from"`
	MaxFrom    string `json:"max_from"`
	Rate       string `json:"rate"`
}

// Estimate calculates the estimated amount and limits for a swap
// @Summary Estimate the amount and limits for a swap
// @Description Estimate the swap amount and limits based on input query parameters
// @Tags swap
// @Produce json
// @Param coin_from query string true "Coin to swap from"
// @Param coin_to query string true "Coin to swap to"
// @Param network_from query string true "Network for coin_from"
// @Param network_to query string true "Network for coin_to"
// @Param amount_from query string false "Amount of coin_from to swap"
// @Param amount_to query string false "Amount of coin_to to swap"
// @Success 200 {object} estimateResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/estimate [get]
func (s *Swap) Estimate(c echo.Context) error {
	coinFrom := c.QueryParam("coin_from")
	coinTo := c.QueryParam("coin_to")
	networkFrom := c.QueryParam("network_from")
	networkTo := c.QueryParam("network_to")
	amountFromStr := c.QueryParam("amount_from")
	amountToStr := c.QueryParam("amount_to")

	if coinFrom == "" || coinTo == "" || networkFrom == "" || networkTo == "" {
		log.Error(c.Request().Context(), "estimate: invalid request")
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	if amountFromStr == "" && amountToStr == "" {
		amountFromStr = "1"
	}

	var (
		amountFrom, amountTo decimal.Decimal
		err                  error
	)

	if amountFromStr != "" {
		amountFrom, err = decimal.NewFromString(amountFromStr)
		if err != nil {
			log.Error(c.Request().Context(), "estimate: invalid amount_from")
			return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
		}
	}

	if amountToStr != "" {
		amountTo, err = decimal.NewFromString(amountToStr)
		if err != nil {
			log.Error(c.Request().Context(), "estimate: invalid amount_to")
			return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
		}
	}

	amountFrom, amountTo, rate, limits, err := s.srv.Estimate(c.Request().Context(), &model.EstimateRequest{
		CoinFrom:    coinFrom,
		CoinTo:      coinTo,
		NetworkFrom: networkFrom,
		NetworkTo:   networkTo,
		AmountFrom:  amountFrom,
		AmountTo:    amountTo,
	})
	if err != nil {
		log.Error(c.Request().Context(), "estimate: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, estimateResponse{
		AmountFrom: amountFrom.String(),
		AmountTo:   amountTo.String(),
		MinFrom:    limits.Min.String(),
		MaxFrom:    limits.Max.String(),
		Rate:       rate.String(),
	})
}

type swapRequest struct {
	CoinFrom    string `json:"coin_from" validate:"required"`
	CoinTo      string `json:"coin_to" validate:"required"`
	NetworkFrom string `json:"network_from" validate:"required"`
	NetworkTo   string `json:"network_to" validate:"required"`
	AmountFrom  string `json:"amount_from"`
	AmountTo    string `json:"amount_to"`
	AddressTo   string `json:"address_to"`
	TagTo       string `json:"tag_to"`
}

type swapResponse struct {
	CoinFrom     string `json:"coin_from" validate:"required"`
	CoinTo       string `json:"coin_to" validate:"required"`
	NetworkFrom  string `json:"network_from" validate:"required"`
	NetworkTo    string `json:"network_to" validate:"required"`
	AmountFrom   string `json:"amount_from"`
	AmountTo     string `json:"amount_to"`
	AddressFrom  string `json:"address_from"`
	AddressTo    string `json:"address_to"`
	TagFrom      string `json:"tag_from"`
	TagTo        string `json:"tag_to"`
	Status       int32  `json:"status"`
	StartTime    int    `json:"start_time"`    // unix
	SwapDuration int    `json:"swap_duration"` // секунд
	Rate         string `json:"rate"`
}

type depositResponse struct {
	ID             string `json:"id"`
	DepositAddress string `json:"deposit_address"`
	DepositTag     string `json:"deposit_tag,omitempty"`
}

// Swap executes the swap and returns deposit address information
// @Summary Execute a swap and get deposit address
// @Description Perform a swap operation and return the deposit address data
// @Tags swap
// @Accept json
// @Produce json
// @Param request body swapRequest true "Swap request"
// @Success 200 {object} depositResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap [post]
func (s *Swap) Swap(c echo.Context) error {
	var req swapRequest
	if err := c.Bind(&req); err != nil {
		log.Error(c.Request().Context(), "swap: bind: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	partnerID := c.Request().Header.Get("X-Partner-ID")

	var (
		amountFrom, amountTo decimal.Decimal
		err                  error
	)

	if req.AmountFrom != "" {
		amountFrom, err = decimal.NewFromString(req.AmountFrom)
		if err != nil {
			log.Error(c.Request().Context(), "swap: newFromString 1: %s", err.Error())
			return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
		}
	}
	if req.AmountTo != "" {
		amountTo, err = decimal.NewFromString(req.AmountTo)
		if err != nil {
			log.Error(c.Request().Context(), "swap: newToString 2: %s", err.Error())
			return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
		}
	}

	id, depositAddress, err := s.srv.PrepareSwap(c.Request().Context(), &model.SwapRequest{
		CoinFrom:    req.CoinFrom,
		CoinTo:      req.CoinTo,
		NetworkFrom: req.NetworkFrom,
		NetworkTo:   req.NetworkTo,
		AmountFrom:  amountFrom,
		AmountTo:    amountTo,
		AddressTo: &model.AddressData{
			Address: req.AddressTo,
			Tag:     req.TagTo,
		},
		PartnerID: partnerID,
	})
	if err != nil {
		log.Error(c.Request().Context(), "swap: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, depositResponse{
		ID:             id.String(),
		DepositAddress: depositAddress.Address,
		DepositTag:     depositAddress.Tag,
	})
}

// Status returns swap status updates via SSE
// @Summary Get swap status updates
// @Description Retrieve real-time swap status updates using Server-Sent Events (SSE)
// @Tags swap
// @Produce text/event-stream
// @Param swapID path string true "Swap ID"
// @Success 200 {string} string "Swap status updates"
// @Failure 400 {object} sdkError.Error "Invalid swap ID"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/status/{swapID} [get]
func (s *Swap) Status(c echo.Context) error {
	swapIDRaw := c.Param("swapID")

	swapID, err := uuid.Parse(swapIDRaw)
	if err != nil {
		log.Error(c.Request().Context(), "swap: invalid swapID: %s", swapIDRaw)
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	ch := make(chan model.Status)

	err = s.srv.Status(c.Request().Context(), swapID, ch)
	if err != nil {
		log.Error(c.Request().Context(), "status: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	var status int
	pingInterval := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-c.Request().Context().Done():
			return c.Request().Context().Err()
		case st, ok := <-ch:
			if !ok {
				return nil
			}
			status = int(st)
			err = sendStatus(c.Response(), status)
			if err != nil {
				log.Error(c.Request().Context(), "sendStatus 1: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, businessError.Internal)
			}

		case <-pingInterval.C:
			err = sendStatus(c.Response(), status)
			if err != nil {
				log.Error(c.Request().Context(), "sendStatus 2: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, businessError.Internal)
			}
		}
	}
}

func sendStatus(response *echo.Response, status int) error {
	_, err := fmt.Fprintf(response, "data: %d\n\n", status)
	if err != nil {
		return fmt.Errorf("fmt.Fprintf: %w", err)
	}
	response.Flush()
	return nil
}

// GetSwapByID returns a swap by transaction ID
// @Summary Get swap by transaction ID
// @Description Get a swap by transaction ID
// @Tags swap
// @Produce json
// @Success 200 {object} swapResponse
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/{swapID} [get]
func (s *Swap) GetSwapByID(c echo.Context) error {
	swp, err := s.srv.GetSwapByID(c.Request().Context(), c.Param("swapID"))
	if err != nil {
		log.Error(c.Request().Context(), "GetSwapByTxID: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, swapResponse{
		CoinFrom:     swp.CoinFrom,
		CoinTo:       swp.CoinTo,
		NetworkFrom:  swp.NetworkFrom,
		NetworkTo:    swp.NetworkTo,
		AmountFrom:   swp.AmountFrom.String(),
		AmountTo:     swp.AmountTo.String(),
		AddressFrom:  swp.AddressFrom.Address,
		AddressTo:    swp.AddressTo.Address,
		TagFrom:      swp.AddressFrom.Tag,
		TagTo:        swp.AddressTo.Tag,
		Status:       int32(swp.Status),
		StartTime:    int(swp.StartTime.Unix()),
		SwapDuration: int(swp.SwapDuration.Seconds()),
		Rate:         swp.Rate.String(),
	})
}

type supportRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Text  string `json:"text" validate:"required,lte=1000"`
}

// SendSupportMessage processes a request to send a message to support
// @Summary Send a support message
// @Description Receive username, email, and message text and send it to support team
// @Tags swapSupport
// @Accept json
// @Produce json
// @Param request body supportRequest true "Support request"
// @Success 200
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/support/message [post]
func (s *Swap) SendSupportMessage(c echo.Context) error {
	req := new(supportRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	err := s.srv.SendSupportMessage(c.Request().Context(), req.Name, req.Email, req.Text)
	if err != nil {
		log.Error(c.Request().Context(), "sendSupportMessage: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.NoContent(http.StatusOK)
}

type swapMessageRequest struct {
	Email  string `json:"email" validate:"required,email"`
	SwapID string `json:"swapID" validate:"required,uuid"`
}

type swapInfoEmailResponse struct {
	Ok     bool   `json:"ok"`
	Email  string `json:"email"`
	SwapID string `json:"swapID"`
}

// SendSwapInfoEmail Send email about swap to user email
// @Summary Send email about swap to user email
// @Description Receive email and send it to user email
// @Tags swap
// @Accept json
// @Produce json
// @Param request body swapInfoEmailResponse true "Swap email message request"
// @Success 200 {object} swapInfoEmailResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/email [post]
func (s *Swap) SendSwapInfoEmail(c echo.Context) error {
	req := new(swapMessageRequest)

	return c.JSON(http.StatusOK, swapInfoEmailResponse{
		Ok:     true,
		Email:  req.Email,  // TODO: get email from request
		SwapID: req.SwapID, // TODO: get swapID from request
	})
}
