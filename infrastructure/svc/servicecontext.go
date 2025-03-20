package svc

import (
	"context"
	"oauth2/infrastructure/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
	Cache  cache.ClusterConf
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := logx.WithContext(context.Background())
	redisConf := redis.RedisConf{
		Host: c.Redis.Host,
		Type: c.Redis.Type,
		Pass: c.Redis.Pass,
		Tls:  c.Redis.Tls,
	}
	cacheConf := cache.CacheConf{
		{
			RedisConf: redisConf,
			Weight:    100,
		},
	}
	conn := sqlx.NewMysql(c.DB.DataSource)

	redisClient, err := redis.NewRedis(redisConf)
	if err != nil {
		logger.Errorf("Failed to create Redis client: %v", err)
	}
	return &ServiceContext{
		Config: c,
		DB:     conn,
		Cache:  cacheConf,
		Redis:  redisClient,
	}
}
