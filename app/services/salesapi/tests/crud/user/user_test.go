package user_test

import (
	"context"
	"runtime/debug"
	"testing"

	"encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/mid"
)

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

	service, err := salesapi.NewService(dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapi", service)

	app := appTest{
		service: service,
	}

	// -------------------------------------------------------------------------

	sd, err := createUserSeed(dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	app.test(t, userQuery200(sd), "user-query-200")
	app.test(t, userQueryByID200(sd), "user-querybyid-200")

	// app.test(t, userCreate200(sd), "user-create-200")
	// app.test(t, userCreate401(sd), "user-create-401")
	// app.test(t, userCreate400(sd), "user-create-400")

	// app.test(t, userUpdate200(sd), "user-update-200")
	// app.test(t, userUpdate401(sd), "user-update-401")
	// app.test(t, userUpdate400(sd), "user-update-400")

	// app.test(t, userDelete200(sd), "user-delete-200")
	// app.test(t, userDelete401(sd), "user-delete-401")
}

// =============================================================================

type appTest struct {
	service *salesapi.Service
}

func (at *appTest) test(t *testing.T, table []tableData, testName string) {
	log := func(got any, exp any) {
		t.Log("GOT")
		t.Logf("%#v", got)
		t.Log("EXP")
		t.Logf("%#v", exp)
		t.Fatalf("Should get the expected response")
	}

	for _, tt := range table {
		f := func(t *testing.T) {
			t.Log("Calling authHandler")
			ctx, err := at.authHandler(context.Background(), tt.token)
			if err != nil {
				diff := tt.cmpFunc(err, tt.expResp)
				if diff != "" {
					log(err, tt.expResp)
				}
				return
			}

			t.Log("Calling excFunc")
			got := tt.excFunc(ctx)

			diff := tt.cmpFunc(got, tt.expResp)
			if diff != "" {
				log(got, tt.expResp)
			}
		}

		t.Run(testName+"-"+tt.name, f)
	}
}

func (at *appTest) authHandler(ctx context.Context, token string) (context.Context, error) {
	uid, claims, err := at.service.AuthHandler(ctx, &mid.AuthParams{
		Authorization: "Bearer " + token,
	})

	if err != nil {
		return ctx, err
	}

	return auth.WithContext(ctx, uid, claims), nil
}
