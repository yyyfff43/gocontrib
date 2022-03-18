/*
* @File : es_test
* @Describe :es接口单元测试
* @Author: yangfan@zongheng.com
* @Date : 2022/1/14 16:57
* @Software: GoLand
 */

package db

import (
	"encoding/json"
	"fmt"
	"git.zhwenxue.com/zhgo/gocontrib/log"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type ForumsAll struct{
	ImgUrl     string   `json:"imgUrl"`
	IndexTime  int64    `json:"indexTime"`
	Name       string   `json:"name"`
	Id         int64    `json:"id"`
	Title      string   `json:"title"`
	Type       int64    `json:"type"`
	BookId     int64    `json:"bookId"`
}

// 索引mapping定义，新建一个索引时使用,包括自定义的字段等,json字符串,对应forums_all索引
const MappingForumsAll = `
  "settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
  },
  "mappings":{
  "all": {
    "properties": {
      "imgUrl": {
        "type": "keyword",
        "store": true
      },
      "indexTime": {
        "type": "long"
      },
       "name": {
        "type": "keyword",
        "store": true
      },
       "id": {
	"type": "long",
	"store": true
      },
      "title": {
        "type": "text",
	"store": true,
	"index_options": "offsets",
        "analyzer": "ik_max_word"
      },
     "type": {
	"type": "long",
	"store": true
     },
     "bookId": {
	"type": "long",
	"store": true
     }
     }
    }
  }`

var file, _ = os.OpenFile("./hachi.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
var logger = log.New(file, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))

//测试初始化客户端
func TestNewEs(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
    //读取配置文件信息
	dir, _ := os.Getwd()
	configPath := dir + "/es_config.yml"
	esConfig, err := EsConfigWithPath(configPath)
	assert.Nil(t, err)

	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         esConfig.EsMaster.DataSource ,
		Log:                *logger,
	}

	_,err = NewEs(EsServerOption)
	assert.Nil(t, err)
}

//测试获取es访问客户端
func TestEsdb_GetEsClient(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//读取配置文件信息
	dir, _ := os.Getwd()
	configPath := dir + "/es_config.yml"
	esConfig, err := EsConfigWithPath(configPath)
	assert.Nil(t, err)
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         esConfig.EsMaster.DataSource,
		Log:                *logger,
	}
	_, err = GetEsClient(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
}

//测试某个指定索引index是否存在
func TestEsdb_IsDocExistsByIndex(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//读取配置文件信息
	dir, _ := os.Getwd()
	configPath := dir + "/es_config.yml"
	esConfig, err := EsConfigWithPath(configPath)
	assert.Nil(t, err)

	EsServerOption := &EsServerOption{
		URL:         esConfig.EsMaster.DataSource,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//测试某个指定索引index是否存在
	resFlag := client.IsDocExistsByIndex("forums_all")
	if resFlag {
		fmt.Println("索引存在")
	}else{
		fmt.Println("索引不存在")
	}
}

//测试某个id的索引是否存在
func TestEsdb_IsDocExistsById(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//测试某个id的索引是否存在
	resFlag := client.IsDocExistsById(88)
	if resFlag {
		fmt.Println("索引存在")
	}else{
		fmt.Println("索引不存在")
	}
}

//测试获取某个指定索引的结果
func TestEsdb_GetDocRowById(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}

	var getRes,err2 = client.GetDocRowById("ForumsAll",4)
	if err2!=nil {
		fmt.Println(err2.Error())
	}else{
		if getRes.Found{
			t := ForumsAll{}
			resSourceJson,_ := getRes.Source.MarshalJSON()//如果GetDocRowById的Get调用加载了Pretty(true).
			_ = json.Unmarshal(resSourceJson,&t)
			fmt.Println(t)
		}
	}
}

//测试获取多个id的指定某索引下记录
func TestEsdb_GetDocsByIds(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}

	ids := []uint64{1,2}
	var tArr []ForumsAll
	var getRess,err3 = client.GetDocsByIds("ForumsAll",ids)
	if err3==nil {
		if getRess.TotalHits() == 0 {
			fmt.Println("未查找到任何记录")
		} else {
			for _, e := range getRess.Each(reflect.TypeOf(ForumsAll{})) {
				us := e.(ForumsAll)
				tArr = append(tArr, us)
			}
			for _, vForumsAll := range tArr {
				fmt.Println(vForumsAll)
			}
		}
	}else{
		fmt.Println(err3.Error())
	}
}

//测试创建一个新的索引，如果索引存在则返回已创建了
func TestEsdb_CreateIndex(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	_,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}

	//var cRes,err4 = client.CreateIndex("forums_all",MappingForumsAll)
	//if err4==nil {
	//	if cRes {
	//		fmt.Println("新索引创建成功")
	//	}else{
	//		fmt.Println(err4.Error())
	//	}
	//}else{
	//	fmt.Println(err4.Error())
	//}
}

//测试新插入一条文档
func TestEsdb_InsertRow(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}

	tStruct := ForumsAll{ImgUrl: "/cover/2009/08/1251427811142.jpg",IndexTime: 1639656000001,Name: "初平纪事",Id: 11,Title: "初平纪事",Type: 0,BookId: 2223}
//	strStuct := `{"properties": {"age": {"fielddata": true}}}`

	var cRes,err5 = client.InsertRow("forums_all","all","11",tStruct)
	cResInt,_ := strconv.Atoi(cRes)
	if err5==nil {
		if cResInt > 0 {
			fmt.Println("插入新文档成功")
		}else{
			fmt.Println("插入新文档失败")
		}
	}else{
		fmt.Println(err5.Error())
	}
}

