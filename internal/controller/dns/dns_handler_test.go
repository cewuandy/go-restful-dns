package dns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
	"time"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/domain/mocks"
)

type dnsHandlerTestSuite struct {
	suite.Suite

	handler dns.Handler

	dnsServer *dns.Server
	dnsClient *dns.Client

	dnsUseCase *mocks.DNSUseCase

	question dns.Question
}

func TestDNSHandler(t *testing.T) {
	suite.Run(t, &dnsHandlerTestSuite{})
}

func (t *dnsHandlerTestSuite) SetupSuite() {
	var err error

	injector := do.New()

	t.question = dns.Question{
		Name:   "test.com.",
		Qtype:  1,
		Qclass: 1,
	}

	t.dnsUseCase = &mocks.DNSUseCase{}
	do.ProvideValue[domain.DNSUseCase](injector, t.dnsUseCase)

	t.handler, err = NewDNSHandler(injector)
	t.Nil(err)

	t.dnsClient = &dns.Client{Net: "udp", DialTimeout: time.Second}
	t.dnsServer = &dns.Server{
		Addr:    "127.0.0.1:53",
		Net:     "udp",
		Handler: t.handler,
	}
	started := new(bool)
	*started = false
	t.dnsServer.NotifyStartedFunc = func() {
		*started = true
		fmt.Printf("Listen and serve at localhost:53 for DNS query\n")
	}

	go func() {
		err = t.dnsServer.ListenAndServe()
		t.Nil(err)
	}()

	for !*started {
		time.Sleep(time.Second)
	}
}

func (t *dnsHandlerTestSuite) SetupTest() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyMsg     = mock.AnythingOfType("*dns.Msg")
	)
	t.dnsUseCase.
		On("QueryRedisCache", anyContext, anyMsg).
		Return(&dns.Msg{Question: []dns.Question{t.question}}, nil)
}

func (t *dnsHandlerTestSuite) TearDownSuite() {
	err := t.dnsServer.Shutdown()
	t.Nil(err)
}

func (t *dnsHandlerTestSuite) TestServeDNS() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyMsg     = mock.AnythingOfType("*dns.Msg")
		req        = &dns.Msg{Question: []dns.Question{t.question}}
		a          = dns.A{
			Hdr: dns.RR_Header{
				Name:   "test.com.",
				Rrtype: 1,
				Class:  1,
				Ttl:    1440,
			},
			A: net.ParseIP("2.2.2.2"),
		}
	)

	rr, _ := dns.NewRR(a.String())

	t.Run(
		"redis_success", func() {
			t.dnsUseCase.ExpectedCalls = nil
			t.dnsUseCase.
				On("QueryRedisCache", anyContext, anyMsg).
				Return(
					&dns.Msg{
						Question: []dns.Question{t.question},
						Answer:   []dns.RR{rr},
					}, nil,
				)
			resp, _, err := t.dnsClient.Exchange(req, "127.0.0.1:53")
			t.Equal("test.com.", resp.Answer[0].Header().Name)
			t.Equal(uint16(1), resp.Answer[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Answer[0].Header().Class)
			t.Equal(uint32(1440), resp.Answer[0].Header().Ttl)
			t.Contains(resp.Answer[0].String(), "2.2.2.2")
			t.Nil(err)
		},
	)

	t.Run(
		"upstream_success", func() {
			t.dnsUseCase.ExpectedCalls = nil
			t.SetupTest()
			t.dnsUseCase.
				On("QueryUpstream", anyContext, anyMsg).
				Return(
					&dns.Msg{
						Question: []dns.Question{t.question},
						Answer:   []dns.RR{rr},
					}, nil,
				)

			resp, _, err := t.dnsClient.Exchange(req, "127.0.0.1:53")
			t.Equal("test.com.", resp.Answer[0].Header().Name)
			t.Equal(uint16(1), resp.Answer[0].Header().Rrtype)
			t.Equal(uint16(1), resp.Answer[0].Header().Class)
			t.Equal(uint32(1440), resp.Answer[0].Header().Ttl)
			t.Contains(resp.Answer[0].String(), "2.2.2.2")
			t.Nil(err)
		},
	)

	t.Run(
		"redis_error", func() {
			t.dnsUseCase.ExpectedCalls = nil
			t.dnsUseCase.
				On("QueryRedisCache", anyContext, anyMsg).
				Return(nil, fmt.Errorf("test-error"))
			resp, _, err := t.dnsClient.Exchange(req, "127.0.0.1:53")
			t.Nil(resp)
			t.NotNil(err)
			t.Contains(err.Error(), "i/o timeout")
		},
	)

	t.Run(
		"upstream_error", func() {
			t.dnsUseCase.ExpectedCalls = nil
			t.SetupTest()
			t.dnsUseCase.
				On("QueryUpstream", anyContext, anyMsg).
				Return(nil, fmt.Errorf("test-error"))

			resp, _, err := t.dnsClient.Exchange(req, "127.0.0.1:53")
			t.Nil(resp)
			t.NotNil(err)
			t.Contains(err.Error(), "i/o timeout")
		},
	)

}
