package domain

import (
	"github.com/miekg/dns"
	"net"
)

type A struct {
	dns.A   `json:"-"`
	Hdr     RR_Header `json:"hdr" binding:"required"`
	Address net.IP    `json:"a" binding:"required" swaggertype:"string"` // should be like "192.168.0.1"
}

func (a *A) String() string {
	a.A.A = a.Address
	a.A.Hdr = a.Hdr.Copy()
	return a.A.String()
}
