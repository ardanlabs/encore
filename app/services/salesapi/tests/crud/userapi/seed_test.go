package user_test

import (
	"fmt"

	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	usrs, err := user.TestGenerateSeedUsers(2, user.RoleAdmin, dbTest.CoreAPIs.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.TokenV1(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
	}

	tu2 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.TokenV1(usrs[1].Email.Address, fmt.Sprintf("Password%s", usrs[1].Name[4:])),
	}

	// -------------------------------------------------------------------------

	usrs, err = user.TestGenerateSeedUsers(2, user.RoleUser, dbTest.CoreAPIs.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := dbtest.User{
		User:  usrs[0],
		Token: dbTest.TokenV1(usrs[0].Email.Address, fmt.Sprintf("Password%s", usrs[0].Name[4:])),
	}

	tu4 := dbtest.User{
		User:  usrs[1],
		Token: dbTest.TokenV1(usrs[1].Email.Address, fmt.Sprintf("Password%s", usrs[1].Name[4:])),
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Users:  []dbtest.User{tu3, tu4},
		Admins: []dbtest.User{tu1, tu2},
	}

	return sd, nil
}
