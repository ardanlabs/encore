package vproductdb

import (
	"time"

	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/google/uuid"
)

type product struct {
	ID          uuid.UUID `db:"product_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Cost        float64   `db:"cost"`
	Quantity    int       `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	UserName    string    `db:"user_name"`
}

func toBusProduct(db product) vproductbus.Product {
	prd := vproductbus.Product{
		ID:          db.ID,
		UserID:      db.UserID,
		Name:        db.Name,
		Cost:        db.Cost,
		Quantity:    db.Quantity,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
		UserName:    db.UserName,
	}

	return prd
}

func toBusProducts(dbs []product) []vproductbus.Product {
	prds := make([]vproductbus.Product, len(dbs))

	for i, db := range dbs {
		prds[i] = toBusProduct(db)
	}

	return prds
}
