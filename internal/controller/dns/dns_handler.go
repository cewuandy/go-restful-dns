package dns

import (
	"context"
	"fmt"

	"github.com/miekg/dns"
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

type dnsHandler struct {
	dnsUseCase domain.DNSUseCase
}

func (d *dnsHandler) ServeDNS(respWriter dns.ResponseWriter, req *dns.Msg) {
	var (
		resp *dns.Msg
		err  error
	)

	defer func() {
		if resp == nil {
			return
		}
		err = respWriter.WriteMsg(resp)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}()

	resp, err = d.dnsUseCase.QueryRedisCache(context.Background(), req)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	if resp.Answer != nil || resp.Ns != nil || resp.Extra != nil {
		return
	}

	resp, err = d.dnsUseCase.QueryUpstream(context.Background(), req)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	if resp.Answer != nil || resp.Ns != nil || resp.Extra != nil {
		return
	}
}

func NewDNSHandler(injector *do.Injector) (dns.Handler, error) {
	return &dnsHandler{do.MustInvoke[domain.DNSUseCase](injector)}, nil
}
