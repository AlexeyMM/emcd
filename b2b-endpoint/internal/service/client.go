package service

import (
	"context"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"code.emcdtech.com/b2b/endpoint/internal/repository"
	"github.com/google/uuid"
)

type Client interface {
	Add(ctx context.Context, name string) (uuid.UUID, error)
}

type client struct {
	repo repository.Client
}

func NewClient(repo repository.Client) *client {
	return &client{
		repo: repo,
	}
}

func (c *client) Add(ctx context.Context, name string) (uuid.UUID, error) {
	clientID := uuid.New()
	err := c.repo.Add(ctx, &model.Client{
		ID:   clientID,
		Name: name,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return clientID, nil
}
