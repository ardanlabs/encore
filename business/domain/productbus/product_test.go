package productbus_test

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"testing"
	"time"

	"encore.dev/et"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/sdk/dbtest"
	"github.com/ardanlabs/encore/business/sdk/page"
	"github.com/ardanlabs/encore/business/sdk/unitest"
	"github.com/google/go-cmp/cmp"
)

func Test_Product(t *testing.T) {
	t.Parallel()

	edb, err := et.NewTestDatabase(context.Background(), "app")
	if err != nil {
		t.Fatalf("Creating new database: %s", err)
	}

	db := dbtest.NewDatabase(t, edb)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		db.Teardown()
	}()

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unitest.Run(t, query(db.BusDomain, sd), "query")
	unitest.Run(t, create(db.BusDomain, sd), "create")
	unitest.Run(t, update(db.BusDomain, sd), "update")
	unitest.Run(t, delete(db.BusDomain, sd), "delete")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := userbus.TestGenerateSeedUsers(ctx, 1, userbus.Roles.User, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := unitest.User{
		User:     usrs[0],
		Products: prds,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestGenerateSeedUsers(ctx, 1, userbus.Roles.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := unitest.User{
		User:     usrs[0],
		Products: prds,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Admins: []unitest.User{tu2},
		Users:  []unitest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	prds := make([]productbus.Product, 0, len(sd.Admins[0].Products)+len(sd.Users[0].Products))
	prds = append(prds, sd.Admins[0].Products...)
	prds = append(prds, sd.Users[0].Products...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: prds,
			ExcFunc: func(ctx context.Context) any {
				filter := productbus.QueryFilter{
					Name: dbtest.ProductNamePointer("Name"),
				}

				resp, err := busDomain.Product.Query(ctx, filter, productbus.DefaultOrderBy, page.MustParse("1", "10"))
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
				resp, err := busDomain.Product.QueryByID(ctx, sd.Users[0].Products[0].ID)
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

func create(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: productbus.Product{
				UserID:   sd.Users[0].ID,
				Name:     productbus.Names.MustParse("Guitar"),
				Cost:     10.34,
				Quantity: 10,
			},
			ExcFunc: func(ctx context.Context) any {
				np := productbus.NewProduct{
					UserID:   sd.Users[0].ID,
					Name:     productbus.Names.MustParse("Guitar"),
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := busDomain.Product.Create(ctx, np)
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

func update(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: productbus.Product{
				ID:          sd.Users[0].Products[0].ID,
				UserID:      sd.Users[0].ID,
				Name:        productbus.Names.MustParse("Guitar"),
				Cost:        10.34,
				Quantity:    10,
				DateCreated: sd.Users[0].Products[0].DateCreated,
				DateUpdated: sd.Users[0].Products[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				up := productbus.UpdateProduct{
					Name:     dbtest.ProductNamePointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := busDomain.Product.Update(ctx, sd.Users[0].Products[0], up)
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

func delete(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name:    "user",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Product.Delete(ctx, sd.Users[0].Products[1]); err != nil {
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
				if err := busDomain.Product.Delete(ctx, sd.Admins[0].Products[1]); err != nil {
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
