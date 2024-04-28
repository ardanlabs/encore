package user_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/api/auth"
	"github.com/ardanlabs/encore/business/api/dbtest"
	"github.com/ardanlabs/encore/business/domain/userbus"
)

func insertSeedData(dbTest *dbtest.Test, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := dbTest.BusDomain

	usrs, err := userbus.TestGenerateSeedUsers(ctx, 2, userbus.RoleAdmin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := apitest.User{
		User: dbtest.User{
			User: usrs[0],
		},
		Token: apitest.Token(dbTest, ath, usrs[0].Email.Address),
	}

	tu2 := apitest.User{
		User: dbtest.User{
			User: usrs[1],
		},
		Token: apitest.Token(dbTest, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestGenerateSeedUsers(ctx, 3, userbus.RoleUser, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := apitest.User{
		User: dbtest.User{
			User: usrs[0],
		},
		Token: apitest.Token(dbTest, ath, usrs[0].Email.Address),
	}

	tu4 := apitest.User{
		User: dbtest.User{
			User: usrs[1],
		},
		Token: apitest.Token(dbTest, ath, usrs[1].Email.Address),
	}

	tu5 := apitest.User{
		User: dbtest.User{
			User: usrs[2],
		},
		Token: apitest.Token(dbTest, ath, usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Users:  []apitest.User{tu3, tu4, tu5},
		Admins: []apitest.User{tu1, tu2},
	}

	return sd, nil
}
