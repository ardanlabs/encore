package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/usergrp"
	v1 "github.com/ardanlabs/encore/business/web/v1"
)

//encore:api auth method=GET path=/v1/token/:kid
func (s *Service) userGrp_Token(ctx context.Context, kid string) (usergrp.Token, error) {
	return s.usrGrp.Token(ctx, kid)
}

//encore:api auth method=POST path=/v1/users tag:auth_admin_only
func (s *Service) userGrp_Create(ctx context.Context, anu usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.usrGrp.Create(ctx, anu)
}

//encore:api auth method=PUT path=/v1/users/:userID tag:auth_user
func (s *Service) userGrp_Update(ctx context.Context, userID string, auu usergrp.AppUpdateUser) (usergrp.AppUser, error) {
	return s.usrGrp.Update(ctx, userID, auu)
}

//encore:api auth method=DELETE path=/v1/users/:userID tag:auth_user
func (s *Service) userGrp_Delete(ctx context.Context, userID string) error {
	return s.usrGrp.Delete(ctx, userID)
}

//encore:api auth method=GET path=/v1/users tag:auth_admin_only
func (s *Service) userGrp_Query(ctx context.Context, qp usergrp.QueryParams) (v1.PageDocument[usergrp.AppUser], error) {
	return s.usrGrp.Query(ctx, qp)
}

//encore:api auth method=GET path=/v1/users/:userID tag:auth_user
func (s *Service) userGrp_QueryById(ctx context.Context, userID string) (usergrp.AppUser, error) {
	return s.usrGrp.QueryByID(ctx)
}
