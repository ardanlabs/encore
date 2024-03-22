package user_test

import (
	"context"
	"time"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/usergrp"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type tableData struct {
	name    string
	token   string
	expResp any
	excFunc func(ctx context.Context) any
	cmpFunc func(x any, y any) string
}

type testUser struct {
	user.User
	token string
}

type seedData struct {
	users  []testUser
	admins []testUser
}

func toPointer(r middleware.Response) *middleware.Response {
	return &r
}

func toAppUser(usr user.User) usergrp.AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return usergrp.AppUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: nil,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []usergrp.AppUser {
	items := make([]usergrp.AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(usr user.User) *usergrp.AppUser {
	appUsr := toAppUser(usr)
	return &appUsr
}
