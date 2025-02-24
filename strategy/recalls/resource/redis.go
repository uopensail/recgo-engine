package resource

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/datastruct"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type redisSource struct {
	connRef utils.Reference
	conn    *redis.Client
}

func getRedisConn(cfg table.RedisResourceConfig) *redis.Client {
	stat := prome.NewStat("getRedisConn")
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

func newRedisSource(cfg table.RedisResourceConfig) redisSource {
	conn := getRedisConn(cfg)
	source := redisSource{

		conn: conn,
	}
	source.connRef.CloseHandler = func() {
		if conn != nil {
			conn.Close()
		}
	}
	return source
}

func (s *redisSource) Close() {
	if s.connRef.CloseHandler != nil {
		s.connRef.LazyFree(10)
	}
}

type redisKVWithScoreSource struct {
	redisSource
}

func newRedisKVWithScoreSource(cfg table.RedisResourceConfig) *redisKVWithScoreSource {
	rs := &redisKVWithScoreSource{
		redisSource: newRedisSource(cfg),
	}

	return rs
}

func (rs *redisKVWithScoreSource) get(pl *pool.Pool, keys []string) [][]datastruct.Tuple[int, float32] {
	stat := prome.NewStat("redisKVWithScoreSource.get")
	defer stat.End()
	rs.connRef.Retain()
	defer rs.connRef.Release()

	retSlice := rs.conn.MGet(keys...)
	result := retSlice.Val()
	if retSlice == nil || len(result) == 0 {
		stat.MarkErr()
		zlog.LOG.Error("redis mget result is nil")
		return nil
	}

	values := make([][]datastruct.Tuple[int, float32], 0, len(keys))
	for i := 0; i < len(result); i++ {
		if result[i] != nil {
			ss := utils.StringSplit(result[i].(string), ",")
			tupleList := make([]datastruct.Tuple[int, float32], 0, len(ss))
			for j := 0; j < len(ss); j++ {
				tmp := utils.StringSplit(ss[j], ":")
				if len(tmp) >= 2 {
					itemID := tmp[0]
					if item := pl.GetByKey(itemID); item != nil {
						tupleList = append(tupleList, datastruct.Tuple[int, float32]{
							item.ID, utils.String2Float32(tmp[1])})
					}

				}
			}
			values = append(values, tupleList)
		}
	}
	return values
}

type RedisResource struct {
	cfg table.RecallResourceMeta
	*redisKVWithScoreSource
}

func NewRedisResource(_ config.EnvConfig, cfg table.RecallResourceMeta, pl *pool.Pool) *RedisResource {
	cfg.ParseRedisSource()
	return &RedisResource{
		cfg:                    cfg,
		redisKVWithScoreSource: newRedisKVWithScoreSource(cfg.RedisResourceConfig),
	}
}

func (res *RedisResource) Get(keys []string, pl *pool.Pool) [][]datastruct.Tuple[int, float32] {
	return res.redisKVWithScoreSource.get(pl, keys)
}

func (res *RedisResource) CheckResourceUpdate(envCfg config.EnvConfig, poolUpdate bool) bool {
	return false
}

func (res *RedisResource) Meta() *table.RecallResourceMeta {
	return &res.cfg
}

func (res *RedisResource) Close() {
	res.redisKVWithScoreSource.Close()
}
