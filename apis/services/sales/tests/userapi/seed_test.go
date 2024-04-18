package user_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/domain/userbus"
)

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	ctx := context.Background()
	busDomain := dbTest.BusDomain

	usrs, err := userbus.TestGenerateSeedUsers(ctx, 2, userbus.RoleAdmin, busDomain.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.Token(usrs[0].Email.Address),
	}

	tu2 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.Token(usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestGenerateSeedUsers(ctx, 3, userbus.RoleUser, busDomain.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.Token(usrs[0].Email.Address),
	}

	tu4 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.Token(usrs[1].Email.Address),
	}

	tu5 := dbtest.User{
		User:  usrs[2],
		Token: dbTest.Token(usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Users:  []dbtest.User{tu3, tu4, tu5},
		Admins: []dbtest.User{tu1, tu2},
	}

	return sd, nil
}
