package encore

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/trangrp"
)

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/tran tag:metrics tag:authorize_user_only
func (s *Service) TranGrpCreate(ctx context.Context, app trangrp.AppNewTran) (trangrp.AppProduct, error) {
	return s.trnGrp.Create(ctx, app)
}
