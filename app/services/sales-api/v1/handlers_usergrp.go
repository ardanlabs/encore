package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/usergrp"
)

//encore:api auth method=GET path=/v1/users/token/:kid
func (s *Service) UserGrp_Token(ctx context.Context, kid string) (usergrp.Token, error) {
	return s.usrGrp.Token(ctx, kid)
}

//encore:api auth method=POST path=/v1/users tag:auth_admin_only
func (s *Service) UserGrp_Create(ctx context.Context, anu usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.usrGrp.Create(ctx, anu)
}

//encore:api auth method=PUT path=/v1/users/:userID tag:auth_user
func (s *Service) UserGrp_Update(ctx context.Context, userID string, auu usergrp.AppUpdateUser) (usergrp.AppUser, error) {
	return s.usrGrp.Update(ctx, userID, auu)
}

//encore:api auth method=DELETE path=/v1/users/:userID tag:auth_user
func (s *Service) UserGrp_Delete(ctx context.Context, userID string) error {
	return s.usrGrp.Delete(ctx, userID)
}
