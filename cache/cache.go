// cache 包提供了缓存扩展
package cache

import (
	"context"
	"io/ioutil"

	"git.zhwenxue.com/zhgo/gocontrib/log"
	"gopkg.in/yaml.v2"
)

type Cache interface {
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, int) error
	Del(context.Context, string) error
	Stats() interface{}
}

type Config struct {
	Cache FreeCacheConfig `yaml:"cache"`
}

func ConfigWithPath(path string) (Config, error) {
	var cnf Config

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return cnf, err
	}

	err = yaml.Unmarshal(f, &cnf)
	if err != nil {
		return cnf, err
	}

	return cnf, nil
}

func NewCache(cnf Config, log *log.Logger) Cache {
	return newFreeCache(cnf.Cache, log)
}
