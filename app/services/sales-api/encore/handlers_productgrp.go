package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/productgrp"
	"github.com/ardanlabs/encore/business/web/page"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/products tag:authorize_user_only
func (s *Service) productGrpCreate(ctx context.Context, app productgrp.AppNewProduct) (productgrp.AppProduct, error) {
	return s.PrdGrp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/products/:productID tag:authorize_product
func (s *Service) productGrpUpdate(ctx context.Context, productID string, app productgrp.AppUpdateProduct) (productgrp.AppProduct, error) {
	return s.PrdGrp.Update(ctx, productID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/products/:productID tag:authorize_product
func (s *Service) productGrpDelete(ctx context.Context, productID string) error {
	return s.PrdGrp.Delete(ctx, productID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products tag:authorize_any
func (s *Service) productGrpQuery(ctx context.Context, qp productgrp.QueryParams) (page.Document[productgrp.AppProduct], error) {
	return s.PrdGrp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products/:productID tag:authorize_product
func (s *Service) productGrpQueryByID(ctx context.Context, productID string) (productgrp.AppProduct, error) {
	return s.PrdGrp.QueryByID(ctx, productID)
}
