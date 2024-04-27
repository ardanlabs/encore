package dbtest

import (
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/userbus"
)

// User represents an app user specified for the test.
type User struct {
	userbus.User
	Products []productbus.Product
	Homes    []homebus.Home
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}
