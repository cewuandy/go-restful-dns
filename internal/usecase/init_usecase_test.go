package usecase

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

type initUseCaseTestSuite struct {
	suite.Suite

	usecase domain.InitUseCase

	redisRepo  *mocks.RedisRepo
	recordRepo *mocks.RecordRepo
}

func TestInitUseCase(t *testing.T) {
	suite.Run(t, &initUseCaseTestSuite{})
}

func (t *initUseCaseTestSuite) SetupSuite() {
	injector := do.New()

	t.redisRepo = &mocks.RedisRepo{}
	t.recordRepo = &mocks.RecordRepo{}
	do.ProvideValue[domain.RedisRepo](injector, t.redisRepo)
	do.ProvideValue[domain.RecordRepo](injector, t.recordRepo)
	t.usecase, _ = NewInitUseCase(injector)
}

func (t *initUseCaseTestSuite) SetupTest() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
		anyTime    = mock.AnythingOfType("time.Duration")
	)

	t.redisRepo.
		On("HSet", anyContext, anyString, anyString, anyString, anyTime).
		Return(nil)
	t.redisRepo.
		On("HGetAll", anyContext, anyString).
		Return(map[string]string{}, nil).Once()
	t.redisRepo.
		On("HGetAll", anyContext, anyString).
		Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
	t.recordRepo.
		On("List", anyContext).
		Return(
			[]*domain.Record{
				{
					Name:   "test.com.",
					RrType: 1,
					Class:  1,
					Record: "test.com.\t1440\tIN\tA\t1.1.1.1",
				},
			}, nil,
		)
	t.redisRepo.
		On("FlushAll", anyContext).
		Return(nil)
}

func (t *initUseCaseTestSuite) SetupErrorTest() {
	t.redisRepo.ExpectedCalls = nil
	t.recordRepo.ExpectedCalls = nil
	t.SetupTest()
	t.redisRepo.ExpectedCalls = nil

	// t.redisRepo.
	// 	On("HGetAll", anyContext, anyString).
	// 	Return(map[string]string{"Answer-0": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)

}

func (t *initUseCaseTestSuite) TestRecoverRecords() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
		anyTime    = mock.AnythingOfType("time.Duration")
	)

	t.Run(
		"success", func() {
			err := t.usecase.RecoverRecords(context.Background())
			t.Nil(err)
		},
	)

	t.Run(
		"List_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil

			t.recordRepo.
				On("List", anyContext).
				Return(nil, fmt.Errorf("test-error"))

			err := t.usecase.RecoverRecords(context.Background())
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HGetAll_record_existed", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil).Once()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil).Once()

			err := t.usecase.RecoverRecords(context.Background())
			t.Nil(err)
		},
	)

	t.Run(
		"HSet_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil).Once()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.RecoverRecords(context.Background())
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HGetAll_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil).Once()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(nil, fmt.Errorf("test-error"))
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)
			err := t.usecase.RecoverRecords(context.Background())
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"getFakeSOA_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil).Once()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)
			err := t.usecase.RecoverRecords(context.Background())
			t.NotNil(err)
			t.Equal("the A record isn't existed", err.Error())
		},
	)

	t.Run(
		"HSet_fakeAAAA_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil).Once()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil).Once()
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.RecoverRecords(context.Background())
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}

func (t *initUseCaseTestSuite) TestClearRedisData() {
	t.Run(
		"success", func() {
			err := t.usecase.ClearRedisData(context.Background())
			t.Nil(err)
		},
	)
}
