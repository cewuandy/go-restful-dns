package do

import (
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/controller/dns"
	"github.com/cewuandy/go-restful-dns/internal/controller/http/middleware"
	v1 "github.com/cewuandy/go-restful-dns/internal/controller/http/v1"
)

func ProvideController(injector *do.Injector) {
	// dns handler
	do.Provide(injector, dns.NewDNSHandler)

	// http handler
	do.Provide(injector, middleware.NewErrorHandler)
	do.Provide(injector, v1.NewRecordHandler)
}
