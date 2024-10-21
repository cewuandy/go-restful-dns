package domain

import (
	"context"
	"time"
)

type RedisRepo interface {
	HSet(ctx context.Context, key string, field string, value string,
		expiration time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
}
