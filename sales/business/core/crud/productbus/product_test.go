package productbus_test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"testing"
	"time"

	"github.com/ardanlabs/encore/sales/business/core/crud/productbus"
	"github.com/ardanlabs/encore/sales/business/core/crud/userbus"
	"github.com/ardanlabs/encore/sales/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

var url string

func TestMain(m *testing.M) {
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

func Test_Product(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_Product")
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

	dbtest.UnitTest(t, query(dbTest, sd), "query")
	dbtest.UnitTest(t, create(dbTest, sd), "create")
	dbtest.UnitTest(t, update(dbTest, sd), "update")
	dbtest.UnitTest(t, delete(dbTest, sd), "delete")
}

// =============================================================================

func insertSeedData(dbTest *dbtest.Test) (dbtest.SeedData, error) {
	ctx := context.Background()
	api := dbTest.Core.BusCrud

	usrs, err := userbus.TestGenerateSeedUsers(ctx, 1, userbus.RoleUser, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 2, api.Product, usrs[0].ID)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := dbtest.User{
		User:     usrs[0],
		Token:    dbTest.Token(usrs[0].Email.Address),
		Products: prds,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestGenerateSeedUsers(ctx, 1, userbus.RoleAdmin, api.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productbus.TestGenerateSeedProducts(ctx, 2, api.Product, usrs[0].ID)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := dbtest.User{
		User:     usrs[0],
		Token:    dbTest.Token(usrs[0].Email.Address),
		Products: prds,
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Admins: []dbtest.User{tu2},
		Users:  []dbtest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func query(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	prds := make([]productbus.Product, 0, len(sd.Admins[0].Products)+len(sd.Users[0].Products))
	prds = append(prds, sd.Admins[0].Products...)
	prds = append(prds, sd.Users[0].Products...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []dbtest.UnitTable{
		{
			Name:    "all",
			ExpResp: prds,
			ExcFunc: func(ctx context.Context) any {
				filter := productbus.QueryFilter{
					Name: dbtest.StringPointer("Name"),
				}

				resp, err := dbt.Core.BusCrud.Product.Query(ctx, filter, productbus.DefaultOrderBy, 1, 10)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]productbus.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]productbus.Product)

				for i := range gotResp {
					if gotResp[i].DateCreated.Format(time.RFC3339) == expResp[i].DateCreated.Format(time.RFC3339) {
						expResp[i].DateCreated = gotResp[i].DateCreated
					}

					if gotResp[i].DateUpdated.Format(time.RFC3339) == expResp[i].DateUpdated.Format(time.RFC3339) {
						expResp[i].DateUpdated = gotResp[i].DateUpdated
					}
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Users[0].Products[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := dbt.Core.BusCrud.Product.QueryByID(ctx, sd.Users[0].Products[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(productbus.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(productbus.Product)

				if gotResp.DateCreated.Format(time.RFC3339) == expResp.DateCreated.Format(time.RFC3339) {
					expResp.DateCreated = gotResp.DateCreated
				}

				if gotResp.DateUpdated.Format(time.RFC3339) == expResp.DateUpdated.Format(time.RFC3339) {
					expResp.DateUpdated = gotResp.DateUpdated
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	table := []dbtest.UnitTable{
		{
			Name: "basic",
			ExpResp: productbus.Product{
				UserID:   sd.Users[0].ID,
				Name:     "Guitar",
				Cost:     10.34,
				Quantity: 10,
			},
			ExcFunc: func(ctx context.Context) any {
				np := productbus.NewProduct{
					UserID:   sd.Users[0].ID,
					Name:     "Guitar",
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := dbt.Core.BusCrud.Product.Create(ctx, np)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(productbus.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(productbus.Product)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(dbt *dbtest.Test, sd dbtest.SeedData) []dbtest.UnitTable {
	table := []dbtest.UnitTable{
		{
			Name: "basic",
			ExpResp: productbus.Product{
				ID:          sd.Users[0].Products[0].ID,
				UserID:      sd.Users[0].ID,
				Name:        "Guitar",
				Cost:        10.34,
				Quantity:    10,
				DateCreated: sd.Users[0].Products[0].DateCreated,
				DateUpdated: sd.Users[0].Products[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				up := productbus.UpdateProduct{
					Name:     dbtest.StringPointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := dbt.Core.BusCrud.Product.Update(ctx, sd.Users[0].Products[0], up)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(productbus.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(productbus.Product)

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
				if err := dbt.Core.BusCrud.Product.Delete(ctx, sd.Users[0].Products[1]); err != nil {
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
				if err := dbt.Core.BusCrud.Product.Delete(ctx, sd.Admins[0].Products[1]); err != nil {
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
