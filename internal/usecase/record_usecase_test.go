package usecase

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/domain/mocks"
)

type recordUseCaseTestSuite struct {
	suite.Suite

	usecase domain.RecordUseCase

	recordRepo *mocks.RecordRepo
	redisRepo  *mocks.RedisRepo
}

func TestRecordUseCase(t *testing.T) {
	suite.Run(t, &recordUseCaseTestSuite{})
}

func (t *recordUseCaseTestSuite) SetupSuite() {
	injector := do.New()
	t.recordRepo = &mocks.RecordRepo{}
	t.redisRepo = &mocks.RedisRepo{}
	do.ProvideValue[domain.RecordRepo](injector, t.recordRepo)
	do.ProvideValue[domain.RedisRepo](injector, t.redisRepo)

	t.usecase, _ = NewRecordUseCase(injector)
}

func (t *recordUseCaseTestSuite) SetupTest() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRecord  = mock.AnythingOfType("*domain.Record")
		anyString  = mock.AnythingOfType("string")
		anyUint16  = mock.AnythingOfType("uint16")
		anyTime    = mock.AnythingOfType("time.Duration")
	)

	t.recordRepo.ExpectedCalls = nil
	t.redisRepo.ExpectedCalls = nil

	t.recordRepo.
		On("Create", anyContext, anyRecord).
		Return(nil)
	t.recordRepo.
		On("Get", anyContext, anyString, anyUint16, anyUint16).
		Return(nil, nil)
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
	t.recordRepo.
		On("Update", anyContext, anyRecord).
		Return(nil)
	t.recordRepo.
		On("Delete", anyContext, anyString, anyUint16, anyUint16).
		Return(nil)
	t.redisRepo.
		On("HSet", anyContext, anyString, anyString, anyString, anyTime).
		Return(nil)
	t.redisRepo.
		On("HGetAll", anyContext, anyString).
		Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
	t.redisRepo.
		On("HDel", anyContext, anyString).
		Return(nil)
}

func (t *recordUseCaseTestSuite) TestCreateRecord() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRecord  = mock.AnythingOfType("*domain.Record")
		anyString  = mock.AnythingOfType("string")
		anyUint16  = mock.AnythingOfType("uint16")
		anyTime    = mock.AnythingOfType("time.Duration")
	)

	a := dns.A{
		Hdr: dns.RR_Header{
			Name:   "test.com.",
			Rrtype: 1,
			Class:  1,
			Ttl:    1440,
		},
		A: net.ParseIP("1.1.1.1"),
	}

	aaaa := dns.AAAA{
		Hdr: dns.RR_Header{
			Name:   "test.com.",
			Rrtype: 28,
			Class:  1,
			Ttl:    1440,
		},
		AAAA: net.ParseIP("b152:e410:14fb:0209:73b2:c3be:c3ad:c0a6"),
	}

	rrAAAA, _ := dns.NewRR(aaaa.String())
	rrA, _ := dns.NewRR(a.String())

	t.Run(
		"success_TypeA", func() {
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.Nil(err)
		},
	)

	t.Run(
		"success_TypeAAAA", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil)
			t.redisRepo.
				On("HDel", anyContext, anyString).
				Return(nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)

			err := t.usecase.CreateRecord(context.Background(), rrAAAA)
			t.Nil(err)
		},
	)

	t.Run(
		"success_TypeAAAA", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-0": ""}, nil)
			t.redisRepo.
				On("HDel", anyContext, anyString).
				Return(nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)

			err := t.usecase.CreateRecord(context.Background(), rrAAAA)
			t.Nil(err)
		},
	)

	t.Run(
		"HDel_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Ns-0": ""}, nil)
			t.redisRepo.
				On("HDel", anyContext, anyString).
				Return(fmt.Errorf("test-error"))
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)

			err := t.usecase.CreateRecord(context.Background(), rrAAAA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"record_existed_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Get", anyContext, anyString, anyUint16, anyUint16).
				Return(
					&domain.Record{
						Name:   "test.com.",
						RrType: 1,
						Class:  1,
						Record: "test.com.\t1440\tIN\tA\t1.1.1.1",
					}, nil,
				)
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Contains(err.Error(), "the record is already existed.")
		},
	)

	t.Run(
		"Get_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Get", anyContext, anyString, anyUint16, anyUint16).
				Return(nil, fmt.Errorf("test-error"))
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"Create_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Get", anyContext, anyString, anyUint16, anyUint16).
				Return(nil, nil)
			t.recordRepo.
				On("Create", anyContext, anyRecord).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HSet_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HGetAll_second_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(nil, fmt.Errorf("test-error"))
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HGetAll_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil)
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Contains(err.Error(), "the A record isn't existed")
		},
	)

	t.Run(
		"HSet_second_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Answer-1": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(nil).Once()
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.CreateRecord(context.Background(), rrA)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}

