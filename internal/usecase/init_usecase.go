package usecase

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

type initUseCase struct {
	redisRepo  domain.RedisRepo
	recordRepo domain.RecordRepo
}

func (i *initUseCase) RecoverRecords(ctx context.Context) error {
	records, err := i.recordRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, r := range records {
		q := dns.Question{
			Name:   r.Name,
			Qtype:  r.RrType,
			Qclass: r.Class,
		}
		record, _ := i.redisRepo.HGetAll(ctx, q.String())
		if len(record) != 0 {
			continue
		}

		field := fmt.Sprintf("%s-0", domain.Answer)
		err = i.redisRepo.HSet(ctx, q.String(), field, r.Record, 0)
		if err != nil {
			return err
		}

		err = i.createFakeAAAA(ctx, q)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *initUseCase) createFakeAAAA(ctx context.Context, q dns.Question) error {
	soa, err := i.getFakeSOA(ctx, q)
	if err != nil {
		return err
	}

	field := fmt.Sprintf("%s-%d", domain.Ns, 0)
	q.Qtype = dns.TypeAAAA
	err = i.redisRepo.HSet(ctx, q.String(), field, soa.String(), 0)
	if err != nil {
		return err
	}

	return nil
}

func (i *initUseCase) getFakeSOA(ctx context.Context, q dns.Question) (*dns.SOA, error) {
	rrMap, err := i.redisRepo.HGetAll(ctx, q.String())
	if err != nil {
		return nil, err
	}
	for _, v := range rrMap {
		rr, _ := dns.NewRR(v)
		rr.Header().Rrtype = dns.TypeSOA
		return &dns.SOA{
			Hdr:     *rr.Header(),
			Ns:      rr.Header().Name,
			Mbox:    rr.Header().Name,
			Serial:  0,
			Refresh: rr.Header().Ttl,
			Retry:   300,
			Expire:  rr.Header().Ttl,
			Minttl:  rr.Header().Ttl,
		}, nil
	}

	return nil, fmt.Errorf("the A record isn't existed")
}

func NewInitUseCase(injector *do.Injector) (domain.InitUseCase, error) {
	return &initUseCase{
		do.MustInvoke[domain.RedisRepo](injector),
		do.MustInvoke[domain.RecordRepo](injector),
	}, nil
}
