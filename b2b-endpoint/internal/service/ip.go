package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"code.emcdtech.com/b2b/endpoint/internal/repository"
	"github.com/google/uuid"
)

type IP interface {
	AddIPs(ctx context.Context, ips []*model.IP) error
	GetIPs(ctx context.Context, apiKey uuid.UUID) ([]*model.IP, error)
	DeleteIP(ctx context.Context, apiKey uuid.UUID, ip string) error
	DeleteAllIPs(ctx context.Context, apiKey uuid.UUID) error
}

type ip struct {
	ipRepo repository.IP
}

func NewIP(ipRepo repository.IP) *ip {
	return &ip{
		ipRepo: ipRepo,
	}
}

func (i *ip) AddIPs(ctx context.Context, ips []*model.IP) error {
	err := i.ipRepo.Add(ctx, ips)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}
	return nil
}

func (i *ip) GetIPs(ctx context.Context, apiKey uuid.UUID) ([]*model.IP, error) {
	ips, err := i.ipRepo.Find(ctx, &model.IPFilter{
		ApiKey: &apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	return ips, nil
}

func (i *ip) DeleteIP(ctx context.Context, apiKey uuid.UUID, ip string) error {
	err := i.ipRepo.Delete(ctx, &model.IPFilter{
		ApiKey:  &apiKey,
		Address: &ip,
	})
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (i *ip) DeleteAllIPs(ctx context.Context, apiKey uuid.UUID) error {
	err := i.ipRepo.Delete(ctx, &model.IPFilter{
		ApiKey: &apiKey,
	})
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
