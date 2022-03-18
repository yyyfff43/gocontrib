/*
* @File : mysql
* @Describe :
* @Author: zhangnaiqian@zongheng.com
* @Date : 2022/1/14 15:55
* @Software: GoLand
 */

package db

import (
	"context"
	"git.zhwenxue.com/zhgo/gocontrib/log"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
)

type MysqlGroupConfig struct {
	MysqlMaster MysqlConfig   `yaml:"mysqlMaster"`
	MysqlSlaves []MysqlConfig `yaml:"mysqlSlaves"`
}

func MysqlConfigWithPath(path string) (MysqlGroupConfig, error) {
	mysqlConfig := MysqlGroupConfig{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return mysqlConfig, err
	}

	err = yaml.Unmarshal(yamlFile, &mysqlConfig)
	if err != nil {
		return mysqlConfig, err
	}

	return mysqlConfig, nil

}

type MysqlConfig struct {
	DataSource   string `yaml:"dataSource"`   // mysql账号密码地址 root:123@localhost/test?charset=utf8
	MaxIdleConns int    `yaml:"maxIdleConns"` // 连接池的空闲数大小
	MaxOpenConns int    `yaml:"maxOpenConns"` // 最大打开连接数
}

// NewMysqlClient
// @Description: 初始化单一的mysql客户端
// @param config
// @return *xorm.Engine
// @return error
func NewMysqlClient(config MysqlGroupConfig, log log.Logger) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", config.MysqlMaster.DataSource)
	if err != nil {
		return nil, err
	}

	setMysqlConfig(engine, config.MysqlMaster)
	engine.AddHook(NewTracingHook(log))
	return engine, nil
}

func NewMysqlGroup(config MysqlGroupConfig, log log.Logger) (*xorm.EngineGroup, error) {
	master, err := xorm.NewEngine("mysql", config.MysqlMaster.DataSource)
	if err != nil {
		return nil, err
	}

	setMysqlConfig(master, config.MysqlMaster)
	master.AddHook(NewTracingHook(log))

	var slaves []*xorm.Engine
	slaveNum := len(config.MysqlSlaves)
	for i := 0; i < slaveNum; i++ {
		slave, err := xorm.NewEngine("mysql", config.MysqlSlaves[i].DataSource)
		if err != nil {
			return nil, err
		}

		setMysqlConfig(slave, config.MysqlSlaves[i])
		slave.AddHook(NewTracingHook(log))
		slaves = append(slaves, slave)
	}

	eg, err := xorm.NewEngineGroup(master, slaves)
	return eg, err
}

func setMysqlConfig(engine *xorm.Engine, config MysqlConfig) {
	if config.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(config.MaxIdleConns)
	}

	if config.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(config.MaxOpenConns)
	}
}

type TracingHook struct {
	// 注意Hook伴随DB实例的生命周期，所以我们不能在Hook里面寄存span变量
	// 否则就会发生并发问题
	Log    log.Logger
	before func(c *contexts.ContextHook) (context.Context, error)
	after  func(c *contexts.ContextHook) error
}

// xorm的hook接口需要满足BeforeProcess和AfterProcess函数
func (h *TracingHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	c.Ctx = context.WithValue(c.Ctx, startKey, time.Now())
	return h.before(c)
}

func (h *TracingHook) AfterProcess(c *contexts.ContextHook) error {
	sql, _ := builder.ConvertToBoundSQL(c.SQL, c.Args)
	use := time.Since(c.Ctx.Value(startKey).(time.Time))
	h.Log.Info(c.Ctx, "SQL",
		log.String("sql", sql),
		log.Any("args", c.Args),
		log.Duration("time", use),
	)
	return h.after(c)
}

// 让编译器知道这个是xorm的Hook，防止编译器无法检查到异常
var _ contexts.Hook = &TracingHook{}

// 使用新定义类型
// context.WithValue方法中注释写到
// 提供的键需要可比性，而且不能是字符串或者任意内建类型，避免不同包之间
// 调用到相同的上下文Key发生碰撞，context的key具体类型通常为struct{}，
// 或者作为外部静态变量（即开头字母为大写）,类型应该是一个指针或者interface类型
func before(c *contexts.ContextHook) (context.Context, error) {
	return c.Ctx, nil
}

func after(c *contexts.ContextHook) error {
	return nil
}

func NewTracingHook(log log.Logger) *TracingHook {
	return &TracingHook{
		Log:    log,
		before: before,
		after:  after,
	}
}
