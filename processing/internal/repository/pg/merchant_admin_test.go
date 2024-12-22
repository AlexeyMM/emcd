package pg

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/processing/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type MerchantAdminTestSuite struct {
	adminRepo  *MerchantAdmin
	commonRepo *Merchant
	suite.Suite
}

func (suite *MerchantAdminTestSuite) SetupSuite() {
	suite.adminRepo = NewMerchantAdmin(db)
	suite.commonRepo = NewMerchant(db)
}

func (suite *MerchantAdminTestSuite) TearDownTest() {
	suite.Require().NoError(truncateServiceTables(context.Background()))
}

func (suite *MerchantAdminTestSuite) TestSaveMerchant() {
	m := &model.Merchant{
		ID: uuid.New(),
		Tariff: &model.Tariff{
			UpperFee: decimal.RequireFromString("1"),
			LowerFee: decimal.RequireFromString("2"),
			MinPay:   decimal.RequireFromString("3"),
			MaxPay:   decimal.RequireFromString("4"),
		},
	}

	suite.Require().NoError(suite.adminRepo.SaveMerchant(context.Background(), m))

	retrievedMerchant, err := suite.commonRepo.Get(context.Background(), m.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(m, retrievedMerchant)
}

func TestMerchantAdminTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantAdminTestSuite))
}
