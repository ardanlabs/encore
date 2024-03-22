package user_test

import (
	"context"

	"encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/usergrp"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/page"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userQuery200(sd seedData) []tableData {
	usrs := make([]user.User, 0, len(sd.admins)+len(sd.users))
	usrsMap := make(map[uuid.UUID]user.User)

	for _, adm := range sd.admins {
		usrsMap[adm.ID] = adm.User
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.users {
		usrsMap[usr.ID] = usr.User
		usrs = append(usrs, usr.User)
	}

	table := []tableData{
		{
			name:  "query",
			token: sd.admins[0].token,
			expResp: page.Document[usergrp.AppUser]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(usrs),
				Items:       toAppUsers(usrs),
			},
			excFunc: func(ctx context.Context, s *salesapi.Service) any {
				qp := usergrp.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "user_id,ASC",
					Name:    "Name",
				}

				resp, err := s.UserGrpQuery(ctx, qp)
				if err != nil {
					return err
				}

				return resp
			},
			cmpFunc: func(got any, exp any) string {
				if errs, exists := got.(*errs.Error); exists {
					return errs.Message
				}

				gotResp := got.(page.Document[usergrp.AppUser])
				expResp := exp.(page.Document[usergrp.AppUser])

				var found int
				for _, r := range gotResp.Items {
					for _, e := range expResp.Items {
						if e.ID == r.ID {
							found++
							break
						}
					}
				}

				if found != len(usrs) {
					return "number of expected users didn't match"
				}

				return ""
			},
		},
	}

	return table
}

func userQueryByID200(sd seedData) []tableData {
	table := []tableData{
		{
			name:    "basic",
			token:   sd.users[0].token,
			expResp: toAppUserPtr(sd.users[0].User),
			excFunc: func(ctx context.Context, s *salesapi.Service) any {
				resp, err := s.UserGrpQueryByID(ctx, sd.users[0].ID.String())
				if err != nil {
					return err
				}

				return resp
			},
			cmpFunc: func(x any, y any) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
