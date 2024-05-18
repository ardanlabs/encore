package user_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/sdk/auth"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/sdk/dbtest"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, userbus.Roles.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db, ath, usrs[0].Email.Address),
	}

	tu2 := apitest.User{
		User:  usrs[1],
		Token: apitest.Token(db, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 3, userbus.Roles.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db, ath, usrs[0].Email.Address),
	}

	tu4 := apitest.User{
		User:  usrs[1],
		Token: apitest.Token(db, ath, usrs[1].Email.Address),
	}

	tu5 := apitest.User{
		User:  usrs[2],
		Token: apitest.Token(db, ath, usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Users:  []apitest.User{tu3, tu4, tu5},
		Admins: []apitest.User{tu1, tu2},
	}

	return sd, nil
}
