package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/testgrp"
)

//encore:api public method=GET path=/test
func (s *Service) TestGrp_Test(ctx context.Context, qp *testgrp.QueryParams) (*testgrp.Status, error) {
	return s.testGrp.Test(ctx, qp)
}

//encore:api auth method=GET path=/testauth/:userID tag:auth_user
func (s *Service) TestGrp_TestAuth(ctx context.Context, userID string, qp *testgrp.QueryParams) (*testgrp.Status, error) {
	return s.testGrp.TestAuth(ctx, qp)
}
