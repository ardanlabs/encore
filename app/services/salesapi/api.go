package salesapi

import (
	"context"
	"net/http"

	"encore.dev"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/tranapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/views/vproductapi"
	"github.com/ardanlabs/encore/business/web/page"
)

// Fallback is called for the debug enpoints.
//
//encore:api public raw path=/!fallback
func (s *Service) Fallback(w http.ResponseWriter, r *http.Request) {

	// If this is a web socket call for statsviz and we are in development.
	if r.URL.String() == "/debug/statsviz/ws" && encore.Meta().Environment.Type == encore.EnvDevelopment {

		// In development the r.Host will be host=127.0.0.1:RandPort while the
		// Origin will be origin=127.0.0.1:4000. These need to match.
		r.Header.Set("Origin", "http://"+r.Host)
	}

	s.debug.ServeHTTP(w, r)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/homes tag:metrics tag:authorize_user_only
func (s *Service) HomeCreate(ctx context.Context, app homeapi.AppNewHome) (homeapi.AppHome, error) {
	return s.api.core.home.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeUpdate(ctx context.Context, homeID string, app homeapi.AppUpdateHome) (homeapi.AppHome, error) {
	return s.api.core.home.Update(ctx, homeID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeDelete(ctx context.Context, homeID string) error {
	return s.api.core.home.Delete(ctx, homeID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes tag:metrics tag:authorize_any
func (s *Service) HomeQuery(ctx context.Context, qp homeapi.QueryParams) (page.Document[homeapi.AppHome], error) {
	return s.api.core.home.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes/:productID tag:metrics tag:authorize_home
func (s *Service) HomeQueryByID(ctx context.Context, productID string) (homeapi.AppHome, error) {
	return s.api.core.home.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/products tag:metrics tag:authorize_user_only
func (s *Service) ProductCreate(ctx context.Context, app productapi.AppNewProduct) (productapi.AppProduct, error) {
	return s.api.core.product.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductUpdate(ctx context.Context, productID string, app productapi.AppUpdateProduct) (productapi.AppProduct, error) {
	return s.api.core.product.Update(ctx, productID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductDelete(ctx context.Context, productID string) error {
	return s.api.core.product.Delete(ctx, productID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products tag:metrics tag:authorize_any
func (s *Service) ProductQuery(ctx context.Context, qp productapi.QueryParams) (page.Document[productapi.AppProduct], error) {
	return s.api.core.product.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductQueryByID(ctx context.Context, productID string) (productapi.AppProduct, error) {
	return s.api.core.product.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/tran tag:metrics tag:authorize_user_only
func (s *Service) TranCreate(ctx context.Context, app tranapi.AppNewTran) (tranapi.AppProduct, error) {
	return s.api.core.tran.Create(ctx, app)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid tag:metrics
func (s *Service) UserToken(ctx context.Context, kid string) (userapi.Token, error) {
	return s.api.core.user.Token(ctx, kid)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/users tag:metrics tag:authorize_admin_only
func (s *Service) UserCreate(ctx context.Context, app userapi.AppNewUser) (userapi.AppUser, error) {
	return s.api.core.user.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserUpdate(ctx context.Context, userID string, app userapi.AppUpdateUser) (userapi.AppUser, error) {
	return s.api.core.user.Update(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserDelete(ctx context.Context, userID string) error {
	return s.api.core.user.Delete(ctx, userID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users tag:metrics tag:authorize_admin_only
func (s *Service) UserQuery(ctx context.Context, qp userapi.QueryParams) (page.Document[userapi.AppUser], error) {
	return s.api.core.user.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserQueryByID(ctx context.Context, userID string) (userapi.AppUser, error) {
	return s.api.core.user.QueryByID(ctx, userID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/vproducts tag:metrics tag:authorize_admin_only
func (s *Service) VProductQuery(ctx context.Context, qp vproductapi.QueryParams) (page.Document[vproductapi.AppProduct], error) {
	return s.api.view.product.Query(ctx, qp)
}
