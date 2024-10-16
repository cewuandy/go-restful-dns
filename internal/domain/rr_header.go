package domain

import (
	"github.com/miekg/dns"
)

type RRType string

const (
	TypeNone       RRType = "None"
	TypeA          RRType = "A"
	TypeNS         RRType = "NS"
	TypeMD         RRType = "MD"
	TypeMF         RRType = "MF"
	TypeCNAME      RRType = "CNAME"
	TypeSOA        RRType = "SOA"
	TypeMB         RRType = "MB"
	TypeMG         RRType = "MG"
	TypeMR         RRType = "MR"
	TypeNULL       RRType = "NULL"
	TypePTR        RRType = "PTR"
	TypeHINFO      RRType = "HINFO"
	TypeMINFO      RRType = "MINFO"
	TypeMX         RRType = "MX"
	TypeTXT        RRType = "TXT"
	TypeRP         RRType = "RP"
	TypeAFSDB      RRType = "AFSDB"
	TypeX25        RRType = "X25"
	TypeISDN       RRType = "ISDN"
	TypeRT         RRType = "RT"
	TypeNSAPPTR    RRType = "NSAPPTR"
	TypeSIG        RRType = "SIG"
	TypeKEY        RRType = "KEY"
	TypePX         RRType = "PX"
	TypeGPOS       RRType = "GPOS"
	TypeAAAA       RRType = "AAAA"
	TypeLOC        RRType = "LOC"
	TypeNXT        RRType = "NXT"
	TypeEID        RRType = "EID"
	TypeNIMLOC     RRType = "NIMLOC"
	TypeSRV        RRType = "SRV"
	TypeATMA       RRType = "ATMA"
	TypeNAPTR      RRType = "NAPTR"
	TypeKX         RRType = "KX"
	TypeCERT       RRType = "CERT"
	TypeDNAME      RRType = "DNAME"
	TypeOPT        RRType = "OPT"
	TypeAPL        RRType = "APL"
	TypeDS         RRType = "DS"
	TypeSSHFP      RRType = "SSHFP"
	TypeIPSECKEY   RRType = "IPSECKEY"
	TypeRRSIG      RRType = "RRSIG"
	TypeNSEC       RRType = "NSEC"
	TypeDNSKEY     RRType = "DNSKEY"
	TypeDHCID      RRType = "DHCID"
	TypeNSEC3      RRType = "NSEC3"
	TypeNSEC3PARAM RRType = "NSEC3PARAM"
	TypeTLSA       RRType = "TLSA"
	TypeSMIMEA     RRType = "SMIMEA"
	TypeHIP        RRType = "HIP"
	TypeNINFO      RRType = "NINFO"
	TypeRKEY       RRType = "RKEY"
	TypeTALINK     RRType = "TALINK"
	TypeCDS        RRType = "CDS"
	TypeCDNSKEY    RRType = "CDNSKEY"
	TypeOPENPGPKEY RRType = "OPENPGPKEY"
	TypeCSYNC      RRType = "CSYNC"
	TypeZONEMD     RRType = "ZONEMD"
	TypeSVCB       RRType = "SVCB"
	TypeHTTPS      RRType = "HTTPS"
	TypeSPF        RRType = "SPF"
	TypeUINFO      RRType = "UINFO"
	TypeUID        RRType = "UID"
	TypeGID        RRType = "GID"
	TypeUNSPEC     RRType = "UNSPEC"
	TypeNID        RRType = "NID"
	TypeL32        RRType = "L32"
	TypeL64        RRType = "L64"
	TypeLP         RRType = "LP"
	TypeEUI48      RRType = "EUI48"
	TypeEUI64      RRType = "EUI64"
	TypeNXNAME     RRType = "NXNAME"
	TypeURI        RRType = "URI"
	TypeCAA        RRType = "CAA"
	TypeAVC        RRType = "AVC"
	TypeAMTRELAY   RRType = "AMTRELAY"
	TypeTKEY       RRType = "TKEY"
	TypeTSIG       RRType = "TSIG"
	TypeIXFR       RRType = "IXFR"
	TypeAXFR       RRType = "AXFR"
	TypeMAILB      RRType = "MAILB"
	TypeMAILA      RRType = "MAILA"
	TypeANY        RRType = "ANY"
	TypeTA         RRType = "TA"
	TypeDLV        RRType = "DLV"
	TypeReserved   RRType = "Reserved"
)

