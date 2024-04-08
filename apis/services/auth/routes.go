package auth

import (
	"context"

	"github.com/ardanlabs/encore/app/api/mid"
)

//lint:ignore U1000 "called by encore"
//encore:api method=GET path=/v1/auth
func (s *Service) AuthHandler(ctx context.Context, ap *mid.AuthParams) error {
	s.log.Info("auth-handler", "status", "started")
	return nil
}
