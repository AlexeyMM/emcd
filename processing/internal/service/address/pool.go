package address

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type PoolService struct {
	addressServiceClient address.AddressServiceClient
	depositAddressPool   repository.DepositAddressPool
}

func NewPoolService(
	addressServiceClient address.AddressServiceClient,
	depositAddressPool repository.DepositAddressPool,
) *PoolService {
	return &PoolService{addressServiceClient: addressServiceClient, depositAddressPool: depositAddressPool}
}

func (s *PoolService) GetOrCreate(
	ctx context.Context,
	merchantID uuid.UUID,
	networkID string,
	idempotencyKey uuid.UUID,
) (string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "get or create address")
	defer span.End()

	a, err := s.depositAddressPool.OccupyAddress(ctx, merchantID, networkID)
	if err == nil {
		return a.Address, nil
	}

	if !errors.Is(err, &model.Error{Code: model.ErrorCodeNoAvailableAddress}) {
		return "", fmt.Errorf("occupyAddress: %w", err)
	}

	newAddress, err := s.createNewAddressAndOccupy(ctx, merchantID, networkID, idempotencyKey)
	if err != nil {
		return "", fmt.Errorf("createNewAddressAndOccupy: %w", err)
	}

	return newAddress, nil
}

func (s *PoolService) createNewAddressAndOccupy(
	ctx context.Context,
	merchantID uuid.UUID,
	networkID string,
	idempotencyKey uuid.UUID,
) (string, error) {
	resp, err := s.addressServiceClient.CreateProcessingAddress(ctx, &address.CreateProcessingAddressRequest{
		UserUuid:       merchantID.String(),
		Network:        networkID,
		ProcessingUuid: idempotencyKey.String(),
	})
	if err != nil {
		return "", fmt.Errorf("createProcessingAddress: %w", err)
	}

	err = s.depositAddressPool.Save(ctx, &model.Address{
		Address:    resp.GetAddress(),
		NetworkID:  networkID,
		MerchantID: merchantID,
		Available:  false, // it is occupied by the invoice for which the address was created
	})
	if err != nil {
		return "", fmt.Errorf("createNewAddressAndOccupy: %w", err)
	}

	return resp.GetAddress(), nil
}
