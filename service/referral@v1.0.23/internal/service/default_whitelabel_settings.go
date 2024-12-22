// Package service defines a set of functions for performing operations related to a service.
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

// DefaultWhitelabelSettings is an interface that defines methods for managing default whitelabel settings.
// It provides functionality to create, update, delete, get, and retrieve all default whitelabel settings.
type DefaultWhitelabelSettings interface {
	Create(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error
	Update(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error
	Delete(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) error
	GetAllWithFilters(
		ctx context.Context,
		skip int32,
		take int32,
		filters map[string]string,
	) ([]*model.DefaultWhitelabelSettingsV2, int, error)
	GetAllWithoutPaginationWithFilters(
		ctx context.Context,
		filters map[string]string,
	) ([]*model.DefaultWhitelabelSettingsV2, error)
	GetV2(ctx context.Context, wlID uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error)
	GetV2ByCoin(ctx context.Context, product, coin string, wlID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)
	Get(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)
}

// defaultWhitelabelSettings represents the default implementation of the DefaultWhitelabelSettings interface.
type defaultWhitelabelSettings struct {
	repo repository.DefaultWhitelabelSettings
}

// NewDefaultWhitelabelSettings is a function that creates a new instance of the DefaultWhitelabelSettings interface.
// It takes a parameter of type repository.DefaultWhitelabelSettings that implements the interface.
// It returns an instance of the DefaultWhitelabelSettings interface implemented by defaultWhitelabelSettings struct.
func NewDefaultWhitelabelSettings(repo repository.DefaultWhitelabelSettings) DefaultWhitelabelSettings {
	return &defaultWhitelabelSettings{repo: repo}
}

// Create is a method of the defaultWhitelabelSettings struct that creates a new default whitelabel setting in the repository.
// It takes a context.Context and a *model.DefaultWhitelabelSettings as parameters.
// It returns an common if the creation of the default whitelabel setting fails.
func (s *defaultWhitelabelSettings) Create(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	return s.repo.Create(ctx, in)
}

// Update is a method that updates the default whitelabel settings for a product and coin.
// It takes a context.Context argument and a *model.DefaultWhitelabelSettings argument, and returns an common.
// It calls the Update method of the repo field to perform the update.
func (s *defaultWhitelabelSettings) Update(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	return s.repo.Update(ctx, in)
}

// Delete is a method of the defaultWhitelabelSettings struct that deletes a default whitelabel setting from the repository.
// It takes a context.Context, a product string, a coin string, and a whitelabelID string as parameters.
// It returns an common if the deletion of the default whitelabel setting fails.
func (s *defaultWhitelabelSettings) Delete(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) error {
	return s.repo.Delete(ctx, product, coin, whitelabelID)
}

// GetAll is a method of the defaultWhitelabelSettings struct that retrieves a list of default whitelabel settings from the repository.
// It takes a context.Context, a skip int32, and a take int32 as parameters.
// It returns a slice of *model.DefaultWhitelabelSettings, the total count of default whitelabel settings, and an common if the retrieval fails.
func (s *defaultWhitelabelSettings) GetAllWithFilters(
	ctx context.Context,
	skip int32,
	take int32,
	filters map[string]string,
) ([]*model.DefaultWhitelabelSettingsV2, int, error) {
	if take <= 0 {
		take = defaultTake
	}

	return s.repo.GetAllWithFilters(ctx, skip, take, filters)
}

func (s *defaultWhitelabelSettings) GetAllWithoutPaginationWithFilters(
	ctx context.Context,
	filters map[string]string,
) ([]*model.DefaultWhitelabelSettingsV2, error) {
	return s.repo.GetAllWithoutPaginationWithFilters(ctx, filters)
}

func (s *defaultWhitelabelSettings) GetV2(ctx context.Context, wlID uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error) {
	settings, err := s.repo.GetV2(ctx, wlID)
	if err != nil {
		return nil, fmt.Errorf("defaultWhitelabelSettings.GetV2: %w", err)
	}
	return settings, nil
}

func (s *defaultWhitelabelSettings) GetV2ByCoin(
	ctx context.Context,
	product, coin string,
	wlID uuid.UUID,
) (*model.DefaultWhitelabelSettingsV2, error) {
	settings, err := s.repo.GetV2ByCoin(ctx, product, coin, wlID)
	if err != nil {
		return nil, fmt.Errorf("defaultWhitelabelSettings.GetV2ByCoin: %w", err)
	}
	return settings, nil
}

// Get retrieves the default settings for a given product and coin from the repository.
// It takes a context, product, and coin as input and returns a pointer to a model.DefaultWhitelabelSettings and an common.
// The input product and coin specify the settings to be retrieved.
// The method uses the repository.DefaultWhitelabelSettings interface to interact with the underlying data store.
// The method delegates the retrieval logic to the repo.Get method.
func (s *defaultWhitelabelSettings) Get(
	ctx context.Context,
	product string,
	coin string,
	whitelabelID uuid.UUID,
) (*model.DefaultWhitelabelSettingsV2, error) {
	return s.repo.Get(ctx, product, coin, whitelabelID)
}