func (t *recordUseCaseTestSuite) TestGetRecord() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
		anyUint16  = mock.AnythingOfType("uint16")
	)

	q := domain.Question{
		Name:   "test.com.",
		Qtype:  domain.TypeA,
		Qclass: domain.ClassINET,
	}

	t.SetupTest()
	t.recordRepo.ExpectedCalls = nil
	t.recordRepo.
		On("Get", anyContext, anyString, anyUint16, anyUint16).
		Return(
			&domain.Record{
				Name:   "test.com.",
				RrType: 1,
				Class:  1,
				Record: "test.com.\t1440\tIN\tA\t1.1.1.1",
			}, nil,
		)

	t.Run(
		"success", func() {
			rr, err := t.usecase.GetRecord(context.Background(), q)
			t.NotNil(rr)
			t.Equal("test.com.", rr.Header().Name)
			t.Equal(uint16(1), rr.Header().Rrtype)
			t.Equal(uint16(1), rr.Header().Class)
			t.Nil(err)
		},
	)

	t.Run(
		"Get_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Get", anyContext, anyString, anyUint16, anyUint16).
				Return(nil, fmt.Errorf("test-error"))
			rr, err := t.usecase.GetRecord(context.Background(), q)
			t.Nil(rr)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"NewRR_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Get", anyContext, anyString, anyUint16, anyUint16).
				Return(&domain.Record{Record: "test-error"}, nil)
			rr, err := t.usecase.GetRecord(context.Background(), q)
			t.Nil(rr)
			t.NotNil(err)
			t.Equal("dns: not a TTL: \"test-error\" at line: 1:10", err.Error())
		},
	)
}

func (t *recordUseCaseTestSuite) TestListRecords() {
	var anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })

	t.Run(
		"success", func() {
			rrs, err := t.usecase.ListRecords(context.Background())
			t.NotNil(rrs)
			t.Equal("test.com.", rrs[0].Header().Name)
			t.Equal(uint16(1), rrs[0].Header().Rrtype)
			t.Equal(uint16(1), rrs[0].Header().Class)
			t.Nil(err)
		},
	)

	t.Run(
		"Get_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("List", anyContext).
				Return(nil, fmt.Errorf("test-error"))
			rr, err := t.usecase.ListRecords(context.Background())
			t.Nil(rr)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"NewRR_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("List", anyContext).
				Return([]*domain.Record{{Record: "test-error"}}, nil)
			rr, err := t.usecase.ListRecords(context.Background())
			t.Nil(rr)
			t.NotNil(err)
			t.Equal("dns: not a TTL: \"test-error\" at line: 1:10", err.Error())
		},
	)
}

func (t *recordUseCaseTestSuite) TestUpdateRecord() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRecord  = mock.AnythingOfType("*domain.Record")
		anyString  = mock.AnythingOfType("string")
	)

	a := dns.A{
		Hdr: dns.RR_Header{
			Name:   "test.com.",
			Rrtype: 1,
			Class:  1,
			Ttl:    1440,
		},
		A: net.ParseIP("1.1.1.1"),
	}

	rr, _ := dns.NewRR(a.String())

	t.Run(
		"success", func() {
			err := t.usecase.UpdateRecord(context.Background(), rr)
			t.Nil(err)
		},
	)

	t.Run(
		"Get_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Update", anyContext, anyRecord).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.UpdateRecord(context.Background(), rr)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HDel_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HDel", anyContext, anyString).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.UpdateRecord(context.Background(), rr)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}

func (t *recordUseCaseTestSuite) TestDeleteRecord() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
		anyUint16  = mock.AnythingOfType("uint16")
	)

	q := domain.Question{
		Name:   "test.com.",
		Qtype:  domain.TypeA,
		Qclass: domain.ClassINET,
	}

	t.Run(
		"success", func() {
			err := t.usecase.DeleteRecord(context.Background(), q)
			t.Nil(err)
		},
	)

	t.Run(
		"Get_error", func() {
			t.SetupTest()
			t.recordRepo.ExpectedCalls = nil
			t.recordRepo.
				On("Delete", anyContext, anyString, anyUint16, anyUint16).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.DeleteRecord(context.Background(), q)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)

	t.Run(
		"HDel_error", func() {
			t.SetupTest()
			t.redisRepo.ExpectedCalls = nil
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Ns-0": ""}, nil)
			t.redisRepo.
				On("HDel", anyContext, anyString).
				Return(fmt.Errorf("test-error"))
			err := t.usecase.DeleteRecord(context.Background(), q)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}
