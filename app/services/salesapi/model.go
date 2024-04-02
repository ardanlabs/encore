package salesapi

import (
	homeapp "github.com/ardanlabs/encore/app/services/salesapi/core/crud/homeapp"
	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/productapp"
	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/tranapp"
	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/userapp"
	"github.com/ardanlabs/encore/app/services/salesapi/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type crudApp struct {
	home    *homeapp.API
	product *productapp.API
	tran    *tranapp.API
	user    *userapp.API
}

type viewApp struct {
	product *vproductapp.Handlers
}

type app struct {
	crud crudApp
	view viewApp
}

type crudBus struct {
	home    *home.Core
	product *product.Core
	user    *user.Core
}

type business struct {
	crud crudBus
}
