package productapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/foundation/validate"
)

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page     int    `query:"page"`
	Rows     int    `query:"rows"`
	OrderBy  string `query:"orderBy"`
	ID       string `query:"product_id"`
	Name     string `query:"name"`
	Cost     string `query:"cost"`
	Quantity string `query:"quantity"`
}

// AppProduct represents information about an individual product.
type AppProduct struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

func toAppProduct(prd product.Product) AppProduct {
	return AppProduct{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []product.Product) []AppProduct {
	items := make([]AppProduct, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}

// AppNewProduct defines the data needed to add a new product.
type AppNewProduct struct {
	Name     string  `json:"name" validate:"required"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
}

func toCoreNewProduct(ctx context.Context, app AppNewProduct) (product.NewProduct, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("getuserid: %w", err)
	}

	prd := product.NewProduct{
		UserID:   userID,
		Name:     app.Name,
		Cost:     app.Cost,
		Quantity: app.Quantity,
	}

	return prd, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewProduct) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(http.StatusBadRequest, "validate: %s", err)
	}

	return nil
}

// AppUpdateProduct defines the data needed to update a product.
type AppUpdateProduct struct {
	Name     *string  `json:"name"`
	Cost     *float64 `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int     `json:"quantity" validate:"omitempty,gte=1"`
}

func toCoreUpdateProduct(app AppUpdateProduct) product.UpdateProduct {
	core := product.UpdateProduct{
		Name:     app.Name,
		Cost:     app.Cost,
		Quantity: app.Quantity,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateProduct) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(http.StatusBadRequest, "validate: %s", err)
	}

	return nil
}