//测试条件更新多条记录
func TestEsdb_UpdateRows(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
    //查找条件
	termQuery := elastic.NewTermQuery("age", "3")
	//要修改的字段
	termScript := elastic.NewScript("ctx._source['about']=6")

	var cRes,_ = client.UpdateRows("forums_all","all",termQuery,termScript)
	if cRes>0 {
		fmt.Println("更新新文档成功")
	}else{
		fmt.Println("更新失败")
	}
}

//测试条件删除多条记录
func TestEsdb_DelRows(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//查找条件
	termQuery := elastic.NewTermQuery("_id", "1")

	var cRes,_ = client.DelRows("forums_all","all",termQuery)
	if cRes>0 {
		fmt.Println("删除文档成功")
	}else{
		fmt.Println("删除失败")
	}
}

//测试精确Term查询
func TestEsdb_TermSearch(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//查找条件,term不带中文分析器，只能解析单个汉字，词语无法解析
	termQuery := elastic.NewTermQuery(
		"name", "初平纪事",
	)

	var tArr []ForumsAll
	var getRess,err3 = client.TermSearch("forums_all","all",termQuery,"age","",0,10)
	if err3==nil {
		if getRess.TotalHits() == 0 {
			fmt.Println("未查找到任何记录")
		} else {
			var tStruct ForumsAll
			for _, item := range getRess.Each(reflect.TypeOf(tStruct)) {
				// 转换成ForumsAll对象
				if t, ok := item.(ForumsAll); ok {
					tArr = append(tArr, t)
				}
			}
            fmt.Println(tArr)
			for _, vForumsAll := range tArr {
				fmt.Println(vForumsAll)
			}
		}
	}else{
		fmt.Println(err3.Error())
	}
}

//测试matchTerm查询
func TestEsdb_MatchSearch(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//查找条件,matchQuery带中文分析器
	matchQuery := elastic.NewMatchQuery("first_name","关羽")

	var tArr []ForumsAll
	var getRess,err3 = client.MatchSearch("ForumsAll","languages",matchQuery,"age","",0,10)
	if err3==nil {
		if getRess.TotalHits() == 0 {
			fmt.Println("未查找到任何记录")
		} else {
			var tStruct ForumsAll
			for _, item := range getRess.Each(reflect.TypeOf(tStruct)) {
				// 转换成ForumsAll对象
				if t, ok := item.(ForumsAll); ok {
					tArr = append(tArr, t)
				}
			}
			fmt.Println(tArr)
			for _, vForumsAll := range tArr {
				fmt.Println(vForumsAll)
			}
		}
	}else{
		fmt.Println(err3.Error())
	}
}

//测试范围查询
func TestEsdb_RangeSearch(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	//查找条件,rangeQuery带中文分析器
	rangeQuery := elastic.NewRangeQuery("about").
		Gte(300).
		Lte(900)

	var tArr []ForumsAll
	var getRess,err3 = client.RangeSearch("ForumsAll","languages",rangeQuery,"desc","",0,10)
	if err3==nil {
		if getRess.TotalHits() == 0 {
			fmt.Println("未查找到任何记录")
		} else {
			var tStruct ForumsAll
			for _, item := range getRess.Each(reflect.TypeOf(tStruct)) {
				// 转换成ForumsAll对象
				if t, ok := item.(ForumsAll); ok {
					tArr = append(tArr, t)
				}
			}
			fmt.Println(tArr)
			for _, vForumsAll := range tArr {
				fmt.Println(vForumsAll)
			}
		}
	}else{
		fmt.Println(err3.Error())
	}
}

//测试多条件查询
func TestEsdb_BoolSearch(t *testing.T) {
	defer func() {
		_ = logger.Sync()//清空log缓冲区，测试结束前应defer调用
	}()
	//测试初始化客户端
	EsServerOption := &EsServerOption{
		URL:         "http://10.3.138.104:9200" ,
		Log:                *logger,
	}

	client,err := NewEs(EsServerOption)
	if err!=nil {
		println(err.Error())
	}
	// 创建bool查询
	boolQuery := elastic.NewBoolQuery().Must()

	// 创建term和match查询
	termQuery := elastic.NewTermQuery("about", "822")
	matchQuery := elastic.NewMatchQuery("first_name", "关羽")
    //可以再增加条件查询
//	rangeQuery := elastic.NewRangeQuery("id").Gte(8).Lte(9)

	// 设置bool查询的must条件, 组合了两个子查询
	boolQuery.Must(termQuery, matchQuery)

	var tArr []ForumsAll
	var getRess,err3 = client.BoolSearch("toruk","languages",boolQuery,"age","ASC",0,10)
	if err3==nil {
		if getRess.TotalHits() == 0 {
			fmt.Println("未查找到任何记录")
		} else {
			var tStruct ForumsAll
			for _, item := range getRess.Each(reflect.TypeOf(tStruct)) {
				// 转换成ForumsAll对象
				if t, ok := item.(ForumsAll); ok {
					tArr = append(tArr, t)
				}
			}
			fmt.Println(tArr)
			for _, vForumsAll := range tArr {
				fmt.Println(vForumsAll)
			}
		}
	}else{
		fmt.Println(err3.Error())
	}
}

//测试获取es服务器配置信息
func TestEsConfigWithPath(t *testing.T) {
	dir, _ := os.Getwd()
	configPath := dir + "/es_config.yml"
	esConfig, err := EsConfigWithPath(configPath)
	assert.Nil(t, err)
	fmt.Println(esConfig)
}