package db

import (
	"git.zhwenxue.com/zhgo/gocontrib/log"
	"time"
)

// Config 主配置文件
type Config struct {
	Redis *RedisConfig `yaml:"redis"`
	//SensorsData *sensorsDataConfig
}

type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	IdleSize     int           `yaml:"idle_size"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ExecSlowTime time.Duration `yaml:"exec_slow_time"`
}

// RedisServerOption redis配置
type RedisServerOption struct {
	Host string
	Port int
	Pwd  string
	DB   int

	// PoolSize 连接池的连接数，默认为16
	PoolSize int

	// IdlePoolSize 空闲连接的数量，默认4
	IdlePoolSize int

	// PoolTimeout 如果连接池所有连接都繁忙，等待获取连接的时间，默认5*time.Second
	PoolTimeout time.Duration

	// IdleTimeout 客户端关闭空闲连接的时间，默认300*time.Second
	IdleTimeout time.Duration

	// IdleCheckFrequency 检查空闲连接的时间，默认60*time.Second
	IdleCheckFrequency time.Duration

	// ExecSlowTime 执行命令未达到超时时间，但是属于"慢"的状态，默认100*time.Millisecond (自定义)
	ExecSlowTime time.Duration

	// Log 日志对象
	Log log.Logger

	// LogLevel 日志级别
	LogLevel LogLevel
}
