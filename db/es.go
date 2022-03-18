/*
* @File : elasticsearch
* @Describe :elasticsearch配置项和基本接口声明
* @Author: yangfan@zongheng.com
* @Date : 2022/1/14 10:33
* @Software: GoLand
 */

package db

import (
	"git.zhwenxue.com/zhgo/gocontrib/log"
	"github.com/olivere/elastic"
)
//
//  EsServerOption
//  @Description: elasticsearch的配置项
//
type EsServerOption struct {
	URL         string //Es服务器地址，包含端口号
	Index       string
	Username    string
	Password    string
	Shards      int
	Replicas    int
	Sniff       *bool  //使客户端去嗅探整个集群的状态，true激活,默认false
	HealthCheck *bool
	InfoLog     string
	ErrorLog    string
	TraceLog    string

	// Log 日志对象
	Log log.Logger

	// LogLevel 日志级别
	LogLevel LogLevel

}

//
//  EsAPI
//  @Description: 外部暴露Es接口声明
//
type EsAPI interface {
	esAPI
}

//
//  esAPI
//  @Description: es相关方法注册声明
//
type esAPI interface {

	//
    //  IsDocExists
    //  @Description: 由索引名判断指定索引是否存在
    //  @param index 索引名称 string
    //  @return bool 存在返回true，不存在返回false，并记录日志
    //
	IsDocExistsByIndex(index string) bool

	//
    //  IsDocExistsById
    //  @Description: 由id判断指定索引是否存在
    //  @param id 索引id
    //  @return bool 存在返回true，不存在返回false，并记录日志
    //
	IsDocExistsById(id int) bool

	//
	//  GetDocByIndex
	//  @Description: 由id和index为条件get出相关记录
	//  @receiver es
	//  @param index 索引名 string
	//  @param id id号 int
	//  @return *elastic.GetResult
	//  @return error
	//
	GetDocRowById(index string,id int) (*elastic.GetResult,error)

	//
    //  GetDocsByIds
    //  @Description: 多个索引id批量获取指定index的记录
    //  @param index 索引名
    //  @param ids id数组，uint64
    //  @return *elastic.SearchResult es的搜索结果对象
    //  @return error
    //
	GetDocsByIds(index string,ids []uint64) (*elastic.SearchResult,error)

	//
    //  CreateIndex
    //  @Description: 创建一个指定名称的索引，并指定字段映射,es BodyString 实现
    //  @param index 索引名
    //  @param fieldsMapping 字段映射 json string格式
    //  @return bool 成功创建true 失败false
    //  @return error
    //
    CreateIndex(index string,fieldsMapping string) (bool,error)

	//
    //  InsertRow
    //  @Description: 新插入一条文档记录，如果已有id则更新这条记录
    //  @param index 索引名
    //  @param typeStr 类名
    //  @param id 索引id
    //  @param body 空接口，此处接受新插入的结构体对象
    //  @return string 返回新插入记录的id
    //  @return error
    //
	InsertRow(index string,typeStr string,id string,body interface{}) (string,error)

	//
    //  UpdateRows
    //  @Description: 按条件更新数据
    //  @param index 索引名
    //  @param typeStr 分类名
    //  @param query 查询query接口类型
    //  @param script 更新字段数据接口类型
    //  @return int64 成功返回更新的文件条数，失败返回0
    //  @return error
    //
    UpdateRows(index string,typeStr string,query *elastic.TermQuery,script *elastic.Script) (int64,error)

	//
    //  DelRows
    //  @Description: 删除单个或者多个索引文档
    //  @param index 索引名称
    //  @param typeStr 类型名称
    //  @param query 查询query接口类型，删除条件
    //  @return int64 返回0 没有这条文档，大于0标识删除一个或者多个成功
    //  @return error
    //
	DelRows(index string,typeStr string,query *elastic.TermQuery) (int64,error)

	//
    //  TermSearch
    //  @Description: 精确查找，不分词，不带分析器，如中文词语是无法用term查到的
    //  @param index 索引名
    //  @param typeStr 类型名
    //  @param query 查找条件体TermQuery
    //  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序
    //  @param page 起始元素数
    //  @param limit 偏移量
    //  @return *elastic.SearchResult
    //  @return error
    //
	TermSearch(index string,typeStr string,query *elastic.TermQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error)

	//
	//  MatchSearch
	//  @Description: 分词查找，带分析器，如中文词语是可查到的
	//  @param index 索引名
	//  @param typeStr 类型名
	//  @param query 查找条件体MatchQuery
	//  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序
	//  @param page 起始元素数
	//  @param limit 偏移量
	//  @return *elastic.SearchResult
	//  @return error
	//
	MatchSearch(index string,typeStr string,query *elastic.MatchQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error)

	//
	//  RangeSearch
	//  @Description: 范围查找，带分析器，如中文词语是可查到的
	//  @param index 索引名
	//  @param typeStr 类型名
	//  @param query 范围查找结构体RangeQuery
	//  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序）
	//  @param page 起始元素数
	//  @param limit 偏移量
	//  @return *elastic.SearchResult
	//  @return error
	//
	RangeSearch(index string,typeStr string,query *elastic.RangeQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error)

    //
    //  BoolSearch
    //  @Description: 多条件查询，组合query条件
    //  @param index 索引名
    //  @param typeStr 类名
    //  @param query query条件
    //  @param sort 排序字段（前提是做过sort的字段设置"fielddata": true），否则不要排序）
    //  @param page 起始元素数
    //  @param limit 偏移量
    //  @return *elastic.SearchResult
    //  @return error
    //
	BoolSearch(index string,typeStr string,query *elastic.BoolQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error)
}
