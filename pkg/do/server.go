package do

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"net/http"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/pkg/gin/routes"
)

func ProvideServer(injector *do.Injector) {
	do.Provide(injector, provideDNSServer)
	do.Provide(injector, provideGinServer)
	do.Provide(injector, provideHTTPServer)
}

func provideDNSServer(injector *do.Injector) (*dns.Server, error) {
	env := do.MustInvoke[*domain.Options](injector)

	dnsServer := &dns.Server{
		Addr:    fmt.Sprintf("%s:%d", env.DnsAddr, env.DnsPort),
		Net:     "udp",
		Handler: do.MustInvoke[dns.Handler](injector),
	}
	dnsServer.NotifyStartedFunc = func() {
		fmt.Printf("Listen and serve at %s for DNS query\n", dnsServer.Addr)
	}
	return dnsServer, nil
}

func provideGinServer(injector *do.Injector) (*gin.Engine, error) {
	r := gin.New()
	// TODO: should assign a real ip:port, that is workaround now
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))
	r.Use(do.MustInvoke[domain.ErrorHandler](injector).HandleError)

	routes.RegisterRecordRoutes(r, do.MustInvoke[domain.RecordHandler](injector))

	return r, nil
}

func provideHTTPServer(injector *do.Injector) (*http.Server, error) {
	env := do.MustInvoke[*domain.Options](injector)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", env.HttpAddr, env.HttpPort),
		Handler: do.MustInvoke[*gin.Engine](injector),
	}
	return httpServer, nil
}
