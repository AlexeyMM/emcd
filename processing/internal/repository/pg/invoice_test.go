package pg

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/gokit"
	"code.emcdtech.com/b2b/processing/pkg/testkit"
)

type InvoiceTestSuite struct {
	repo *Invoice
	suite.Suite
	merchantID     uuid.UUID
	depositAddress string
}

func (suite *InvoiceTestSuite) SetupSuite() {
	suite.repo = NewInvoice(db)
}

func (suite *InvoiceTestSuite) SetupTest() {
	q := sqlc.New(suite.repo.Runner(context.Background()))

	// Create merchant in the database
	suite.merchantID = uuid.New()
	suite.Require().NoError(q.SaveMerchantID(context.Background(), suite.merchantID))

	// Create deposit address in the database
	suite.depositAddress = "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	depositAddressParams := &sqlc.SaveDepositAddressParams{
		Address:    suite.depositAddress,
		NetworkID:  "mainnet",
		MerchantID: suite.merchantID,
		Available:  true,
	}
	suite.Require().NoError(q.SaveDepositAddress(context.Background(), depositAddressParams))
}

func (suite *InvoiceTestSuite) TearDownTest() {
	suite.Require().NoError(truncateServiceTables(context.Background()))
}

func (suite *InvoiceTestSuite) TestSaveInvoice() {
	invoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice",
		Description:    "This is a test invoice",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		PaidAmount:     decimal.NewFromFloat(0.0),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
		CreatedAt:      time.Now(),
	}

	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice))

	retrievedInvoice, err := suite.repo.GetInvoice(context.Background(), invoice.ID)
	suite.Require().NoError(err)

	testkit.RequireEqualCmp(suite.T(), invoice, retrievedInvoice, cmpopts.EquateApproxTime(time.Second))
}

func (suite *InvoiceTestSuite) TestCreateInvoiceForm() {
	form := &model.InvoiceForm{
		ID:          uuid.New(),
		MerchantID:  suite.merchantID,
		Title:       nil, // optional field not filled
		Description: gokit.Ptr("This is a test form"),
		CoinID:      nil, // optional field not filled
		NetworkID:   gokit.Ptr("mainnet"),
		Amount:      gokit.Ptr(decimal.NewFromFloat(0.01)),
		BuyerEmail:  nil, // optional field not filled
		CheckoutURL: "http://example.com/checkout",
		ExternalID:  gokit.Ptr("ext-123"),
		ExpiresAt:   gokit.Ptr(time.Now().Add(24 * time.Hour)),
	}

	suite.Require().NoError(suite.repo.CreateInvoiceForm(context.Background(), form))

	retrievedForm, err := suite.repo.GetInvoiceForm(context.Background(), form.ID)
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), form, retrievedForm, cmpopts.EquateApproxTime(time.Second))
}

func (suite *InvoiceTestSuite) TestCreateInvoiceForm_NoSuchMerchant() {
	form := &model.InvoiceForm{
		ID:          uuid.New(),
		MerchantID:  uuid.New(), // random non-existent merchant ID
		CheckoutURL: "http://example.com/checkout",
	}

	err := suite.repo.CreateInvoiceForm(context.Background(), form)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, &model.Error{Code: model.ErrorCodeNoSuchMerchant})
}

func (suite *InvoiceTestSuite) TestGetInvoiceForm_NotFound() {
	randomID := uuid.New()
	_, err := suite.repo.GetInvoiceForm(context.Background(), randomID)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, &model.Error{Code: model.ErrorCodeNoSuchInvoiceForm})
}

func (suite *InvoiceTestSuite) TestGetActiveInvoiceByDepositAddressForUpdate_Found() {
	// Create an active invoice first
	invoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice",
		Description:    "This is a test invoice",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
		CreatedAt:      time.Now(),
	}

	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice))

	// Try to get the active invoice
	retrievedInvoice, err := suite.repo.GetActiveInvoiceByDepositAddressForUpdate(
		context.Background(),
		suite.depositAddress,
	)
	suite.Require().NoError(err)

	testkit.RequireEqualCmp(suite.T(), invoice, retrievedInvoice, cmpopts.EquateApproxTime(time.Second))
}

func (suite *InvoiceTestSuite) TestGetActiveInvoiceByDepositAddressForUpdate_NotFound() {
	// Try to get a non-existent invoice
	nonExistentAddress := "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"

	_, err := suite.repo.GetActiveInvoiceByDepositAddressForUpdate(context.Background(), nonExistentAddress)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, &model.Error{Code: model.ErrorCodeNoSuchInvoice})
}

func (suite *InvoiceTestSuite) TestGetActiveInvoiceByDepositAddressForUpdate_InactiveStatus() {
	// Create an invoice with finished status
	invoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice",
		Description:    "This is a test invoice",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusFinished,
		DepositAddress: suite.depositAddress,
	}
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice))

	// Try to get the invoice - should return not found error since status is finished
	_, err := suite.repo.GetActiveInvoiceByDepositAddressForUpdate(context.Background(), suite.depositAddress)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, &model.Error{Code: model.ErrorCodeNoSuchInvoice})
}

