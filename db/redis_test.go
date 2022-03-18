package db

import (
	"context"
	"fmt"
	hctx "git.zhwenxue.com/zhgo/gocontrib/context"
	"git.zhwenxue.com/zhgo/gocontrib/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	logger := log.New(os.Stdout, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	ctx := hctx.GetContext(context.Background(), "")
	// RedisTestOption 测试
	//RedisTestOption := &RedisServerOption{
	//	Host:               "redis",
	//	Port:               6379,
	//	Pwd:                "",
	//	DB:                 0,
	//	PoolSize:           3,
	//	IdlePoolSize:       1,
	//	PoolTimeout:        2 * time.Second,
	//	IdleTimeout:        100 * time.Second,
	//	IdleCheckFrequency: 40 * time.Second,
	//	ExecSlowTime:       80 * time.Millisecond,
	//	Log:                *logger,
	//}
	dir, _ := os.Getwd()
	configPath := dir + "/redis_config.yml"
	rConfig, err := RedisConfigWithPath(configPath)
	assert.Nil(t, err)
	//redis
	redisDB, err := NewRedis(rConfig, *logger)
	assert.Nil(t, err)

	rStatus := redisDB.Set(ctx, "hachi", "hachi123", time.Second*20)
	assert.Nil(t, rStatus.Err())
	info := redisDB.Get(ctx, "hachi")
	assert.Nil(t, info.Err())
	fmt.Println(info.Val())
}
