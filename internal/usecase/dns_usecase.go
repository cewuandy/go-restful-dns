package usecase

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/cewuandy/go-restful-dns/internal/domain"

	"github.com/miekg/dns"
	"github.com/samber/do"
)

type dnsUseCase struct {
	redisRepo domain.RedisRepo
	upstreams []string
}

func (d *dnsUseCase) QueryRedisCache(ctx context.Context, req *dns.Msg) (resp *dns.Msg, err error) {
	var rrMap map[string]string

	resp = d.initRespMsg(req, resp)
	q := req.Question[0]
	rrMap, err = d.redisRepo.HGetAll(ctx, q.String())
	if err != nil {
		return nil, err
	}

	for k, v := range rrMap {
		var (
			rr            dns.RR
			answerPattern = regexp.MustCompile(fmt.Sprintf("%s-*", domain.Answer))
			nsPattern     = regexp.MustCompile(fmt.Sprintf("%s-*", domain.Ns))
			extraPattern  = regexp.MustCompile(fmt.Sprintf("%s-*", domain.Extra))
		)

		rr, _ = dns.NewRR(v)
		switch {
		case answerPattern.Match([]byte(k)):
			resp.Answer = append(resp.Answer, rr)
		case nsPattern.Match([]byte(k)):
			resp.Ns = append(resp.Ns, rr)
		case extraPattern.Match([]byte(k)):
			resp.Extra = append(resp.Extra, rr)
		}
	}

	return resp, nil
}

func (d *dnsUseCase) QueryUpstream(ctx context.Context, req *dns.Msg) (resp *dns.Msg, err error) {
	resp = d.initRespMsg(req, resp)
	client := &dns.Client{Net: "udp", DialTimeout: time.Second}

	for _, server := range d.upstreams {
		resp, _, err = client.Exchange(req, server)
		if err != nil || (resp.Answer == nil && resp.Ns == nil && resp.Extra == nil) {
			continue
		}

		if resp.Answer != nil || resp.Ns != nil || resp.Extra != nil {
			err = d.cacheRecord(ctx, resp)
			return resp, err
		}
	}

	return nil, domain.Error{Message: "upstream forwarder error"}
}

func (d *dnsUseCase) initRespMsg(req *dns.Msg, resp *dns.Msg) *dns.Msg {
	resp = new(dns.Msg)
	resp.SetReply(req)
	resp.Compress = false
	return resp
}

func (d *dnsUseCase) cacheRecord(ctx context.Context, resp *dns.Msg) error {
	q := domain.Question{Question: resp.Question[0]}
	for _, t := range domain.ResponseTypeMap {
		value := reflect.ValueOf(*resp).FieldByName(string(t)).Interface()
		rr := value.([]dns.RR)
		for i, rl := range rr {
			field := fmt.Sprintf("%s-%d", t, i)
			ttl := time.Duration(rl.Header().Ttl) * time.Second
			err := d.redisRepo.HSet(ctx, q.String(), field, rl.String(), ttl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewDNSUseCase(injector *do.Injector) (domain.DNSUseCase, error) {
	return &dnsUseCase{
		do.MustInvoke[domain.RedisRepo](injector),
		do.MustInvoke[[]string](injector),
	}, nil
}
