package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/usergrp"
)

//encore:api auth method=POST path=/v1/users tag:auth_admin_only
func (s *Service) UserGrp_Create(ctx context.Context, app usergrp.AppNewUser) (usergrp.AppUser, error) {
	return s.usrGrp.Create(ctx, app)
}
