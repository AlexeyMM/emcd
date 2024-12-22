package controller

import (
	"errors"
	"net/http"

	"code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	secretSrv service.Secret
}

func NewAuth(secretSrv service.Secret) *Auth {
	return &Auth{
		secretSrv: secretSrv,
	}
}

type rotateAPIKeyResponse struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

// RotateKeys rotates the API keys using an existing key
// @Summary Rotate API Key
// @Description Generates a new API keys using the old one for authentication
// @Tags API Key
// @Produce json
// @Success 200 {object} rotateAPIKeyResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /auth/rotate-keys [post]
func (a *Auth) RotateKeys(c echo.Context) error {
	apiKeyStr := c.Request().Header.Get("X-API-Key")
	apiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	key, err := a.secretSrv.RotateKeys(c.Request().Context(), apiKey)
	if err != nil {
		log.Error(c.Request().Context(), "rotateKeys: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, rotateAPIKeyResponse{
		ApiKey:    key.ApiKey.String(),
		ApiSecret: key.ApiSecret.String(),
	})
}
