package initializer

import (
	"context"
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

type initHandler struct {
	initUseCase domain.InitUseCase
}

func (i *initHandler) Initialize(ctx context.Context) error {
	err := i.initUseCase.ClearRedisData(ctx)
	if err != nil {
		return err
	}
	return i.initUseCase.RecoverRecords(ctx)
}

func NewInitHandler(injector *do.Injector) (domain.InitHandler, error) {
	return &initHandler{do.MustInvoke[domain.InitUseCase](injector)}, nil
}
