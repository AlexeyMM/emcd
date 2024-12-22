package common

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/protocol/commonpb"
)

func ConvertModelInvoiceToProto(invoice *model.Invoice) *commonpb.Invoice {
	protoInvoice := &commonpb.Invoice{
		Id:                    invoice.ID.String(),
		ExternalId:            invoice.ExternalID,
		Title:                 invoice.Title,
		Description:           invoice.Description,
		ExpiresAt:             timestamppb.New(invoice.ExpiresAt),
		CoinId:                invoice.CoinID,
		NetworkId:             invoice.NetworkID,
		PaymentAmount:         invoice.PaymentAmount.String(),
		PaidAmount:            invoice.PaidAmount.String(),
		BuyerEmail:            invoice.BuyerEmail,
		CheckoutUrl:           invoice.CheckoutURL,
		Status:                convertModelInvoiceStatusToProto(invoice.Status),
		DepositAddress:        invoice.DepositAddress,
		BuyerFee:              invoice.BuyerFee.String(),
		MerchantFee:           invoice.MerchantFee.String(),
		RequiredPaymentAmount: invoice.RequiredPayment.String(),
		CreatedAt:             timestamppb.New(invoice.CreatedAt),
		FinishedAt:            timestamppb.New(invoice.FinishedAt),
	}

	protoInvoice.Transactions = make([]*commonpb.Transaction, 0, len(invoice.Transactions))
	for _, tx := range invoice.Transactions {
		protoInvoice.Transactions = append(protoInvoice.Transactions, &commonpb.Transaction{
			Hash:        tx.Hash,
			Amount:      tx.Amount.String(),
			Address:     tx.Address,
			IsConfirmed: tx.IsConfirmed,
			CreatedAt:   timestamppb.New(tx.CreatedAt),
		})
	}

	return protoInvoice
}

var modelInvoiceStatusToProto = map[model.InvoiceStatus]commonpb.Invoice_InvoiceStatus{
	model.InvoiceStatusUnknown:             commonpb.Invoice_UNKNOWN,
	model.InvoiceStatusWaitingForDeposit:   commonpb.Invoice_WAITING_FOR_DEPOSIT,
	model.InvoiceStatusPaymentConfirmation: commonpb.Invoice_PAYMENT_ACCEPTED,
	model.InvoiceStatusPartiallyPaid:       commonpb.Invoice_PARTIALLY_PAID,
	model.InvoiceStatusPaymentAccepted:     commonpb.Invoice_PAYMENT_ACCEPTED,
	model.InvoiceStatusFinished:            commonpb.Invoice_FINISHED,
	model.InvoiceStatusCancelled:           commonpb.Invoice_CANCELLED,
	model.InvoiceStatusExpired:             commonpb.Invoice_EXPIRED,
}

func convertModelInvoiceStatusToProto(status model.InvoiceStatus) commonpb.Invoice_InvoiceStatus {
	if protoStatus, exists := modelInvoiceStatusToProto[status]; exists {
		return protoStatus
	}

	return commonpb.Invoice_UNKNOWN
}
