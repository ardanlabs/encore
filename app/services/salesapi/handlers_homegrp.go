package salesapi

import (
	"context"

	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/homegrp"
	"github.com/ardanlabs/encore/business/web/page"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/homes tag:metrics tag:authorize_user_only
func (s *Service) HomeGrpCreate(ctx context.Context, app homegrp.AppNewHome) (homegrp.AppHome, error) {
	return s.hmeGrp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeGrpUpdate(ctx context.Context, homeID string, app homegrp.AppUpdateHome) (homegrp.AppHome, error) {
	return s.hmeGrp.Update(ctx, homeID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeGrpDelete(ctx context.Context, homeID string) error {
	return s.hmeGrp.Delete(ctx, homeID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes tag:metrics tag:authorize_any
func (s *Service) HomeGrpQuery(ctx context.Context, qp homegrp.QueryParams) (page.Document[homegrp.AppHome], error) {
	return s.hmeGrp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes/:productID tag:metrics tag:authorize_home
func (s *Service) HomeGrpQueryByID(ctx context.Context, productID string) (homegrp.AppHome, error) {
	return s.hmeGrp.QueryByID(ctx, productID)
}
