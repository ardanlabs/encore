package home_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/business/core/crud/home"
)

func toAppHome(hme home.Home) homeapi.AppHome {
	return homeapi.AppHome{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type.Name(),
		Address: homeapi.AppAddress{
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

func toAppHomes(homes []home.Home) []homeapi.AppHome {
	items := make([]homeapi.AppHome, len(homes))
	for i, hme := range homes {
		items[i] = toAppHome(hme)
	}

	return items
}

func toAppHomePtr(hme home.Home) *homeapi.AppHome {
	appHme := toAppHome(hme)
	return &appHme
}
