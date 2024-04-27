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
	homeApp     *homeapp.Core
	productApp  *productapp.Core
	tranApp     *tranapp.Core
	userApp     *userapp.Core
	vproductApp *vproductapp.Core
}

type busDomain struct {
	delegate   *delegate.Delegate
	homeBus    *homebus.Core
	productBus *productbus.Core
	userBus    *userbus.Core
}
