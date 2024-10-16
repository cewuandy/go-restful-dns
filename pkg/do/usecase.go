package do

import (
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/usecase"
)

func ProvideUseCase(injector *do.Injector) {
	do.Provide(injector, usecase.NewDNSUseCase)

	do.Provide(injector, usecase.NewRecordUseCase)
}
