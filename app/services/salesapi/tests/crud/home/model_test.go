package home_test

import (
	"time"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/homegrp"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type tableData struct {
	name       string
	url        string
	token      string
	method     string
	statusCode int
	model      any
	resp       any
	expResp    any
	excFunc    func()
	cmpFunc    func(x interface{}, y interface{}) string
}

type testUser struct {
	user.User
	token string
	homes []home.Home
}

type seedData struct {
	users  []testUser
	admins []testUser
}

func toPointer(r middleware.Response) *middleware.Response {
	return &r
}

func toAppHome(hme home.Home) homegrp.AppHome {
	return homegrp.AppHome{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type.Name(),
		Address: homegrp.AppAddress{
			Address1: hme.Address.Address1,
			Address2: hme.Address.Address2,
			ZipCode:  hme.Address.ZipCode,
			City:     hme.Address.City,
			State:    hme.Address.State,
			Country:  hme.Address.Country,
		},
		DateCreated: hme.DateCreated.Format(time.RFC3339),
		DateUpdated: hme.DateUpdated.Format(time.RFC3339),
	}
}

func toAppHomes(homes []home.Home) []homegrp.AppHome {
	items := make([]homegrp.AppHome, len(homes))
	for i, hme := range homes {
		items[i] = toAppHome(hme)
	}

	return items
}

func toAppHomePtr(hme home.Home) *homegrp.AppHome {
	appHme := toAppHome(hme)
	return &appHme
}
