// Package apitest contains supporting code for running app layer tests.
package apitest

import (
	"context"
	"net/mail"
	"testing"
	"time"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/api/auth"
	"github.com/ardanlabs/encore/business/api/dbtest"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/golang-jwt/jwt/v4"
)

// User extends the dbtest user for app test support.
type User struct {
	dbtest.User
	Token string
}

// SeedData represents users for app tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an app test.
type Table struct {
	Name    string
	Token   string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}

// AuthParams rerpesents the parameters provided to the auth handler.
type AuthParams struct {
	Authorization string `header:"Authorization"`
}

// AuthHandler defines a function that can be called to handle authentication.
type AuthHandler func(ctx context.Context, ap *AuthParams) (eauth.UID, *auth.Claims, error)

// =============================================================================

// AppTest contains functions for executing an app test.
type AppTest struct {
	DBTest  *dbtest.Test
	Auth    *auth.Auth
	handler AuthHandler
}

// New constructs an app test for use.
func New(dbTest *dbtest.Test, ath *auth.Auth, handler AuthHandler) *AppTest {
	return &AppTest{
		DBTest:  dbTest,
		Auth:    ath,
		handler: handler,
	}
}

// Run performs the actual test logic based on the table data.
func (at *AppTest) Run(t *testing.T, table []Table, testName string) {
	log := func(diff string, got any, exp any) {
		t.Log("DIFF")
		t.Logf("%s", diff)
		t.Log("GOT")
		t.Logf("%#v", got)
		t.Log("EXP")
		t.Logf("%#v", exp)
		t.Fatalf("Should get the expected response")
	}

	for _, tt := range table {
		f := func(t *testing.T) {
			ctx := context.Background()

			t.Log("Calling authHandler")
			ctx, err := at.authHandler(ctx, tt.Token)
			if err != nil {
				diff := tt.CmpFunc(err, tt.ExpResp)
				if diff != "" {
					log(diff, err, tt.ExpResp)
				}
				return
			}

			t.Log("Calling excFunc")
			got := tt.ExcFunc(ctx)

			diff := tt.CmpFunc(got, tt.ExpResp)
			if diff != "" {
				log(diff, got, tt.ExpResp)
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

func (at *AppTest) authHandler(ctx context.Context, token string) (context.Context, error) {
	uid, claims, err := at.handler(ctx, &AuthParams{
		Authorization: "Bearer " + token,
	})

	if err != nil {
		return ctx, err
	}

	return eauth.WithContext(ctx, uid, claims), nil
}

// =============================================================================

// CmpAppErrors compares two encore error values. If they are not equal, the
// reason is returned.
func CmpAppErrors(got any, exp any) string {
	expResp := exp.(*eerrs.Error)

	gotResp, exists := got.(*eerrs.Error)
	if !exists {
		return "no error occurred"
	}

	if gotResp.Code != expResp.Code {
		return "code does not match"
	}

	if gotResp.Message != expResp.Message {
		return "message does not match"
	}

	return ""
}

// Token generates an authenticated token for a user.
func Token(dbTest *dbtest.Test, ath *auth.Auth, email string) string {
	addr, _ := mail.ParseAddress(email)

	store := userdb.NewStore(dbTest.Log, dbTest.DB)
	dbUsr, err := store.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	token, err := ath.GenerateToken(kid, claims)
	if err != nil {
		return ""
	}

	return token
}

// =============================================================================

// KeyStore represents a mock keystore with a hardcoded key.
type KeyStore struct{}

// PrivateKey implements the auth interface.
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	return privateKeyPEM, nil
}

// PublicKey implements the auth interface.
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	return publicKeyPEM, nil
}

const (
	kid = "s4sKIjD9kIRjxs2tulPqGLdxSfgPErRN1Mu3Hd9k9NQ"

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEpQIBAAKCAQEAvMAHb0IoLvoYuW2kA+LTmnk+hfnBq1eYIh4CT/rMPCxgtzjq
U0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiMQl/bxAhgLkGrUhgpZ10Bt6GzSEgw
QNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNcLxPUmfjbbH64g+seg28pbgQPwKFK
tF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R42LH9c75JsRzz4FeN+VzwMAL6yQn
ZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLRBzQGN00Z7TuBRIfDJWU21qQN4Cq7
keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMGvwIDAQABAoIBAFQmQKpHkmavNYql
6POaksBRwaA1YzSijr7XJizGIXvKRSwqgb2zdnuTSgpspAx09Dr/aDdy7rZ0DAJt
fk2mInINDottOIQm3txwzTS58GQQAT/+fxTKWJMqwPfxYFPWqbbU76T8kXYna0Gs
OcK36GdMrgIfQqQyMs0Na8MpMg1LmkAxuqnFCXS/NMyKl9jInaaTS+Kz+BSzUMGQ
zebfLFsf2N7sLZuimt9zlRG30JJTfBlB04xsYMo734usA2ITe8U0XqG6Og0qc6ev
6lsoM8hpvEUsQLcjQQ5up7xx3S2stZJ8o0X8GEX5qUMaomil8mZ7X5xOlEqf7p+v
lXQ46cECgYEA2lbZQON6l3ZV9PCn9j1rEGaXio3SrAdTyWK3D1HF+/lEjClhMkfC
XrECOZYj+fiI9n+YpSog+tTDF7FTLf7VP21d2gnhQN6KAXUnLIypzXxodcC6h+8M
ZGJh/EydLvC7nPNoaXx96bohxzS8hrOlOlkCbr+8gPYKf8qkbe7HyxECgYEA3U6e
x9g4FfTvI5MGrhp2BIzoRSn7HlNQzjJ71iMHmM2kBm7TsER8Co1PmPDrP8K/UyGU
Q25usTsPSrHtKQEV6EsWKaP/6p2Q82sDkT9bZlV+OjRvOfpdO5rP6Q95vUmMGWJ/
S6oimbXXL8p3gDafw3vC1PCAhoaxMnGyKuZwlM8CgYEAixT1sXr2dZMg8DV4mMfI
8pqXf+AVyhWkzsz+FVkeyAKiIrKdQp0peI5C/5HfevVRscvX3aY3efCcEfSYKt2A
07WEKkdO4LahrIoHGT7FT6snE5NgfwTMnQl6p2/aVLNun20CHuf5gTBbIf069odr
Af7/KLMkjfWs/HiGQ6zuQjECgYEAv+DIvlDz3+Wr6dYyNoXuyWc6g60wc0ydhQo0
YKeikJPLoWA53lyih6uZ1escrP23UOaOXCDFjJi+W28FR0YProZbwuLUoqDW6pZg
U3DxWDrL5L9NqKEwcNt7ZIDsdnfsJp5F7F6o/UiyOFd9YQb7YkxN0r5rUTg7Lpdx
eMyv0/UCgYEAhX9MPzmTO4+N8naGFof1o8YP97pZj0HkEvM0hTaeAQFKJiwX5ijQ
xumKGh//G0AYsjqP02ItzOm2mWnbI3FrNlKmGFvR6VxIZMOyXvpLofHucjJ5SWli
eYjPklKcXaMftt1FVO4n+EKj1k1+Tv14nytq/J5WN+r4FBlNEYj/6vg=
-----END PRIVATE KEY-----
`
	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvMAHb0IoLvoYuW2kA+LT
mnk+hfnBq1eYIh4CT/rMPCxgtzjqU0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiM
Ql/bxAhgLkGrUhgpZ10Bt6GzSEgwQNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNc
LxPUmfjbbH64g+seg28pbgQPwKFKtF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R
42LH9c75JsRzz4FeN+VzwMAL6yQnZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLR
BzQGN00Z7TuBRIfDJWU21qQN4Cq7keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMG
vwIDAQAB
-----END PUBLIC KEY-----`
)
