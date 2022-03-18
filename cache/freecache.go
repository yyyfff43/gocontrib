package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"runtime/debug"
	"time"

	"git.zhwenxue.com/zhgo/gocontrib/log"
	"github.com/coocood/freecache"
)

type FreeCacheConfig struct {
	// 缓存大小，单位：MB，预分配内存
	CacheSizeMB int `yaml:"cacheSizeMB"`
	// 当设置了较大的缓存大小时，应调用 debug.SetGCPercent()，设置一个较小值，限制内存消耗和 GC 暂停时间
	GCPercent int `yaml:"gcPercent"`
}

type freeCache struct {
	cache  *freecache.Cache
	logger *log.Logger
}

type FreeCacheStats struct {
	// 驱逐发生的次数
	EvacuateCount int64 `json:"evacuate_count"`
	// 过期发生的次数
	ExpiredCount int64 `json:"expired_count"`
	// 当前缓存中的项目数
	EntryCount int64 `json:"entry_count"`
	// 访问条目时的平均 unix 时间戳
	AverageAccessTime int64 `json:"average_access_time"`
	// 在缓存中找到键的次数
	HitCount int64 `json:"hit_count"`
	// 缓存中发生未命中的次数
	MissCount int64 `json:"miss_count"`
	// 对给定键的查找发生的次数
	LookupCount int64 `json:"lookup_count"`
	// 命中与查找的比率
	HitRate float64 `json:"hit_rate"`
	// 条目被覆盖的次数
	OverwriteCount int64 `json:"overwrite_count"`
	// 条目的过期时间延长的次数
	TouchedCount int64 `json:"touched_count"`
}

func newFreeCache(cnf FreeCacheConfig, log *log.Logger) *freeCache {
	c := freecache.NewCache(cnf.CacheSizeMB * 1024 * 1024) // MB => Byte
	debug.SetGCPercent(cnf.GCPercent)                      // !!!

	cache := &freeCache{
		cache:  c,
		logger: log,
	}

	return cache
}

func (c *freeCache) Get(ctx context.Context, key string) (interface{}, error) {
	defer c.trace(ctx, "get", key, time.Now())

	valueBytes, err := c.cache.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	return c.deserialize(valueBytes)
}

func (c *freeCache) Set(ctx context.Context, key string, value interface{}, expireSeconds int) error {
	defer c.trace(ctx, "set", key, time.Now())

	valueBytes, err := c.serialize(value)
	if err != nil {
		return err
	}

	return c.cache.Set([]byte(key), valueBytes, expireSeconds)
}

func (c *freeCache) Del(ctx context.Context, key string) error {
	defer c.trace(ctx, "del", key, time.Now())

	ok := c.cache.Del([]byte(key))
	if !ok {
		return errors.New("freecache delete failed")
	}

	return nil
}

func (c *freeCache) Stats() interface{} {
	return FreeCacheStats{
		EvacuateCount:     c.cache.EvacuateCount(),
		ExpiredCount:      c.cache.ExpiredCount(),
		EntryCount:        c.cache.EntryCount(),
		AverageAccessTime: c.cache.AverageAccessTime(),
		HitCount:          c.cache.HitCount(),
		MissCount:         c.cache.MissCount(),
		LookupCount:       c.cache.LookupCount(),
		HitRate:           c.cache.HitRate(),
		OverwriteCount:    c.cache.OverwriteCount(),
		TouchedCount:      c.cache.TouchedCount(),
	}
}

func (c *freeCache) trace(ctx context.Context, cmd, key string, start time.Time) {
	c.logger.Info(ctx, "Cache",
		log.String("cmd", cmd),
		log.String("key", key),
		log.Duration("time", time.Since(start)),
	)
}

func (c *freeCache) serialize(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *freeCache) deserialize(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
