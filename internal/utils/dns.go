package utils

import (
	"fmt"
	"strings"
)

// GetFQDNFromDomainName get fqdn from domain name
func GetFQDNFromDomainName(domainName string) (fqdn string) {
	domainName = strings.TrimSuffix(domainName, ".")
	fqdn = fmt.Sprintf("%s.", domainName)
	return
}
