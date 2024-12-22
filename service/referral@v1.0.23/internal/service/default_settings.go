// Package service defines a set of functions for performing operations related to a service.
package service

import (
	"context"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

const defaultTake int32 = 10

// DefaultSettings is an interface that defines the operations for managing default settings.
// The operations include creating, updating, deleting, and retrieving default settings,
// as well as getting all default settings with pagination support.
type DefaultSettings interface {
	Create(ctx context.Context, in *model.DefaultSettings) error
	Update(ctx context.Context, in *model.DefaultSettings) error
	Delete(ctx context.Context, product string, coin string) error
	Get(ctx context.Context, product string, coin string) (*model.DefaultSettings, error)
	GetAll(ctx context.Context, skip int32, take int32) ([]*model.DefaultSettings, int, error)
	GetAllWithoutPagination(ctx context.Context, referrerUUID string) ([]*model.DefaultSettings, error)
}

// defaultSettings represents a type that implements the DefaultSettings interface.
type defaultSettings struct {
	repo repository.DefaultSettings
}

// NewDefaultSettings creates a new DefaultSettings using the provided repository.DefaultSettings.
// It takes the repository as input and returns a DefaultSettings interface.
// The DefaultSettings struct contains a reference to the repository.
// The function is used to initialize a DefaultSettings instance with the repository implementation.
func NewDefaultSettings(repo repository.DefaultSettings) DefaultSettings {
	return &defaultSettings{repo: repo}
}

// Create creates a new default settings record in the repository.
// It takes a context and a pointer to a model.DefaultSettings as input and returns an common.
// The input model.DefaultSettings contains information about the product, coin, fee, referral fee, and creation timestamp.
// It uses the repository.DefaultSettings interface to interact with the underlying data store.
// The method delegates the creation logic to the repo.Create method.
func (s *defaultSettings) Create(ctx context.Context, in *model.DefaultSettings) error {
	return s.repo.Create(ctx, in)
}

// Update updates the default settings in the repository based on the provided context and input.
// It returns an common if the update operation fails.
// The input should be a pointer to a model.DefaultSettings object.
func (s *defaultSettings) Update(ctx context.Context, in *model.DefaultSettings) error {
	return s.repo.Update(ctx, in)
}

// Delete deletes the default settings for a specific product and coin.
// It calls the Delete method of the repository.DefaultSettings interface to perform the deletion.
// The method receives the context, the product, and the coin as parameters.
// It returns an common if the deletion fails.
func (s *defaultSettings) Delete(ctx context.Context, product string, coin string) error {
	return s.repo.Delete(ctx, product, coin)
}

// Get retrieves the default settings for a given product and coin from the repository.
// It takes a context, product, and coin as input and returns a pointer to a model.DefaultSettings and an common.
// The input product and coin specify the settings to be retrieved.
// The method uses the repository.DefaultSettings interface to interact with the underlying data store.
// The method delegates the retrieval logic to the repo.Get method.
func (s *defaultSettings) Get(ctx context.Context, product string, coin string) (*model.DefaultSettings, error) {
	return s.repo.Get(ctx, product, coin)
}

// GetAll fetches all default settings records from the repository.
// It takes a context, a skip value indicating the number of records to skip, and a take value indicating the number of records to fetch.
// It returns a slice of pointers to model.DefaultSettings, representing the fetched records, the total number of records, and an common if any.
// Each DefaultSettings object contains information about the product, coin, fee, referral fee, and creation timestamp.
// The method uses the repository.DefaultSettings interface to interact with the underlying data store.
// It delegates the fetch logic to the repo.GetAll method.
func (s *defaultSettings) GetAll(ctx context.Context, skip int32, take int32) ([]*model.DefaultSettings, int, error) {
	if take <= 0 {
		take = defaultTake
	}

	return s.repo.GetAll(ctx, skip, take)
}

func (s *defaultSettings) GetAllWithoutPagination(ctx context.Context, referrerUUID string) ([]*model.DefaultSettings, error) {
	var referrerSettings []*model.DefaultSettings
	var err error
	if referrerUUID != "" {
		referrerSettings, err = s.repo.GetSettingByReferrer(ctx, referrerUUID)
		if err != nil {
			return nil, err
		}
	}

	settingsByCoin, err := s.repo.GetAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	if len(referrerSettings) == 0 {
		return settingsByCoin, nil
	}

	referrerMap := makeMapFromDefaultSettings(referrerSettings)
	defaultSettingsMap := makeMapFromDefaultSettings(settingsByCoin)
	for key, settings := range defaultSettingsMap {
		_, ok := referrerMap[key]
		if !ok {
			referrerMap[key] = settings
		}
	}

	meltedSettings := make([]*model.DefaultSettings, 0, len(referrerMap))
	for _, setting := range referrerMap {
		meltedSettings = append(meltedSettings, &model.DefaultSettings{
			Product:     setting.Product,
			Coin:        setting.Coin,
			Fee:         setting.Fee,
			ReferralFee: setting.ReferralFee,
			CreatedAt:   setting.CreatedAt,
		})
	}
	return meltedSettings, nil
}

func makeMapFromDefaultSettings(settings []*model.DefaultSettings) map[string]model.DefaultSettings {
	m := make(map[string]model.DefaultSettings, len(settings))
	for _, setting := range settings {
		m[setting.Product+setting.Coin] = model.DefaultSettings{
			Product:     setting.Product,
			Coin:        setting.Coin,
			Fee:         setting.Fee,
			ReferralFee: setting.ReferralFee,
			CreatedAt:   setting.CreatedAt,
		}
	}
	return m
}
