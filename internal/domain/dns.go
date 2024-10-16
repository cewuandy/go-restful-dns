package domain

import (
	"context"
	"github.com/miekg/dns"
)

type DNSUseCase interface {
	QueryRedisCache(ctx context.Context, req *dns.Msg) (resp *dns.Msg, err error)

	QueryUpstream(ctx context.Context, req *dns.Msg) (resp *dns.Msg, err error)
}
