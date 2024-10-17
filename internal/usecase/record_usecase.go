package usecase

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/utils"
)

type recordUseCase struct {
	redisRepo domain.RedisRepo

	recordRepo domain.RecordRepo
}

func (r *recordUseCase) CreateRecord(ctx context.Context, rr dns.RR) error {
	header := rr.Header()
	q := dns.Question{
		Name:   header.Name,
		Qtype:  header.Rrtype,
		Qclass: header.Class,
	}
	field := fmt.Sprintf("%s-%d", domain.Answer, 0)

	if header.Rrtype == dns.TypeAAAA {
		if r.isNsExisted(ctx, q.String()) {
			err := r.redisRepo.HDel(ctx, q.String())
			if err != nil {
				return err
			}
		}
	}

	record, err := r.recordRepo.Get(ctx, q.Name, q.Qtype, q.Qclass)
	if err != nil && !strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
		return err
	}
	if record != nil {
		return &domain.Error{
			Message:    "the record is already existed.",
			StatusCode: http.StatusBadRequest,
		}
	}

	record = &domain.Record{
		Name:   header.Name,
		RrType: header.Rrtype,
		Class:  header.Class,
		Record: rr.String(),
	}

	err = r.recordRepo.Create(ctx, *record)
	if err != nil {
		return err
	}

	err = r.redisRepo.HSet(ctx, q.String(), field, rr.String(), 0)
	if err != nil {
		return err
	}

	if header.Rrtype != dns.TypeA {
		return nil
	}

	return r.createFakeAAAA(ctx, header)
}

func (r *recordUseCase) GetRecord(ctx context.Context, question domain.Question) (dns.RR, error) {
	var (
		record *domain.Record
		rr     dns.RR
		err    error
	)
	question.Name = utils.GetFQDNFromDomainName(question.Name)
	t := domain.RRTypeMap[question.Qtype]
	c := domain.ClassMap[question.Qclass]
	record, err = r.recordRepo.Get(ctx, question.Name, t, c)
	if err != nil {
		return nil, err
	}

	rr, err = dns.NewRR(record.Record)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (r *recordUseCase) ListRecords(ctx context.Context) ([]dns.RR, error) {
	var (
		records []*domain.Record
		rrs     []dns.RR
		err     error
	)

	records, err = r.recordRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		var rr dns.RR
		rr, err = dns.NewRR(record.Record)
		if err != nil {
			return nil, err
		}
		rrs = append(rrs, rr)
	}

	return rrs, nil
}

func (r *recordUseCase) UpdateRecord(ctx context.Context, rr dns.RR) error {
	record := &domain.Record{
		Name:   utils.GetFQDNFromDomainName(rr.Header().Name),
		RrType: rr.Header().Rrtype,
		Class:  rr.Header().Class,
		Record: rr.String(),
	}

	err := r.recordRepo.Update(ctx, record)
	if err != nil {
		return err
	}

	err = r.redisRepo.HDel(ctx, rr.String())
	if err != nil {
		return err
	}

	q := dns.Question{
		Name:   rr.Header().Name,
		Qtype:  rr.Header().Rrtype,
		Qclass: rr.Header().Class,
	}
	field := fmt.Sprintf("%s-%d", domain.Answer, 0)
	return r.redisRepo.HSet(ctx, q.String(), field, rr.String(), 0)
}

func (r *recordUseCase) DeleteRecord(ctx context.Context, question domain.Question) error {
	question.Name = utils.GetFQDNFromDomainName(question.Name)
	t := domain.RRTypeMap[question.Qtype]
	c := domain.ClassMap[question.Qclass]
	err := r.recordRepo.Delete(ctx, question.Name, t, c)
	if err != nil {
		return err
	}

	err = r.deleteFakeAAAA(ctx, question)
	if err != nil {
		return err
	}

	return r.redisRepo.HDel(ctx, question.String())
}

func (r *recordUseCase) createFakeAAAA(ctx context.Context, header *dns.RR_Header) error {
	q := dns.Question{
		Name:   header.Name,
		Qtype:  header.Rrtype,
		Qclass: header.Class,
	}
	soa, err := r.getFakeSOA(ctx, q)
	if err != nil {
		return err
	}
	field := fmt.Sprintf("%s-%d", domain.Ns, 0)
	q.Qtype = dns.TypeAAAA
	err = r.redisRepo.HSet(ctx, q.String(), field, soa.String(), 0)
	if err != nil {
		return err
	}

	return nil
}

func (r *recordUseCase) deleteFakeAAAA(ctx context.Context, question domain.Question) error {
	question.Qtype = domain.TypeAAAA
	if !r.isNsExisted(ctx, question.String()) {
		return nil
	}

	return r.redisRepo.HDel(ctx, question.String())
}

func (r *recordUseCase) getFakeSOA(ctx context.Context, q dns.Question) (*dns.SOA, error) {
	rrMap, err := r.redisRepo.HGetAll(ctx, q.String())
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

func (r *recordUseCase) isNsExisted(ctx context.Context, key string) bool {
	rrMap, _ := r.redisRepo.HGetAll(ctx, key)
	if len(rrMap) == 0 {
		return false
	}
	for k := range rrMap {
		if strings.Contains(k, string(domain.Answer)) {
			return false
		}
	}

	return true
}

func NewRecordUseCase(injector *do.Injector) (domain.RecordUseCase, error) {
	return &recordUseCase{
		do.MustInvoke[domain.RedisRepo](injector),
		do.MustInvoke[domain.RecordRepo](injector),
	}, nil
}
