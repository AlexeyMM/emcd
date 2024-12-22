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

type SwapUser struct {
	srv service.Swap
}

func NewSwapUser(srv service.Swap) *SwapUser {
	return &SwapUser{
		srv: srv,
	}
}

type createUserBySwapIDRequest struct {
	SwapID   uuid.UUID `json:"swap_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Language string    `json:"language"`
}

// CreateUserBySwapID creates a new user based on the provided SwapID
// @Summary Create user by SwapID
// @Description Create a new user using SwapID and email
// @Tags swapUsers
// @Accept json
// @Produce json
// @Param request body createUserBySwapIDRequest true "User creation request"
// @Success 201 "User successfully created"
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /swap/user [post]
func (s *SwapUser) CreateUserBySwapID(c echo.Context) error {
	var req createUserBySwapIDRequest

	if err := c.Bind(&req); err != nil {
		log.Error(c.Request().Context(), "createUserBySwapID: bind: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	if err := c.Validate(&req); err != nil {
		log.Error(c.Request().Context(), "createUserBySwapID: validate: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	//Поддержка только 2 языков, default - en
	if req.Language != "ru" && req.Language != "en" {
		req.Language = "en"
	}

	err := s.srv.StartSwap(c.Request().Context(), req.SwapID, req.Email, req.Language)
	if err != nil {
		log.Error(c.Request().Context(), "startSwap: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.NoContent(http.StatusCreated)
}
