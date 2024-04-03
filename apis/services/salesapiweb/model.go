package salesapiweb

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

type appCrud struct {
	home    *homeapp.Core
	product *productapp.Core
	tran    *tranapp.Core
	user    *userapp.Core
}

type appView struct {
	product *vproductapp.Core
}

type busCrud struct {
	home    *home.Core
	product *product.Core
	user    *user.Core
}
