package grpc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"code.emcdtech.com/emcd/service/email/internal/berrors"
	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
	pb "code.emcdtech.com/emcd/service/email/protocol/email"
)

type paginationListSettings struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func paginationToBase64(pagination repository.Pagination) (string, error) {
	p := paginationListSettings{
		Page: pagination.Page,
		Size: pagination.Size,
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("pagination to json: %w", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func paginationFromBase64(s string) (repository.Pagination, error) {
	var (
		pagination repository.Pagination
		pg         paginationListSettings
	)
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return pagination, fmt.Errorf("base64 to json: %w", err)
	}
	err = json.Unmarshal(b, &pg)
	if err != nil {
		return pagination, fmt.Errorf("pagination from json: %w", err)
	}
	pagination.Page = pg.Page
	pagination.Size = pg.Size
	return pagination, nil
}

func convertError(err error) error {
	switch {
	case errors.Is(err, berrors.ErrTemplateNotFound):
		return errTemplateNotFound
	case errors.Is(err, berrors.ErrNoMailSenderAvailable):
		return errNoMailSenderAvailable
	case errors.Is(err, berrors.ErrProviderSettingNotFound):
		return errProviderSettingNotFound
	default:
		return err
	}
}

func convertReportTypeToModel(reportType pb.ReportType) (model.CodeTemplate, error) {
	switch reportType {
	case pb.ReportType_INCOME:
		return model.IncomeReport, nil
	case pb.ReportType_PAYOUT:
		return model.PayoutReport, nil
	case pb.ReportType_WORKER:
		return model.WorkerReport, nil
	default:
		return "", errors.New("invalid report type")
	}
}
