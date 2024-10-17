// @title go-restful-dns API
// @version 1.0
// @description Documentation of go-restful-dns API
// @host localhost:8081
// @BasePath /api/v1

package main

import (
	"context"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	pkgDo "github.com/cewuandy/go-restful-dns/pkg/do"
)

func main() {
	injector := pkgDo.Injector
	pkgDo.ProvideThirdPartyElement(injector)
	pkgDo.ProvideServer(injector)
	pkgDo.ProvideRepository(injector)
	pkgDo.ProvideUseCase(injector)
	pkgDo.ProvideController(injector)

	err := do.MustInvoke[domain.InitHandler](injector).Initialize(context.Background())
	if err != nil {
		panic(err)
	}

	dnsServer := do.MustInvoke[*dns.Server](injector)
	httpServer := do.MustInvoke[*http.Server](injector)

	startServer(dnsServer.ListenAndServe, httpServer.ListenAndServe)
	startWaitForShutdown(
		func() error {
			return httpServer.Shutdown(context.Background())
		},
		dnsServer.Shutdown,
	)
}

func startServer(start ...func() error) {
	for _, f := range start {
		go func() {
			err := f()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func startWaitForShutdown(shutdown ...func() error) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-sigChan:
			for _, f := range shutdown {
				err := f()
				if err != nil {
					log.Fatal(err)
				}
			}
			return
		}
	}
}
