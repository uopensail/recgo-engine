package resource

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type RedisResource struct {
	connRef utils.Reference
	redis   *redis.Client
	table.FilterResourceMeta
}

// 获得redis的连接
func getRedisConn(cfg table.RedisConfigure) *redis.Client {
	stat := prome.NewStat("getRedisConnection")
	defer stat.End()
	minIdleConns := cfg.MinIdleConns
	if minIdleConns <= 0 {
		minIdleConns = 1
	}
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		zlog.LOG.Error("redis parse url error", zap.Error(err))
		return nil
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30
	}
	opt.ReadTimeout = time.Millisecond * time.Duration(timeout)
	opt.DialTimeout = time.Millisecond * time.Duration(timeout)

	return redis.NewClient(opt)
}

func NewRedisResource(cfg table.FilterResourceMeta, _ config.EnvConfig) *RedisResource {
	conn := getRedisConn(cfg.Redis)
	if conn == nil {
		return nil
	}
	res := &RedisResource{
		redis:              conn,
		FilterResourceMeta: cfg,
	}
	res.connRef.CloseHandler = func() {
		if conn != nil {
			conn.Close()
		}
	}
	return res
}

func (s *RedisResource) Do(ctx context.Context, key string, param map[string]string) ([]string, error) {
	stat := prome.NewStat("RedisSource.LRange")
	defer stat.End()
	s.connRef.Retain()
	defer s.connRef.Release()
	length := utils.GetInt64Param(param, "max_length", 100)
	list, err := s.redis.LRange(key, 0, length).Result()
	if err != nil {
		stat.MarkErr()
		zlog.LOG.Error("redis lrange error", zap.Error(err))
		return nil, err
	}
	return list, nil
}
func (s *RedisResource) Meta() *table.FilterResourceMeta {
	return &s.FilterResourceMeta
}
func (s *RedisResource) Close() {
	if s.connRef.CloseHandler != nil {
		s.connRef.LazyFree(10)
	}
}
