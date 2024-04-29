package sales

import (
	homeapp "github.com/ardanlabs/encore/app/domain/homeapp"
	productapp "github.com/ardanlabs/encore/app/domain/productapp"
	tranapp "github.com/ardanlabs/encore/app/domain/tranapp"
	userapp "github.com/ardanlabs/encore/app/domain/userapp"
	vproductapp "github.com/ardanlabs/encore/app/domain/vproductapp"
	"github.com/ardanlabs/encore/business/api/delegate"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/userbus"
)

type appDomain struct {
	homeApp     *homeapp.App
	productApp  *productapp.App
	tranApp     *tranapp.App
	userApp     *userapp.App
	vproductApp *vproductapp.App
}

type busDomain struct {
	delegate   *delegate.Delegate
	homeBus    *homebus.Business
	productBus *productbus.Business
	userBus    *userbus.Business
}
