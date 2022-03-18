package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	hctx "git.zhwenxue.com/zhgo/gocontrib/context"
	"git.zhwenxue.com/zhgo/gocontrib/log"
)

func TestCache(t *testing.T) {
	logger := log.New(os.Stdout, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	ctx := hctx.GetContext(context.Background(), "")

	dir, _ := os.Getwd()
	configPath := dir + "/cache_config_sample.yml"
	config, err := ConfigWithPath(configPath)
	assert.Nil(t, err)
	fmt.Printf("config %+v\n", config)

	cache := NewCache(config, logger)

	key := "test_key"
	value := config
	expireSeconds := 30 * 60

	// Set
	err = cache.Set(ctx, key, config, expireSeconds)
	assert.Nil(t, err)
	fmt.Printf("set %+v %+v %+v\n", key, value, err)

	// Get
	val, err := cache.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, value, val)
	fmt.Printf("get %+v %+v %+v\n", key, val, err)
	v, _ := val.(FreeCacheConfig)
	fmt.Printf("v : %+v\n", v)

	// Del
	err = cache.Del(ctx, key)
	assert.Nil(t, err)

	// Stats
	stats := cache.Stats()
	fmt.Printf("stats : %+v\n", stats)
}