var RRTypeMap = map[RRType]uint16{
	TypeNone:       dns.TypeNone,
	TypeA:          dns.TypeA,
	TypeNS:         dns.TypeNS,
	TypeMD:         dns.TypeMD,
	TypeMF:         dns.TypeMF,
	TypeCNAME:      dns.TypeCNAME,
	TypeSOA:        dns.TypeSOA,
	TypeMB:         dns.TypeMB,
	TypeMG:         dns.TypeMG,
	TypeMR:         dns.TypeMR,
	TypeNULL:       dns.TypeNULL,
	TypePTR:        dns.TypePTR,
	TypeHINFO:      dns.TypeHINFO,
	TypeMINFO:      dns.TypeMINFO,
	TypeMX:         dns.TypeMX,
	TypeTXT:        dns.TypeTXT,
	TypeRP:         dns.TypeRP,
	TypeAFSDB:      dns.TypeAFSDB,
	TypeX25:        dns.TypeX25,
	TypeISDN:       dns.TypeISDN,
	TypeRT:         dns.TypeRT,
	TypeNSAPPTR:    dns.TypeNSAPPTR,
	TypeSIG:        dns.TypeSIG,
	TypeKEY:        dns.TypeKEY,
	TypePX:         dns.TypePX,
	TypeGPOS:       dns.TypeGPOS,
	TypeAAAA:       dns.TypeAAAA,
	TypeLOC:        dns.TypeLOC,
	TypeNXT:        dns.TypeNXT,
	TypeEID:        dns.TypeEID,
	TypeNIMLOC:     dns.TypeNIMLOC,
	TypeSRV:        dns.TypeSRV,
	TypeATMA:       dns.TypeATMA,
	TypeNAPTR:      dns.TypeNAPTR,
	TypeKX:         dns.TypeKX,
	TypeCERT:       dns.TypeCERT,
	TypeDNAME:      dns.TypeDNAME,
	TypeOPT:        dns.TypeOPT,
	TypeAPL:        dns.TypeAPL,
	TypeDS:         dns.TypeDS,
	TypeSSHFP:      dns.TypeSSHFP,
	TypeIPSECKEY:   dns.TypeIPSECKEY,
	TypeRRSIG:      dns.TypeRRSIG,
	TypeNSEC:       dns.TypeNSEC,
	TypeDNSKEY:     dns.TypeDNSKEY,
	TypeDHCID:      dns.TypeDHCID,
	TypeNSEC3:      dns.TypeNSEC3,
	TypeNSEC3PARAM: dns.TypeNSEC3PARAM,
	TypeTLSA:       dns.TypeTLSA,
	TypeSMIMEA:     dns.TypeSMIMEA,
	TypeHIP:        dns.TypeHIP,
	TypeNINFO:      dns.TypeNINFO,
	TypeRKEY:       dns.TypeRKEY,
	TypeTALINK:     dns.TypeTALINK,
	TypeCDS:        dns.TypeCDS,
	TypeCDNSKEY:    dns.TypeCDNSKEY,
	TypeOPENPGPKEY: dns.TypeOPENPGPKEY,
	TypeCSYNC:      dns.TypeCSYNC,
	TypeZONEMD:     dns.TypeZONEMD,
	TypeSVCB:       dns.TypeSVCB,
	TypeHTTPS:      dns.TypeHTTPS,
	TypeSPF:        dns.TypeSPF,
	TypeUINFO:      dns.TypeUINFO,
	TypeUID:        dns.TypeUID,
	TypeGID:        dns.TypeGID,
	TypeUNSPEC:     dns.TypeUNSPEC,
	TypeNID:        dns.TypeNID,
	TypeL32:        dns.TypeL32,
	TypeL64:        dns.TypeL64,
	TypeLP:         dns.TypeLP,
	TypeEUI48:      dns.TypeEUI48,
	TypeEUI64:      dns.TypeEUI64,
	TypeNXNAME:     dns.TypeNXNAME,
	TypeURI:        dns.TypeURI,
	TypeCAA:        dns.TypeCAA,
	TypeAVC:        dns.TypeAVC,
	TypeAMTRELAY:   dns.TypeAMTRELAY,
	TypeTKEY:       dns.TypeTKEY,
	TypeTSIG:       dns.TypeTSIG,
	TypeIXFR:       dns.TypeIXFR,
	TypeAXFR:       dns.TypeAXFR,
	TypeMAILB:      dns.TypeMAILB,
	TypeMAILA:      dns.TypeMAILA,
	TypeANY:        dns.TypeANY,
	TypeTA:         dns.TypeTA,
	TypeDLV:        dns.TypeDLV,
	TypeReserved:   dns.TypeReserved,
}

type Class string

const (
	ClassINET   Class = "INET"
	ClassCSNET  Class = "CSNET"
	ClassCHAOS  Class = "CHAOS"
	ClassHESIOD Class = "HESIOD"
	ClassNONE   Class = "NONE"
	ClassANY    Class = "ANY"
)

var ClassMap = map[Class]uint16{
	ClassINET:   dns.ClassINET,
	ClassCSNET:  dns.ClassCSNET,
	ClassCHAOS:  dns.ClassCHAOS,
	ClassHESIOD: dns.ClassHESIOD,
	ClassNONE:   dns.ClassNONE,
	ClassANY:    dns.ClassANY,
}

type RR_Header struct {
	dns.RR_Header
	Name     string `json:"name" binding:"required"`
	Rrtype   RRType `json:"rrtype" binding:"required"`
	Class    Class  `json:"class" binding:"required"`
	Ttl      uint32 `json:"ttl" binding:"required"`
	Rdlength uint16
}

func (h *RR_Header) String() string {
	h.RR_Header.Name = h.Name
	h.RR_Header.Rrtype = RRTypeMap[h.Rrtype]
	h.RR_Header.Class = ClassMap[h.Class]
	h.RR_Header.Ttl = h.Ttl
	h.RR_Header.Rdlength = h.Rdlength
	return h.RR_Header.String()
}

func (h *RR_Header) Copy() dns.RR_Header {
	h.RR_Header.Name = h.Name
	h.RR_Header.Rrtype = RRTypeMap[h.Rrtype]
	h.RR_Header.Class = ClassMap[h.Class]
	h.RR_Header.Ttl = h.Ttl
	h.RR_Header.Rdlength = h.Rdlength
	return h.RR_Header
}
