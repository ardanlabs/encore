package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/usergrp"
	"github.com/ardanlabs/encore/business/web"
)

//encore:api auth method=GET path=/v1/token/:kid
func (s *Service) userGrp_Token(ctx context.Context, kid string) (usergrp.Token, error) {
	return s.UsrGrp.Token(ctx, kid)
}

//encore:api auth method=POST path=/v1/users tag:authorize_admin_only
func (s *Service) userGrp_Create(ctx context.Context, anu usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.UsrGrp.Create(ctx, anu)
}

//encore:api auth method=PUT path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrp_Update(ctx context.Context, userID string, auu usergrp.AppUpdateUser) (usergrp.AppUser, error) {
	return s.UsrGrp.Update(ctx, userID, auu)
}

//encore:api auth method=DELETE path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrp_Delete(ctx context.Context, userID string) error {
	return s.UsrGrp.Delete(ctx, userID)
}

//encore:api auth method=GET path=/v1/users tag:authorize_admin_only
func (s *Service) userGrp_Query(ctx context.Context, qp usergrp.QueryParams) (web.PageDocument[usergrp.AppUser], error) {
	return s.UsrGrp.Query(ctx, qp)
}

//encore:api auth method=GET path=/v1/users/:userID tag:authorize_user
func (s *Service) userGrp_QueryById(ctx context.Context, userID string) (usergrp.AppUser, error) {
	return s.UsrGrp.QueryByID(ctx)
}
