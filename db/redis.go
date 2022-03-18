package db

import (
	"context"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	"git.zhwenxue.com/zhgo/gocontrib/log"
	goRedis "github.com/go-redis/redis/v8"
	guuid "github.com/google/uuid"
	"gopkg.in/yaml.v2"
)


// Redis redis对象
type Redis struct {
	*goRedis.Client
}

func RedisConfigWithPath(path string) (*RedisConfig, error) {
	var cfg *Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Redis, nil
}

func NewRedis(opt *RedisConfig, log log.Logger) (*Redis, error) {
	// 去掉参数检测
	//if err := opt.perfect(); err != nil {
	//	return nil, err
	//}
	redisOpt := &goRedis.Options{
		Network:         "tcp",
		Addr:            net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port)),
		Password:        opt.Password,
		DB:              opt.DB,
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolSize:        opt.PoolSize,
		//MinIdleConns:       opt.,
		PoolTimeout: opt.PoolTimeout,
		IdleTimeout: opt.IdleTimeout,
		//IdleCheckFrequency: opt.IdleCheckFrequency,
	}
	return newRedis(goRedis.NewClient(redisOpt), log), nil
}

// LogLevel 日志级别
type LogLevel int

type timeKey string

const (
	startKey timeKey = "start-time"
)

func newRedis(client *goRedis.Client, log log.Logger) *Redis {
	client.AddHook(&hook{
		log: log,
	})
	return &Redis{client}
}

type hook struct {
	log   log.Logger
	level LogLevel //nolint
	//slow  time.Duration
}

func (h *hook) BeforeProcess(ctx context.Context, _ goRedis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey, time.Now()), nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd goRedis.Cmder) error {
	h.sweep(ctx, cmd)
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, _ []goRedis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey, time.Now()), nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []goRedis.Cmder) error {
	h.sweepPipeline(ctx, cmds)
	return nil
}

func (h *hook) sweep(ctx context.Context, cmd goRedis.Cmder) {
	use := time.Since(ctx.Value(startKey).(time.Time))
	err := cmd.Err()
	h.log.Info(
		ctx,
		"redis_cmd",
		log.Any("args", cmd.Args()),
		log.Duration("time", use),
		log.ErrorType("err", err),
	)
}

func (h *hook) sweepPipeline(ctx context.Context, cmds []goRedis.Cmder) {
	use := time.Since(ctx.Value(startKey).(time.Time))
	pid := guuid.New().String()

	var firstErr error
	cmdsArgs := make([][]interface{}, 0, len(cmds))
	for _, cmd := range cmds {
		if firstErr == nil {
			if err := cmd.Err(); err != nil && err != goRedis.Nil {
				firstErr = err
			}
		}
		cmdsArgs = append(cmdsArgs, cmd.Args())
	}
	h.log.Info(
		ctx,
		"redislog",
		log.Any("args", cmdsArgs),
		log.Duration("time", use),
		log.String("err", firstErr.Error()),
		log.String("pipeline-id", pid),
	)
}
