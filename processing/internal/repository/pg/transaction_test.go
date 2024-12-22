package pg

import (
	"context"
	"testing"
	"time"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/testkit"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	suite.Suite
	repo           *Transaction
	merchantID     uuid.UUID
	invoiceID      uuid.UUID
	depositAddress string
}

func (suite *TransactionTestSuite) SetupSuite() {
	suite.repo = NewTransaction(db)
}

func (suite *TransactionTestSuite) SetupTest() {
	q := sqlc.New(suite.repo.Runner(context.Background()))

	// Create merchant in the database
	suite.merchantID = uuid.New()
	suite.Require().NoError(q.SaveMerchantID(context.Background(), suite.merchantID))

	// Create deposit address
	suite.depositAddress = "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	depositAddressParams := &sqlc.SaveDepositAddressParams{
		Address:    suite.depositAddress,
		NetworkID:  "eth",
		MerchantID: suite.merchantID,
		Available:  true,
	}
	suite.Require().NoError(q.SaveDepositAddress(context.Background(), depositAddressParams))

	// Create an invoice
	suite.invoiceID = uuid.New()
	invoice := &sqlc.CreateInvoiceParams{
		ID:             suite.invoiceID,
		MerchantID:     suite.merchantID,
		CoinID:         "eth",
		NetworkID:      "eth",
		DepositAddress: suite.depositAddress,
		Amount:         decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		Title:          "Test Invoice",
		Description:    "Test Description",
		CheckoutUrl:    "http://example.com",
		Status:         sqlc.InvoiceStatusWaitingForDeposit,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		ExternalID:     "ext-123",
		BuyerEmail:     "test@example.com",
	}
	suite.Require().NoError(q.CreateInvoice(context.Background(), invoice))
}

func (suite *TransactionTestSuite) TearDownTest() {
	suite.Require().NoError(truncateServiceTables(context.Background()))
}

func (suite *TransactionTestSuite) TestGetInvoiceTransactions() {
	expectedTxs := []*model.Transaction{
		{
			InvoiceID:   suite.invoiceID,
			Hash:        "0x123abc",
			Amount:      decimal.NewFromFloat(0.01),
			Address:     suite.depositAddress,
			IsConfirmed: false,
			CreatedAt:   time.Now(),
		},
		{
			InvoiceID:   suite.invoiceID,
			Hash:        "0x456def",
			Amount:      decimal.NewFromFloat(0.02),
			Address:     suite.depositAddress,
			IsConfirmed: true,
			CreatedAt:   time.Now(),
		},
	}

	// Save all transactions
	for _, tx := range expectedTxs {
		suite.Require().NoError(suite.repo.SaveTransaction(context.Background(), tx))
	}

	// Retrieve and verify transactions
	actualTxs, err := suite.repo.GetInvoiceTransactions(context.Background(), suite.invoiceID)
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), expectedTxs, actualTxs, cmpopts.EquateApproxTime(time.Second))
}

func (suite *TransactionTestSuite) TestSaveTransaction_NewTransaction() {
	tx := &model.Transaction{
		InvoiceID:   suite.invoiceID,
		Hash:        "0x123abc",
		Amount:      decimal.NewFromFloat(0.01),
		Address:     suite.depositAddress,
		IsConfirmed: false,
		CreatedAt:   time.Now(),
	}

	err := suite.repo.SaveTransaction(context.Background(), tx)
	suite.Require().NoError(err)

	// Verify transaction was saved
	transactions, err := suite.repo.GetInvoiceTransactions(context.Background(), suite.invoiceID)
	suite.Require().NoError(err)
	suite.Require().Len(transactions, 1)
	testkit.RequireEqualCmp(suite.T(), tx, transactions[0], cmpopts.EquateApproxTime(time.Second))
}

func (suite *TransactionTestSuite) TestSaveTransaction_UpdateExisting() {
	// Create initial transaction
	tx := &model.Transaction{
		InvoiceID:   suite.invoiceID,
		Hash:        "0x123abc",
		Amount:      decimal.NewFromFloat(0.01),
		Address:     suite.depositAddress,
		IsConfirmed: false,
		CreatedAt:   time.Now(),
	}
	suite.Require().NoError(suite.repo.SaveTransaction(context.Background(), tx))

	// Modify and save the transaction
	tx.Amount = decimal.NewFromFloat(0.02)
	tx.IsConfirmed = true

	err := suite.repo.SaveTransaction(context.Background(), tx)
	suite.Require().NoError(err)

	// Verify transaction was updated
	transactions, err := suite.repo.GetInvoiceTransactions(context.Background(), suite.invoiceID)
	suite.Require().NoError(err)
	suite.Require().Len(transactions, 1)
	testkit.RequireEqualCmp(suite.T(), tx, transactions[0], cmpopts.EquateApproxTime(time.Second))
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
