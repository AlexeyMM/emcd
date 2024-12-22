package buyer

import (
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
)

type Controller struct {
	buyerpb.UnimplementedInvoiceBuyerServiceServer
	buyerInvoiceService service.BuyerInvoiceService
}

func NewController(buyerInvoiceService service.BuyerInvoiceService) *Controller {
	return &Controller{buyerInvoiceService: buyerInvoiceService}
}
