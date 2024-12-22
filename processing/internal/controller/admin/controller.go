package admin

import (
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/protocol/adminpb"
)

type Controller struct {
	adminpb.UnimplementedMerchantAdminServiceServer
	merchantAdminService service.MerchantAdminService
}

func NewController(merchantAdminService service.MerchantAdminService) *Controller {
	return &Controller{merchantAdminService: merchantAdminService}
}
