package do

import (
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/repository/db"
	"github.com/cewuandy/go-restful-dns/internal/repository/redis"
)

func ProvideRepository(injector *do.Injector) {
	do.Provide(injector, redis.NewRedisRepo)

	do.Provide(injector, db.NewRecordsRepo)
}
