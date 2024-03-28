package user_test

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"runtime/debug"
	"sort"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
)

var url string

func TestMain(m *testing.M) {
	et.EnableServiceInstanceIsolation()

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	url, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}

	defer func() {
		err = dbtest.StopDB()
	}()

	return m.Run(), nil
}

// =============================================================================

func Test_User(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_User")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	sd, err := insertSeedData(dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	var app dbtest.UnitTest
	app.Test(t, query(dbTest, sd), "query")
	app.Test(t, create(dbTest), "create")
	app.Test(t, update(dbTest, sd), "update")
	app.Test(t, delete(dbTest, sd), "delete")
}

// =============================================================================

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	ctx := context.Background()
	api := dbTest.Core.Crud

	usrs, err := user.TestGenerateSeedUsers(ctx, 2, user.RoleAdmin, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.User{
		User: usrs[0],
	}

	tu2 := dbtest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	usrs, err = user.TestGenerateSeedUsers(ctx, 2, user.RoleUser, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := dbtest.User{
		User: usrs[0],
	}

	tu4 := dbtest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Users:  []dbtest.User{tu3, tu4},
		Admins: []dbtest.User{tu1, tu2},
	}

	return sd, nil
}

// =============================================================================

func query(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	usrs := make([]user.User, 0, len(sd.Admins)+len(sd.Users))

	for _, adm := range sd.Admins {
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.Users {
		usrs = append(usrs, usr.User)
	}

	sort.Slice(usrs, func(i, j int) bool {
		return usrs[i].ID.String() <= usrs[j].ID.String()
	})

	table := []dbtest.UnitTable{
		{
			Name:    "all",
			ExpResp: usrs,
			ExcFunc: func(ctx context.Context) any {
				filter := user.QueryFilter{
					Name: dbtest.StringPointer("Name"),
				}

				resp, err := dbt.Core.Crud.User.Query(ctx, filter, user.DefaultOrderBy, 1, 10)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Users[0].User,
			ExcFunc: func(ctx context.Context) any {
				resp, err := dbt.Core.Crud.User.QueryByID(ctx, sd.Users[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create(dbt *dbtest.Test) []dbtest.UnitTable {
	email, _ := mail.ParseAddress("bill@ardanlabs.com")

	table := []dbtest.UnitTable{
		{
			Name: "basic",
			ExpResp: user.User{
				Name:       "Bill Kennedy",
				Email:      *email,
				Roles:      []user.Role{user.RoleAdmin},
				Department: "IT",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				nu := user.NewUser{
					Name:            "Bill Kennedy",
					Email:           *email,
					Roles:           []user.Role{user.RoleAdmin},
					Department:      "IT",
					Password:        "123",
					PasswordConfirm: "123",
				}

				resp, err := dbt.Core.Crud.User.Create(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(user.User)
				if !exists {
					return "error occurred"
				}

				if err := bcrypt.CompareHashAndPassword(gotResp.PasswordHash, []byte("123")); err != nil {
					return err.Error()
				}

				expResp := exp.(user.User)

				expResp.ID = gotResp.ID
				expResp.PasswordHash = gotResp.PasswordHash
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	email, _ := mail.ParseAddress("jack@ardanlabs.com")

	table := []dbtest.UnitTable{
		{
			Name: "basic",
			ExpResp: user.User{
				ID:          sd.Users[0].ID,
				Name:        "Jack Kennedy",
				Email:       *email,
				Roles:       []user.Role{user.RoleAdmin},
				Department:  "IT",
				Enabled:     true,
				DateCreated: sd.Users[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uu := user.UpdateUser{
					Name:            dbtest.StringPointer("Jack Kennedy"),
					Email:           email,
					Roles:           []user.Role{user.RoleAdmin},
					Department:      dbtest.StringPointer("IT"),
					Password:        dbtest.StringPointer("1234"),
					PasswordConfirm: dbtest.StringPointer("1234"),
				}

				resp, err := dbt.Core.Crud.User.Update(ctx, sd.Users[0].User, uu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(user.User)
				if !exists {
					return "error occurred"
				}

				if err := bcrypt.CompareHashAndPassword(gotResp.PasswordHash, []byte("1234")); err != nil {
					return err.Error()
				}

				expResp := exp.(user.User)

				expResp.PasswordHash = gotResp.PasswordHash
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func delete(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	table := []dbtest.UnitTable{
		{
			Name:    "user",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := dbt.Core.Crud.User.Delete(ctx, sd.Users[1].User); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "admin",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := dbt.Core.Crud.User.Delete(ctx, sd.Admins[1].User); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
