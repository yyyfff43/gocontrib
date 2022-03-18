/*
* @File : mysql_test
* @Describe :
* @Author: zhangnaiqian@zongheng.com
* @Date : 2022/1/14 18:25
* @Software: GoLand
 */

package db

import (
	"context"
	"fmt"
	hctx "git.zhwenxue.com/zhgo/gocontrib/context"
	"git.zhwenxue.com/zhgo/gocontrib/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMysqlConfigWithPath(t *testing.T) {
	dir, _ := os.Getwd()
	configPath := dir + "/mysql_config_sample.yml"
	mysqlConfig, err := MysqlConfigWithPath(configPath)
	assert.Nil(t, err)
	fmt.Println(mysqlConfig)
}

func TestMysqlClient(t *testing.T) {
	dir, _ := os.Getwd()
	configPath := dir + "/mysql_config_sample.yml"
	mysqlConfig, err := MysqlConfigWithPath(configPath)
	assert.Nil(t, err)
	logger := log.New(os.Stdout, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	engine, err := NewMysqlClient(mysqlConfig, *logger)
	assert.Nil(t, err)
	type hachi struct {
		Id     int
		Test_a string
		Test_b string
	}
	myHachi := new(hachi)
	ctx := hctx.GetContext(context.Background(), "")
	has, err := engine.Context(ctx).Table("hachi").Get(myHachi)
	assert.Nil(t, err)
	fmt.Println(has)
	fmt.Println(myHachi)
}
