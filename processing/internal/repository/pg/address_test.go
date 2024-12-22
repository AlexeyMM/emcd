package pg

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AddressTestSuite struct {
	suite.Suite
	repo       *DepositAddressPool
	merchantID uuid.UUID
}

func (suite *AddressTestSuite) SetupSuite() {
	suite.repo = NewDepositAddressPool(db)
}

func (suite *AddressTestSuite) SetupTest() {
	q := sqlc.New(suite.repo.Runner(context.Background()))

	// Create merchant in the database
	suite.merchantID = uuid.New()
	suite.Require().NoError(q.SaveMerchantID(context.Background(), suite.merchantID))
}

func (suite *AddressTestSuite) TearDownTest() {
	suite.Require().NoError(truncateServiceTables(context.Background()))
}

func (suite *AddressTestSuite) TestOccupyAddress_Success() {
	address := &model.Address{
		Address:    "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		NetworkID:  "mainnet",
		MerchantID: suite.merchantID,
		Available:  true,
	}

	suite.Require().NoError(suite.repo.Save(context.Background(), address))

	retrievedAddress, err := suite.repo.OccupyAddress(context.Background(), suite.merchantID, "mainnet")
	suite.Require().NoError(err)
	suite.Require().Equal(address, retrievedAddress)
}

func (suite *AddressTestSuite) TestOccupyAddress_NoAvailableAddress() {
	address := &model.Address{
		Address:    "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		NetworkID:  "mainnet",
		MerchantID: suite.merchantID,
		Available:  false,
	}

	suite.Require().NoError(suite.repo.Save(context.Background(), address))

	_, err := suite.repo.OccupyAddress(context.Background(), suite.merchantID, "mainnet")
	suite.Require().ErrorIs(err, &model.Error{Code: model.ErrorCodeNoAvailableAddress})
}

func (suite *AddressTestSuite) TestSave() {
	address := &model.Address{
		Address:    "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		NetworkID:  "mainnet",
		MerchantID: suite.merchantID,
		Available:  true,
	}

	suite.Require().NoError(suite.repo.Save(context.Background(), address))

	retrievedAddress, err := suite.repo.OccupyAddress(context.Background(), suite.merchantID, "mainnet")
	suite.Require().NoError(err)
	suite.Require().Equal(address, retrievedAddress)
}

func TestAddressTestSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}
