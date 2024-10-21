package initializer

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/domain/mocks"
)

type initHandlerTestSuite struct {
	suite.Suite

	handler domain.InitHandler

	initUseCase *mocks.InitUseCase
}

func TestInitHandler(t *testing.T) {
	suite.Run(t, &initHandlerTestSuite{})
}

func (t *initHandlerTestSuite) SetupSuite() {
	injector := do.New()
	t.initUseCase = &mocks.InitUseCase{}
	do.ProvideValue[domain.InitUseCase](injector, t.initUseCase)
	t.handler, _ = NewInitHandler(injector)
}

func (t *initHandlerTestSuite) SetupTest() {
	var anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })

	t.initUseCase.
		On("ClearRedisData", anyContext).
		Return(nil)
	t.initUseCase.
		On("RecoverRecords", anyContext).
		Return(nil)
}

func (t *initHandlerTestSuite) TestInitialize() {
	var anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })

	t.Run(
		"success", func() {
			err := t.handler.Initialize(context.Background())
			t.Nil(err)
		},
	)

	t.Run(
		"ClearRedisData_error", func() {
			t.SetupTest()
			t.initUseCase.ExpectedCalls = nil
			t.initUseCase.
				On("ClearRedisData", anyContext).
				Return(fmt.Errorf("test-error"))
			err := t.handler.Initialize(context.Background())
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}
