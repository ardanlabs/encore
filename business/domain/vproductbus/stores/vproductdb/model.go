package vproductdb

import (
	"fmt"
	"time"

	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/userbus"
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

func toBusProduct(db product) (vproductbus.Product, error) {
	userName, err := userbus.ParseName(db.UserName)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse user name: %w", err)
	}

	name, err := productbus.ParseName(db.Name)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse name: %w", err)
	}

	bus := vproductbus.Product{
		ID:          db.ID,
		UserID:      db.UserID,
		Name:        name,
		Cost:        db.Cost,
		Quantity:    db.Quantity,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
		UserName:    userName,
	}

	return bus, nil
}

func toBusProducts(dbPrds []product) ([]vproductbus.Product, error) {
	bus := make([]vproductbus.Product, len(dbPrds))

	for i, dbPrd := range dbPrds {
		var err error
		bus[i], err = toBusProduct(dbPrd)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
