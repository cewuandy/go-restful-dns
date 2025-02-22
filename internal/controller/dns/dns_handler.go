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

	resp = new(dns.Msg)

	defer func() {
		if respWriter == nil {
			fmt.Printf("Error: respWriter is nil\n")
			return
		}
		if resp == nil {
			fmt.Printf("Error: resp is nil\n")
			return
		}
		if len(resp.Answer) > 0 || len(resp.Ns) > 0 || len(resp.Extra) > 0 {
			err = respWriter.WriteMsg(resp)
			if err != nil {
				fmt.Printf("Error writing response: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Warning: No data to write in response\n")
		}
	}()

	resp, err = d.dnsUseCase.QueryRedisCache(context.Background(), req)
	if err != nil {
		fmt.Printf("Error querying Redis cache: %s\n", err.Error())
		return
	}
	if resp != nil && (len(resp.Answer) > 0 || len(resp.Ns) > 0 || len(resp.Extra) > 0) {
		return
	}

	resp, err = d.dnsUseCase.QueryUpstream(context.Background(), req)
	if err != nil {
		fmt.Printf("Error querying upstream: %s\n", err.Error())
		return
	}
}

func NewDNSHandler(injector *do.Injector) (dns.Handler, error) {
	return &dnsHandler{do.MustInvoke[domain.DNSUseCase](injector)}, nil
}
