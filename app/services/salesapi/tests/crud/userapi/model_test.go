package user_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/userapp"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

func toAppUser(usr user.User) userapp.User {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return userapp.User{
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

func toAppUsers(users []user.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}
