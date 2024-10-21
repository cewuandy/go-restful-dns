package usecase

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/domain/mocks"
)

type dnsUseCaseTestSuite struct {
	suite.Suite

	usecase domain.DNSUseCase

	redisRepo *mocks.RedisRepo
}

func TestDnsUseCase(t *testing.T) {
	suite.Run(t, &dnsUseCaseTestSuite{})
}

func (t *dnsUseCaseTestSuite) SetupSuite() {
	injector := do.New()
	t.redisRepo = &mocks.RedisRepo{}
	do.ProvideValue[domain.RedisRepo](injector, t.redisRepo)
	do.ProvideValue[[]string](injector, []string{"1.1.1.1:53"})

	t.usecase, _ = NewDNSUseCase(injector)
}

func (t *dnsUseCaseTestSuite) SetupTest() {
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
		Return(map[string]string{"Answer-0": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
	t.redisRepo.
		On("HDel", anyContext, anyString).
		Return(nil)
}

func (t *dnsUseCaseTestSuite) SetupErrorTest() {
	t.SetupTest()
	t.redisRepo.ExpectedCalls = nil
}

func (t *dnsUseCaseTestSuite) TestQueryRedisCache() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
	)

	req := &dns.Msg{
		Question: []dns.Question{
			{
				Name:   "test.com.",
				Qtype:  1,
				Qclass: 1,
			},
		},
	}

	t.Run(
		"success_answer", func() {
			resp, err := t.usecase.QueryRedisCache(context.Background(), req)
			t.NotNil(resp)
			t.Equal("test.com.", resp.Answer[0].Header().Name)
			t.Equal(uint16(1), resp.Answer[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Answer[0].Header().Class)
			t.Equal(uint32(1440), resp.Answer[0].Header().Ttl)
			t.Nil(err)
		},
	)

	t.Run(
		"success_ns", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Ns-0": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			resp, err := t.usecase.QueryRedisCache(context.Background(), req)
			t.NotNil(resp)
			t.Equal("test.com.", resp.Ns[0].Header().Name)
			t.Equal(uint16(1), resp.Ns[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Ns[0].Header().Class)
			t.Equal(uint32(1440), resp.Ns[0].Header().Ttl)
			t.Nil(err)
		},
	)

	t.Run(
		"success_extra", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(map[string]string{"Extra-0": "test.com.\t1440\tIN\tA\t1.1.1.1"}, nil)
			resp, err := t.usecase.QueryRedisCache(context.Background(), req)
			t.NotNil(resp)
			t.Equal("test.com.", resp.Extra[0].Header().Name)
			t.Equal(uint16(1), resp.Extra[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Extra[0].Header().Class)
			t.Equal(uint32(1440), resp.Extra[0].Header().Ttl)
			t.Nil(err)
		},
	)

	t.Run(
		"HGetAll_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HGetAll", anyContext, anyString).
				Return(nil, fmt.Errorf("test-error"))
			resp, err := t.usecase.QueryRedisCache(context.Background(), req)
			t.Nil(resp)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}

func (t *dnsUseCaseTestSuite) TestQueryUpstream() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyString  = mock.AnythingOfType("string")
		anyTime    = mock.AnythingOfType("time.Duration")
	)
	req := &dns.Msg{
		Question: []dns.Question{
			{
				Name:   "google.com.",
				Qtype:  1,
				Qclass: 1,
			},
		},
	}

	t.Run(
		"success", func() {
			resp, err := t.usecase.QueryUpstream(context.Background(), req)
			t.NotNil(resp)
			t.Equal("google.com.", resp.Answer[0].Header().Name)
			t.Equal(uint16(1), resp.Answer[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Answer[0].Header().Class)
			t.Nil(err)
		},
	)

	t.Run(
		"exchange_error", func() {
			request := &dns.Msg{
				Question: []dns.Question{
					{
						Name:   "notexisted.test.",
						Qtype:  1,
						Qclass: 1,
					},
				},
			}

			resp, err := t.usecase.QueryUpstream(context.Background(), request)
			t.Nil(resp)
			t.NotNil(err)
			t.Contains(err.Error(), "upstream forwarder error")
		},
	)

	t.Run(
		"HSet_error", func() {
			t.SetupErrorTest()
			t.redisRepo.
				On("HSet", anyContext, anyString, anyString, anyString, anyTime).
				Return(fmt.Errorf("test-error"))
			resp, err := t.usecase.QueryUpstream(context.Background(), req)
			t.NotNil(resp)
			t.NotNil(err)
			t.Equal("test-error", err.Error())
		},
	)
}
