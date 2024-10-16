package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"time"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

type redisRepository struct {
	client *redis.Client
}

func (r *redisRepository) HSet(ctx context.Context, key string, field string, value string,
	expiration time.Duration) error {
	var err error

	err = r.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		err = &domain.Error{
			Message: fmt.Sprintf("Redis error: %s", err.Error()),
		}
		return err
	}

	if expiration == 0 {
		return nil
	}

	err = r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		err = &domain.Error{
			Message: fmt.Sprintf("Redis error: %s", err.Error()),
		}
		return err
	}

	return nil
}

func (r *redisRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	var (
		val map[string]string
		err error
	)

	val, err = r.client.HGetAll(ctx, key).Result()
	if err != nil {
		err = &domain.Error{
			Message: fmt.Sprintf("Redis error: %s", err.Error()),
		}
		return nil, err
	}

	return val, nil
}

func (r *redisRepository) HDel(ctx context.Context, key string) error {
	var (
		val map[string]string
		err error
	)

	val, err = r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return &domain.Error{
			Message: fmt.Sprintf("Redis error: %s", err.Error()),
		}
	}
	for field := range val {
		_, err = r.client.HDel(ctx, key, field).Result()
		if err != nil {
			err = &domain.Error{
				Message: fmt.Sprintf("Redis error: %s", err.Error()),
			}
			return err
		}
	}

	return nil
}

// NewRedisRepo init redisRepository
func NewRedisRepo(injector *do.Injector) (domain.RedisRepo, error) {
	return &redisRepository{do.MustInvoke[*redis.Client](injector)}, nil
}
