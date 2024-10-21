package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redismock/v8"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"
)

type redisRepoTestSuite struct {
	suite.Suite

	injector *do.Injector

	key   string
	field string
	value string
}

func TestRedisRepo(t *testing.T) {
	suite.Run(t, &redisRepoTestSuite{})
}

func (t *redisRepoTestSuite) SetupSuite() {
	t.injector = do.New()
	t.key = "key"
	t.field = "field"
	t.value = "value"
}

func (t *redisRepoTestSuite) SetupTest() {

}

func (t *redisRepoTestSuite) TestHSet() {
	t.Run(
		"success", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHSet(t.key, t.field, t.value).
				SetVal(0)
			client.
				ExpectExpire(t.key, time.Minute).
				SetVal(true)

			err := repo.HSet(context.Background(), t.key, t.field, t.value, time.Minute)
			t.Nil(err)
		},
	)

	t.Run(
		"success_expire_0", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHSet(t.key, t.field, t.value).
				SetVal(0)
			client.
				ExpectExpire(t.key, 0).
				SetVal(true)

			err := repo.HSet(context.Background(), t.key, t.field, t.value, 0)
			t.Nil(err)
		},
	)

	t.Run(
		"HSet_error", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHSet(t.key, t.field, t.value).
				SetErr(fmt.Errorf("test-error"))

			err := repo.HSet(context.Background(), t.key, t.field, t.value, time.Minute)
			t.NotNil(err)
			t.Contains(err.Error(), "Redis error:")
		},
	)

	t.Run(
		"Expire_error", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHSet(t.key, t.field, t.value).
				SetVal(0)
			client.
				ExpectExpire(t.key, time.Minute).
				SetErr(fmt.Errorf("test-error"))

			err := repo.HSet(context.Background(), t.key, t.field, t.value, time.Minute)
			t.NotNil(err)
			t.Contains(err.Error(), "Redis error:")
		},
	)
}

func (t *redisRepoTestSuite) TestHGetAll() {
	t.Run(
		"success", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHGetAll(t.key).
				SetVal(map[string]string{"key": "value"})

			val, err := repo.HGetAll(context.Background(), t.key)
			t.Equal("value", val["key"])
			t.Nil(err)
		},
	)

	t.Run(
		"HGetAll_error", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHGetAll(t.key).
				SetErr(fmt.Errorf("test-error"))

			val, err := repo.HGetAll(context.Background(), t.key)
			t.Nil(val)
			t.NotNil(err)
			t.Contains(err.Error(), "Redis error:")
		},
	)
}

func (t *redisRepoTestSuite) TestHDel() {
	t.Run(
		"success", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHGetAll(t.key).
				SetVal(map[string]string{t.field: t.value})
			client.
				ExpectHDel(t.key, t.field).
				SetVal(0)

			err := repo.HDel(context.Background(), t.key)
			t.Nil(err)
		},
	)

	t.Run(
		"HGetAll_error", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHGetAll(t.key).
				SetErr(fmt.Errorf("HGetAll_error"))

			err := repo.HDel(context.Background(), t.key)
			t.NotNil(err)
			t.Contains(err.Error(), "Redis error:")
			t.Contains(err.Error(), "HGetAll_error")
		},
	)

	t.Run(
		"HDel_error", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectHGetAll(t.key).
				SetVal(map[string]string{t.field: t.value})
			client.
				ExpectHDel(t.key, t.field).
				SetErr(fmt.Errorf("HDel_error"))

			err := repo.HDel(context.Background(), t.key)
			t.NotNil(err)
			t.Contains(err.Error(), "Redis error:")
			t.Contains(err.Error(), "HDel_error")
		},
	)
}

func (t *redisRepoTestSuite) TestFlushAll() {
	t.Run(
		"success", func() {
			redisClient, client := redismock.NewClientMock()
			do.OverrideValue[*redis.Client](t.injector, redisClient)
			repo, _ := NewRedisRepo(t.injector)
			client.
				ExpectFlushAll().
				SetVal("")
			err := repo.FlushAll(context.Background())
			t.Nil(err)
		},
	)
}
