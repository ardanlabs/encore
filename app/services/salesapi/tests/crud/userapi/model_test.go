package user_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

func toAppUser(usr user.User) userapi.AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return userapi.AppUser{
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

func toAppUsers(users []user.User) []userapi.AppUser {
	items := make([]userapi.AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(usr user.User) *userapi.AppUser {
	appUsr := toAppUser(usr)
	return &appUsr
}
