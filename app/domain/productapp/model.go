package productapp

import (
	"context"
	"fmt"
	"time"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/app/sdk/mid"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/foundation/validate"
)

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page     string `query:"page"`
	Rows     string `query:"rows"`
	OrderBy  string `query:"orderBy"`
	ID       string `query:"product_id"`
	Name     string `query:"name"`
	Cost     string `query:"cost"`
	Quantity string `query:"quantity"`
}

// Product represents information about an individual product.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

func toAppProduct(prd productbus.Product) Product {
	return Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []productbus.Product) []Product {
	items := make([]Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}

// NewProduct defines the data needed to add a new product.
type NewProduct struct {
	Name     string  `json:"name" validate:"required"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
}

func toBusNewProduct(ctx context.Context, app NewProduct) (productbus.NewProduct, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return productbus.NewProduct{}, fmt.Errorf("getuserid: %w", err)
	}

	name, err := productbus.Names.Parse(app.Name)
	if err != nil {
		return productbus.NewProduct{}, fmt.Errorf("parse name: %w", err)
	}

	prd := productbus.NewProduct{
		UserID:   userID,
		Name:     name,
		Cost:     app.Cost,
		Quantity: app.Quantity,
	}

	return prd, nil
}

// Validate checks the data in the model is considered clean.
func (app NewProduct) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(eerrs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}

// UpdateProduct defines the data needed to update a product.
type UpdateProduct struct {
	Name     *string  `json:"name"`
	Cost     *float64 `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int     `json:"quantity" validate:"omitempty,gte=1"`
}

func toBusUpdateProduct(app UpdateProduct) (productbus.UpdateProduct, error) {
	var name *productbus.Name
	if app.Name != nil {
		nm, err := productbus.Names.Parse(*app.Name)
		if err != nil {
			return productbus.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		name = &nm
	}

	bus := productbus.UpdateProduct{
		Name:     name,
		Cost:     app.Cost,
		Quantity: app.Quantity,
	}

	return bus, nil
}

// Validate checks the data in the model is considered clean.
func (app UpdateProduct) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(eerrs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}
