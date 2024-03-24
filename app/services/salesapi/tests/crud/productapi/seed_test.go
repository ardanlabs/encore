package product_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/encore/business/api/order"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	api := dbTest.Core.Crud

	usrs, err := api.User.Query(context.Background(), user.QueryFilter{}, order.By{Field: user.OrderByName, Direction: order.ASC}, 1, 2)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// -------------------------------------------------------------------------

	tu1 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.TokenV1(usrs[0].Email.Address, "gophers"),
	}

	tu2 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.TokenV1(usrs[1].Email.Address, "gophers"),
	}

	// -------------------------------------------------------------------------

	usrs, err = user.TestGenerateSeedUsers(1, user.RoleUser, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := product.TestGenerateSeedProducts(2, api.Product, usrs[0].ID)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu3 := dbtest.User{
		User:     usrs[0],
		Token:    dbTest.TokenV1(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
		Products: prds,
	}

	// -------------------------------------------------------------------------

	usrs, err = user.TestGenerateSeedUsers(1, user.RoleAdmin, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = product.TestGenerateSeedProducts(2, api.Product, usrs[0].ID)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu4 := dbtest.User{
		User:     usrs[0],
		Token:    dbTest.TokenV1(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
		Products: prds,
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Admins: []dbtest.User{tu1, tu4},
		Users:  []dbtest.User{tu2, tu3},
	}

	return sd, nil
}
