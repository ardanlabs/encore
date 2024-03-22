package salesapi

import (
	"context"

	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/views/vproductgrp"
	"github.com/ardanlabs/encore/business/web/page"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/vproducts tag:metrics tag:authorize_admin_only
func (s *Service) VProductGrpQuery(ctx context.Context, qp vproductgrp.QueryParams) (page.Document[vproductgrp.AppProduct], error) {
	return s.vprdGrp.Query(ctx, qp)
}
