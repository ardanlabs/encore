package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/usergrp"
	"github.com/ardanlabs/encore/business/web/page"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid tag:metrics
func (s *Service) userGrpToken(ctx context.Context, kid string) (usergrp.Token, error) {
	return s.usrGrp.Token(ctx, kid)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/users tag:metrics tag:authorize_admin_only
func (s *Service) userGrpCreate(ctx context.Context, app usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.usrGrp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) userGrpUpdate(ctx context.Context, userID string, app usergrp.AppUpdateUser) (usergrp.AppUser, error) {
	return s.usrGrp.Update(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) userGrpDelete(ctx context.Context, userID string) error {
	return s.usrGrp.Delete(ctx, userID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users tag:metrics tag:authorize_admin_only
func (s *Service) userGrpQuery(ctx context.Context, qp usergrp.QueryParams) (page.Document[usergrp.AppUser], error) {
	return s.usrGrp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) userGrpQueryByID(ctx context.Context, userID string) (usergrp.AppUser, error) {
	return s.usrGrp.QueryByID(ctx, userID)
}
