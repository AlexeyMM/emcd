package service

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"code.emcdtech.com/b2b/endpoint/internal/repository"
	"github.com/google/uuid"
)

type Secret interface {
	GenerateKeys(ctx context.Context, clientID uuid.UUID) (*model.Secret, error)
	RotateKeys(ctx context.Context, apiKey uuid.UUID) (*model.Secret, error)
	GetActiveKeys(ctx context.Context, clientID uuid.UUID) ([]*model.Secret, error)
	DeactivateKey(ctx context.Context, clientID, apiKey uuid.UUID) error
	DeactivateAllKeys(ctx context.Context, clientID uuid.UUID) error
}

type secret struct {
	repo   repository.Secret
	ipRepo repository.IP
}

func NewSecret(repo repository.Secret, ipRepo repository.IP) *secret {
	return &secret{
		repo:   repo,
		ipRepo: ipRepo,
	}
}

func (s *secret) GenerateKeys(ctx context.Context, clientID uuid.UUID) (*model.Secret, error) {
	return s.generateKeys(ctx, clientID)
}

func (s *secret) generateKeys(ctx context.Context, clientID uuid.UUID) (*model.Secret, error) {
	newKey := model.Secret{
		ClientID:  clientID,
		ApiKey:    uuid.New(),
		ApiSecret: uuid.New(),
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		LastUsed:  time.Now().UTC(),
	}

	err := s.repo.Add(ctx, &newKey)
	if err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}

	return &newKey, nil
}

func (s *secret) RotateKeys(ctx context.Context, apiKey uuid.UUID) (*model.Secret, error) {
	oldSecret, err := s.repo.FindOne(ctx, &model.SecretFilter{
		ApiKey: &apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("findOne: %w", err)
	}

	isActive := false
	err = s.repo.Update(ctx, oldSecret,
		&model.SecretFilter{
			ApiKey: &oldSecret.ApiKey,
		}, &model.SecretPartial{
			IsActive: &isActive,
		})
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	sec, err := s.generateKeys(ctx, oldSecret.ClientID)
	if err != nil {
		return nil, fmt.Errorf("generateKeys: %w", err)
	}

	// Копируем все разрешённые ip для нового apiKey
	ips, err := s.ipRepo.Find(ctx, &model.IPFilter{
		ApiKey: &apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	for _, one := range ips {
		one.ApiKey = sec.ApiKey
	}
	err = s.ipRepo.Add(ctx, ips)
	if err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}

	return sec, nil
}

func (s *secret) GetActiveKeys(ctx context.Context, clientID uuid.UUID) ([]*model.Secret, error) {
	isActive := true
	secrets, err := s.repo.Find(ctx, &model.SecretFilter{
		ClientID: &clientID,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("findOne: %w", err)
	}
	return secrets, nil
}

func (s *secret) DeactivateKey(ctx context.Context, clientID, apiKey uuid.UUID) error {
	isActive := false
	err := s.repo.Update(ctx, &model.Secret{},
		&model.SecretFilter{
			ApiKey:   &apiKey,
			ClientID: &clientID,
		},
		&model.SecretPartial{
			IsActive: &isActive,
		})
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *secret) DeactivateAllKeys(ctx context.Context, clientID uuid.UUID) error {
	isActive := false
	err := s.repo.Update(ctx, &model.Secret{},
		&model.SecretFilter{
			ClientID: &clientID,
		},
		&model.SecretPartial{
			IsActive: &isActive,
		})
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
