package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	businessError "code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"code.emcdtech.com/b2b/endpoint/internal/repository"
	"code.emcdtech.com/b2b/endpoint/internal/repository/pg"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	ValidateRequest(next echo.HandlerFunc) echo.HandlerFunc
}

type auth struct {
	secretRepo repository.Secret
	ipRepo     repository.IP
	logRepo    repository.RequestLog
}

func NewAuth(
	secretRepo repository.Secret,
	ipRepo repository.IP,
	logRepo repository.RequestLog,
) Auth {
	return &auth{
		secretRepo: secretRepo,
		ipRepo:     ipRepo,
		logRepo:    logRepo,
	}
}

const clientIDKey = "client_id"

func (a *auth) ValidateRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKeyStr := c.Request().Header.Get("X-API-Key")
		timestampStr := c.Request().Header.Get("X-Timestamp")
		signature := c.Request().Header.Get("X-Signature")

		apiKey, err := uuid.Parse(apiKeyStr)
		if err != nil {
			log.Error(c.Request().Context(), "auth: invalid api key: %s, err: %s", apiKeyStr, err.Error())
			return c.JSON(http.StatusBadRequest, businessError.InvalidApiKey)
		}
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			log.Error(c.Request().Context(), "auth: invalid timestamp: %s, err: %s", timestampStr, err.Error())
			return c.JSON(http.StatusBadRequest, businessError.InvalidTimestamp)
		}

		secret, err := a.secretRepo.FindOne(c.Request().Context(), &model.SecretFilter{
			ApiKey: &apiKey,
		})
		if err != nil {
			log.Error(c.Request().Context(), "auth: secretRepo.FindOne: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, businessError.Internal)
		}

		if !secret.IsActive {
			log.Error(c.Request().Context(), "auth: api_key is not active: %s", secret.ApiKey)
			return c.JSON(http.StatusBadRequest, businessError.InvalidApiKey)
		}

		reqTime := timeFromMilliseconds(timestamp)

		// reqTime строго больше чем в предыдущий раз и меньше текущего времени
		if !reqTime.After(secret.LastUsed) || reqTime.After(time.Now().UTC()) {
			log.Error(c.Request().Context(), "auth: reqTime is invalid")
			return c.JSON(http.StatusBadRequest, businessError.InvalidTimestamp)
		}

		var params string

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Error(c.Request().Context(), "auth: read body: %s", err.Error())
			return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
		}
		if len(body) > 0 {
			// Восстанавливаем body, чтобы его можно было прочитать в хэндлере
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			params = string(body)
		}

		rawQuery := c.Request().URL.RawQuery
		if rawQuery != "" {
			params = fmt.Sprintf("%s%s", params, rawQuery)
		}

		expectedSignature := generateHMAC(secret.ApiSecret.String(), timestampStr, params)

		// Сравниваем с переданной подписью
		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			log.Error(c.Request().Context(), "auth: invalid signature")
			return c.JSON(http.StatusBadRequest, businessError.InvalidSignature)
		}

		isUnique, err := a.validateRequestUnique(c.Request().Context(), secret.ApiKey, signature)
		if err != nil {
			log.Error(c.Request().Context(), "auth: validateRequestUnique: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, businessError.Internal)
		}
		if !isUnique {
			log.Error(c.Request().Context(), "auth: request is not unique")
			return c.JSON(http.StatusBadRequest, businessError.DoubleRequest)
		}

		isValid, err := a.validateIP(c, secret.ApiKey)
		if err != nil {
			log.Error(c.Request().Context(), "auth: validateIP: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, businessError.Internal)
		}
		if !isValid {
			log.Error(c.Request().Context(), "auth: invalid ip: %s", getIPFromRequest(c))
			return c.JSON(http.StatusBadRequest, businessError.InvalidIP)
		}

		c.Set(clientIDKey, secret.ClientID)

		// Передаем управление хендлеру
		if err = next(c); err != nil {
			return err
		}

		// c.JSON не возвращает ошибку сюда. Поэтому логи пишутся и для запросов, которые упали с ошибкой.
		// Но, вроде, это хорошо?

		err = a.logRepo.Add(c.Request().Context(), &model.RequestLog{
			ApiKey:      apiKey,
			RequestHash: signature,
			CreatedAt:   reqTime,
		})
		if err != nil {
			log.Error(c.Request().Context(), "auth: logRepo.Add: %s", err.Error())
			// Думаю, не нужно возвращать error после успешноё обработки хендлером
			//return c.JSON(http.StatusInternalServerError, businessError.Internal)
		}

		err = a.secretRepo.Update(c.Request().Context(), secret,
			&model.SecretFilter{
				ApiKey: &secret.ApiKey,
			},
			&model.SecretPartial{
				LastUsed: &reqTime,
			})
		if err != nil {
			log.Error(c.Request().Context(), "auth: secretRepo.Update: %s", err.Error())
			//return c.JSON(http.StatusInternalServerError, businessError.Internal)
		}
		return nil
	}
}

func (a *auth) validateIP(c echo.Context, apiKey uuid.UUID) (bool, error) {
	ips, err := a.ipRepo.Find(c.Request().Context(), &model.IPFilter{
		ApiKey: &apiKey,
	})
	if err != nil {
		return false, fmt.Errorf("find: %w", err)
	}

	ip := getIPFromRequest(c)

	for i := range ips {
		if ip == ips[i].Address {
			return true, nil
		}
	}
	return false, nil
}

func (a *auth) validateRequestUnique(ctx context.Context, apiKey uuid.UUID, requestHash string) (bool, error) {
	_, err := a.logRepo.FindOne(ctx, &model.RequestLogFilter{
		ApiKey:      &apiKey,
		RequestHash: &requestHash,
	})
	if err != nil {
		if errors.Is(err, pg.LogsNotFoundError) {
			return true, nil
		}
		return false, fmt.Errorf("find: %w", err)
	}
	return false, nil

}

func generateHMAC(apiSecret, timestamp, params string) string {
	message := apiSecret + timestamp + params
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func getIPFromRequest(c echo.Context) string {
	ip := c.Request().Header.Get("X-Forwarded-For") // ip1,ip2,ip3,...
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 1 {
			ip = ips[0]
		}

		return ip
	}

	ip = c.Request().RemoteAddr
	ip = ip[:strings.LastIndex(ip, ":")]

	if ip == "[::1]" || ip == "127.0.0.1" {
		ip = "localhost"
	}

	return ip
}

func timeFromMilliseconds(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}
