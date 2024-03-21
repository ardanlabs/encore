package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/usergrp"
	"github.com/ardanlabs/encore/business/web/page"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid
func (s *Service) userGrpToken(ctx context.Context, kid string) (usergrp.Token, error) {
	return s.UsrGrp.Token(ctx, kid)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/users tag:authorize_admin_only
func (s *Service) userGrpCreate(ctx context.Context, anu usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.UsrGrp.Create(ctx, anu)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrpUpdate(ctx context.Context, userID string, auu usergrp.AppUpdateUser) (usergrp.AppUser, error) {
	return s.UsrGrp.Update(ctx, userID, auu)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrpDelete(ctx context.Context, userID string) error {
	return s.UsrGrp.Delete(ctx, userID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users tag:authorize_admin_only
func (s *Service) userGrpQuery(ctx context.Context, qp usergrp.QueryParams) (page.Document[usergrp.AppUser], error) {
	return s.UsrGrp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrpQueryByID(ctx context.Context, userID string) (usergrp.AppUser, error) {
	return s.UsrGrp.QueryByID(ctx, userID)
}
