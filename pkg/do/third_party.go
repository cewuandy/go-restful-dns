package do

import (
	"flag"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm/logger"
	"os"
	"strings"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	pkgGorm "github.com/cewuandy/go-restful-dns/pkg/gorm"
	"github.com/cewuandy/go-restful-dns/pkg/options"

	redisLib "github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ProvideThirdPartyElement(injector *do.Injector) {
	do.Provide(injector, provideEnv)
	do.Provide(injector, provideUpstreams)
	do.Provide(injector, provideRedisClient)
	do.Provide[*gorm.DB](injector, provideSqliteClient)
}

func provideEnv(*do.Injector) (*domain.Options, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.Usage = func() {}
	env := &domain.Options{}
	_ = options.LoadDefaultConfig(flagSet, env)
	_ = options.LoadCliFlagConfigs(flagSet)
	return env, nil
}

func provideUpstreams(injector *do.Injector) ([]string, error) {
	env := do.MustInvoke[*domain.Options](injector)
	return strings.Split(env.UpstreamForwarders, ","), nil
}

func provideRedisClient(injector *do.Injector) (*redis.Client, error) {
	env := do.MustInvoke[*domain.Options](injector)
	return redisLib.NewClient(
		&redisLib.Options{
			Addr: env.RedisAddr,
			DB:   0,
		},
	), nil
}

func provideSqliteClient(injector *do.Injector) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open("dns.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return nil, err
	}

	err = pkgGorm.AutoMigrate(db)
	if err != nil {
		return db, err
	}

	return db, nil
}
