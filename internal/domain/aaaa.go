package domain

import (
	"github.com/miekg/dns"
	"net"
)

type AAAA struct {
	dns.AAAA
	Hdr     RR_Header `json:"hdr" binding:"required"`
	Address net.IP    `json:"aaaa" binding:"required"`
}

func (a *AAAA) String() string {
	a.AAAA.AAAA = a.Address
	a.AAAA.Hdr = a.Hdr.Copy()
	return a.AAAA.String()
}
