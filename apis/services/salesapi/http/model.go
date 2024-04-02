package salesapi

import (
	homeapp "github.com/ardanlabs/encore/app/core/crud/homeapp"
	"github.com/ardanlabs/encore/app/core/crud/productapp"
	"github.com/ardanlabs/encore/app/core/crud/tranapp"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/app/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type crudApp struct {
	home    *homeapp.Core
	product *productapp.Core
	tran    *tranapp.Core
	user    *userapp.Core
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
