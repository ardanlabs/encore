package user_test

import (
	"time"

	"github.com/ardanlabs/encore/app/domain/userapp"
	"github.com/ardanlabs/encore/business/domain/userbus"
)

func toAppUser(usr userbus.User) userapp.User {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.String()
	}

	return userapp.User{
		ID:           usr.ID.String(),
		Name:         usr.Name.String(),
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: nil,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}
