package merchant

import (
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/protocol/merchantpb"
)

type Controller struct {
	merchantpb.UnimplementedInvoiceMerchantServiceServer
	merchantInvoiceService service.MerchantInvoiceService
}

func NewController(merchantInvoiceService service.MerchantInvoiceService) *Controller {
	return &Controller{merchantInvoiceService: merchantInvoiceService}
}