func (suite *InvoiceTestSuite) TestUpdateStatusWithFinishedAt() {
	// Create an invoice first
	invoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice",
		Description:    "This is a test invoice",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice))

	// Update the status
	newStatus := model.InvoiceStatusFinished
	err := suite.repo.UpdateStatus(context.Background(), invoice.ID, newStatus)
	suite.Require().NoError(err)

	// Verify the status was updated
	updatedInvoice, err := suite.repo.GetInvoice(context.Background(), invoice.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(newStatus, updatedInvoice.Status)

	testkit.RequireEqualCmp(suite.T(), time.Now(), updatedInvoice.FinishedAt, cmpopts.EquateApproxTime(time.Second))
}

func (suite *InvoiceTestSuite) TestUpdateStatusWithoutFinishedAt() {
	// Create an invoice first
	invoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice",
		Description:    "This is a test invoice",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice))

	// Update the status
	newStatus := model.InvoiceStatusPaymentAccepted
	err := suite.repo.UpdateStatus(context.Background(), invoice.ID, newStatus)
	suite.Require().NoError(err)

	// Verify the status was updated
	updatedInvoice, err := suite.repo.GetInvoice(context.Background(), invoice.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(newStatus, updatedInvoice.Status)
	suite.Require().True(updatedInvoice.FinishedAt.IsZero())
}

func (suite *InvoiceTestSuite) TestSetInvoicesExpired() {
	// Create two invoices - one that should expire and one that shouldn't
	expiredInvoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Expired Invoice",
		Description:    "This invoice should expire",
		ExpiresAt:      time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	activeInvoice := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-456",
		Title:          "Active Invoice",
		Description:    "This invoice should not expire",
		ExpiresAt:      time.Now().Add(24 * time.Hour), // Expires in 24 hours
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	// Create both invoices
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), expiredInvoice))
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), activeInvoice))

	// Run expiration
	suite.Require().NoError(suite.repo.SetInvoicesExpired(context.Background()))

	// Verify expired invoice status was updated
	expiredInvoiceAfter, err := suite.repo.GetInvoice(context.Background(), expiredInvoice.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(model.InvoiceStatusExpired, expiredInvoiceAfter.Status)
	suite.Require().Equal(time.Now().Truncate(time.Second), expiredInvoiceAfter.FinishedAt.Truncate(time.Second))

	// Verify active invoice status was not changed
	activeInvoiceAfter, err := suite.repo.GetInvoice(context.Background(), activeInvoice.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(model.InvoiceStatusWaitingForDeposit, activeInvoiceAfter.Status)
	suite.Require().True(activeInvoiceAfter.FinishedAt.IsZero())
}

func (suite *InvoiceTestSuite) TestGetInvoiceCountByStatus() {
	// Создание инвойсов со статусом "waiting_for_deposit"
	invoice1 := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-123",
		Title:          "Test Invoice 1",
		Description:    "This is a test invoice with status 'waiting_for_deposit'",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.01),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		PaidAmount:     decimal.NewFromFloat(0.0),
		BuyerEmail:     "buyer@example.com",
		CheckoutURL:    "http://example.com/checkout",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	invoice2 := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-124",
		Title:          "Test Invoice 2",
		Description:    "This is a test invoice with status 'finished'",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.02),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		PaidAmount:     decimal.NewFromFloat(0.02),
		BuyerEmail:     "buyer2@example.com",
		CheckoutURL:    "http://example.com/checkout2",
		Status:         model.InvoiceStatusWaitingForDeposit,
		DepositAddress: suite.depositAddress,
	}

	// Создание инвойса с статусом "expired"
	invoice3 := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-125",
		Title:          "Test Invoice 3",
		Description:    "This is a test invoice with status 'expired'",
		ExpiresAt:      time.Now().Add(-24 * time.Hour), // Это уже просроченный инвойс
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.03),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		PaidAmount:     decimal.NewFromFloat(0.0),
		BuyerEmail:     "buyer3@example.com",
		CheckoutURL:    "http://example.com/checkout3",
		Status:         model.InvoiceStatusExpired,
		DepositAddress: suite.depositAddress,
	}

	// Создание инвойса с статусом "cancelled"
	invoice4 := &model.Invoice{
		ID:             uuid.New(),
		MerchantID:     suite.merchantID,
		ExternalID:     "ext-126",
		Title:          "Test Invoice 4",
		Description:    "This is a test invoice with status 'cancelled'",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		CoinID:         "BTC",
		NetworkID:      "mainnet",
		PaymentAmount:  decimal.NewFromFloat(0.04),
		BuyerFee:       decimal.NewFromFloat(0.001),
		MerchantFee:    decimal.NewFromFloat(0.002),
		PaidAmount:     decimal.NewFromFloat(0.0),
		BuyerEmail:     "buyer4@example.com",
		CheckoutURL:    "http://example.com/checkout4",
		Status:         model.InvoiceStatusCancelled,
		DepositAddress: suite.depositAddress,
	}

	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice1))
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice2))
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice3))
	suite.Require().NoError(suite.repo.CreateInvoice(context.Background(), invoice4))

	statuses, err := suite.repo.CountInvoiceByStatus(context.Background())
	suite.Require().NoError(err)

	suite.Require().Len(statuses, 3)

	expectedStatuses := map[model.InvoiceStatus]int{
		model.InvoiceStatusWaitingForDeposit: 2,
		model.InvoiceStatusExpired:           1,
		model.InvoiceStatusCancelled:         1,
	}
	suite.Require().Equal(expectedStatuses, statuses)
}

func TestInvoiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvoiceTestSuite))
}
