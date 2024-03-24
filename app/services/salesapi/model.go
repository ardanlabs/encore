package salesapi

import (
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/tranapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/views/vproductapi"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type coreAPI struct {
	home    *homeapi.Handlers
	product *productapi.Handlers
	tran    *tranapi.Handlers
	user    *userapi.Handlers
}

type viewAPI struct {
	product *vproductapi.Handlers
}

type api struct {
	core coreAPI
	view viewAPI
}

type core struct {
	home    *home.Core
	product *product.Core
	user    *user.Core
}
