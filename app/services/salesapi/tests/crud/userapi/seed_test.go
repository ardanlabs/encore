package user_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	ctx := context.Background()
	api := dbTest.Core.Crud

	usrs, err := user.TestGenerateSeedUsers(ctx, 2, user.RoleAdmin, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.Token(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
	}

	tu2 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.Token(usrs[1].Email.Address, fmt.Sprintf("Password%s", usrs[1].Name[4:])),
	}

	// -------------------------------------------------------------------------

	usrs, err = user.TestGenerateSeedUsers(ctx, 2, user.RoleUser, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.Token(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
	}

	tu4 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.Token(usrs[1].Email.Address, fmt.Sprintf("Password%s", usrs[1].Name[4:])),
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Users:  []dbtest.User{tu3, tu4},
		Admins: []dbtest.User{tu1, tu2},
	}

	return sd, nil
}
