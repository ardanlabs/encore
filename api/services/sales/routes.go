package sales

import (
	"context"
	"net/http"

	"encore.dev"
	"github.com/ardanlabs/encore/app/domain/homeapp"
	"github.com/ardanlabs/encore/app/domain/productapp"
	"github.com/ardanlabs/encore/app/domain/tranapp"
	"github.com/ardanlabs/encore/app/domain/userapp"
	"github.com/ardanlabs/encore/app/domain/vproductapp"
	"github.com/ardanlabs/encore/app/sdk/page"
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
//encore:api auth method=POST path=/v1/homes tag:metrics tag:authorize tag:as_user_role
func (s *Service) HomeCreate(ctx context.Context, app homeapp.NewHome) (homeapp.Home, error) {
	return s.homeApp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeUpdate(ctx context.Context, homeID string, app homeapp.UpdateHome) (homeapp.Home, error) {
	return s.homeApp.Update(ctx, homeID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeDelete(ctx context.Context, homeID string) error {
	return s.homeApp.Delete(ctx, homeID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes tag:metrics tag:authorize tag:as_any_role
func (s *Service) HomeQuery(ctx context.Context, qp homeapp.QueryParams) (page.Document[homeapp.Home], error) {
	return s.homeApp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes/:productID tag:metrics tag:authorize_home
func (s *Service) HomeQueryByID(ctx context.Context, productID string) (homeapp.Home, error) {
	return s.homeApp.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/products tag:metrics tag:authorize tag:as_user_role
func (s *Service) ProductCreate(ctx context.Context, app productapp.NewProduct) (productapp.Product, error) {
	return s.productApp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductUpdate(ctx context.Context, productID string, app productapp.UpdateProduct) (productapp.Product, error) {
	return s.productApp.Update(ctx, productID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductDelete(ctx context.Context, productID string) error {
	return s.productApp.Delete(ctx, productID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products tag:metrics tag:authorize tag:as_any_role
func (s *Service) ProductQuery(ctx context.Context, qp productapp.QueryParams) (page.Document[productapp.Product], error) {
	return s.productApp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductQueryByID(ctx context.Context, productID string) (productapp.Product, error) {
	return s.productApp.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/tran tag:transaction tag:metrics tag:authorize tag:as_admin_role
func (s *Service) TranCreate(ctx context.Context, app tranapp.NewTran) (tranapp.Product, error) {
	return s.tranApp.Create(ctx, app)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/users tag:metrics tag:authorize tag:as_admin_role
func (s *Service) UserCreate(ctx context.Context, app userapp.NewUser) (userapp.User, error) {
	return s.userApp.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserUpdate(ctx context.Context, userID string, app userapp.UpdateUser) (userapp.User, error) {
	return s.userApp.Update(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/role/:userID tag:metrics tag:authorize_user tag:as_admin_role
func (s *Service) UserUpdateRole(ctx context.Context, userID string, app userapp.UpdateUserRole) (userapp.User, error) {
	return s.userApp.UpdateRole(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserDelete(ctx context.Context, userID string) error {
	return s.userApp.Delete(ctx, userID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users tag:metrics tag:authorize tag:as_admin_role
func (s *Service) UserQuery(ctx context.Context, qp userapp.QueryParams) (page.Document[userapp.User], error) {
	return s.userApp.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserQueryByID(ctx context.Context, userID string) (userapp.User, error) {
	return s.userApp.QueryByID(ctx, userID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/vproducts tag:metrics tag:authorize tag:as_admin_role
func (s *Service) VProductQuery(ctx context.Context, qp vproductapp.QueryParams) (page.Document[vproductapp.Product], error) {
	return s.vproductApp.Query(ctx, qp)
}
