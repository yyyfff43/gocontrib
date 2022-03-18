/*
* @File : es_func
* @Describe : elasticsearch方法实现结构体
* @Author: yangfan@zongheng.com
* @Date : 2022/1/14 11:03
* @Software: GoLand
 */

package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	hctx "git.zhwenxue.com/zhgo/gocontrib/context"
	"git.zhwenxue.com/zhgo/gocontrib/log"

	"github.com/olivere/elastic"
	"gopkg.in/yaml.v2"


)

//
//  esdb
//  @Description: es操作结构体，包含获取的客户端和日志
//
type esdb struct {
	esClient *elastic.Client
	log   log.Logger
}

//context获取，使用默认的UUID
var ctx = hctx.GetContext(context.Background(), "")

//
//  EsGroupConfig
//  @Description: 配置文件结构体
//
type EsGroupConfig struct {
	EsMaster EsConfig   `yaml:"esMaster"`
//	EsSlaves []EsConfig `yaml:"esSlaves"`
}

//
//  EsConfigWithPath
//  @Description: 获取配置文件信息并映射成结构体
//  @param path 配置文件yml路径
//  @return EsGroupConfig
//  @return error
//
func EsConfigWithPath(path string) (EsGroupConfig, error) {
	esConfig := EsGroupConfig{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return esConfig, err
	}

	err = yaml.Unmarshal(yamlFile, &esConfig)
	if err != nil {
		return esConfig, err
	}

	return esConfig, nil

}

type EsConfig struct {
	DataSource   string `yaml:"dataSource"`   //es服务器地址
}

//
//  NewEs
//  @Description: 初始化一个ElasticSearch客户端
//  @param esOption es配置项参数
//  @param log 日志
//  @return EsAPI 返回es实例化接口
//  @return error 错误信息
//
func NewEs(esOption *EsServerOption) (EsAPI, error) {
	var host = esOption.URL
	var err error
	var startTime = time.Now()//起始调用时间
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		// 设置错误日志输出(会在控制台输出)
		//      elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出(会在控制台输出)
		//		elastic.SetInfoLog(log.New(os.Stdout, "info:", log.LstdFlags)),
		// 设置trace日志输出(会在控制台输出，包括查询的dsl语句，跟踪过程等)
		//		elastic.SetTraceLog(log.New(os.Stdout, "trace", log.LstdFlags)),
	)
//	var esLog = elastic.Logger() //是否需要另加es自己的log?
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		//服务端连接出错
		esOption.Log.Error(ctx, err.Error())
	}

	info, code, err := client.Ping(host).Do(ctx)
	if err != nil {
		//服务端ping命令后报错
		esOption.Log.Error(ctx, "elasticsearch服务端ping命令报错:"+err.Error())
	}else {
		esOption.Log.Info(ctx, "Elasticsearch returned with code "+strconv.Itoa(code)+"and version "+info.Version.Number,
			log.String("状态码：", strconv.Itoa(code)),
			log.String("版本号：", info.Version.Number),
			log.Duration("time", usedTime),
		)
	}

	return newEs(client,esOption.Log),err
}

//
//  newEs
//  @Description: es客户端工厂实例化结构体
//  @param esClient
//  @param log
//  @return EsAPI
//
func newEs(esClient *elastic.Client, log  log.Logger) EsAPI {
	return &esdb{esClient: esClient, log: log}
}

//
//  GetEsClient
//  @Description: 获取一个es的客户端实例，注意：此实例需自己打日志跟踪并且使用完毕后要关闭es链接:defer es.esClient.Stop()
//  @param esOption
//  @return *elastic.Client
//  @return error
//
func GetEsClient(esOption *EsServerOption) (*elastic.Client,error) {
	var host = esOption.URL
	var err error
	var startTime = time.Now()//起始调用时间
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		// 设置错误日志输出(会在控制台输出)
		//      elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出(会在控制台输出)
		//		elastic.SetInfoLog(log.New(os.Stdout, "info:", log.LstdFlags)),
		// 设置trace日志输出(会在控制台输出，包括查询的dsl语句，跟踪过程等)
		//		elastic.SetTraceLog(log.New(os.Stdout, "trace", log.LstdFlags)),
	)
	//	var esLog = elastic.Logger() //是否需要另加es自己的log?
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		//服务端连接出错
		esOption.Log.Error(ctx, err.Error())
	}

	info, code, err := client.Ping(host).Do(ctx)
	if err != nil {
		//服务端ping命令后报错
		esOption.Log.Error(ctx, "elasticsearch服务端ping命令报错:"+err.Error())
	}else {
		esOption.Log.Info(ctx, "Elasticsearch returned with code "+strconv.Itoa(code)+"and version "+info.Version.Number,
			log.String("状态码：", strconv.Itoa(code)),
			log.String("版本号：", info.Version.Number),
			log.Duration("time", usedTime),
		)
	}

	return client,err
}

//
//  IsDocExistsByIndex
//  @Description: 判断指定索引是否存在
//  @receiver es
//  @param index 索引名称
//  @return bool 存在返回true，不存在返回false，并记录日志
//
func (es *esdb) IsDocExistsByIndex(index string) bool {
	defer es.esClient.Stop()//函数退出时关闭es连接
	var startTime = time.Now()//起始调用时间
	exist,_ := es.esClient.IndexExists(index).Do(ctx)
	usedTime := time.Since(startTime)//结束调用时间
	if !exist{
		es.log.Warn(ctx,"索引：“"+index+"“不存在",log.Duration("time",usedTime))
//		log.Println("ID may be incorrect! ",id)
		return false
	}else{
		es.log.Info(ctx,index+"index索引存在",log.Duration("time",usedTime))
	}
	return true
}

//
//  IsDocExistsById
//  @Description: 判断指定索引是否存在
//  @receiver es
//  @param id 索引id
//  @return bool 存在返回true，不存在返回false，记录日志
//
func (es *esdb) IsDocExistsById(id int) bool{
	defer es.esClient.Stop()//函数退出时关闭es连接
	var startTime = time.Now()//起始调用时间
	exist,_ := es.esClient.Exists().Id(strconv.Itoa(id)).Do(ctx)
	usedTime := time.Since(startTime)//结束调用时间
	if !exist{
		es.log.Warn(ctx,"_id为：“"+strconv.Itoa(id)+"“的索引不存在",log.Duration("time",usedTime))
		return false
	}else{
		es.log.Info(ctx,strconv.Itoa(id)+"id索引存在",log.Duration("time",usedTime))
	}
	return true
}

//
//  GetDocByIndex
//  @Description: 由id和index为条件get出相关记录
//  @receiver es
//  @param index 索引名 string
//  @param id id号 int
//  @return *elastic.GetResult
//  @return error
//
func (es *esdb) GetDocRowById(index string,id int) (*elastic.GetResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	// 使用文档id查询
	var startTime = time.Now()//起始调用时间
	getRes, err := es.esClient.Get().Index(index).Id(strconv.Itoa(id)).
		Pretty(true).//查询结果返回较好的json格式，但好像不加这个照样MarshalJSON能处理
		Do(ctx)
	fmt.Println(getRes)
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil{
		es.log.Error(
			ctx,
			"index为：“"+index+"“ id为"+strconv.Itoa(id)+"的索引不存在",
			log.Duration("time",usedTime),
			)
	}else{
		es.log.Info(ctx,"GetDocRowById获取单一id记录成功",log.Duration("time",usedTime))
	}

	return getRes,err
}

//
//  GetDocsByIds
//  @Description: 多个索引id批量获取指定index的记录
//  @param index 索引名
//  @param ids id数组，uint64
//  @return *elastic.SearchResult es的搜索结果对象
//  @return error
//
func (es *esdb) GetDocsByIds(index string,ids []uint64)(*elastic.SearchResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	idStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idStr = append(idStr, strconv.FormatUint(id, 10))
	}
	var startTime = time.Now()//起始调用时间
	resp, err := es.esClient.Search().Index(index).Query(
		elastic.NewIdsQuery().Ids(idStr...)).Size(len(ids)).Do(ctx)
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		es.log.Error(
			ctx,
			"索引为：“"+index+"“的多id查询报错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
			)
	}else{
		es.log.Info(ctx,"GetDocsByIds批量获取记录成功",log.Duration("time",usedTime))
	}

	return resp,err
}

//
//  CreateIndex
//  @Description: 创建一个指定名称的索引，并指定字段映射，创建一个指定名称的索引，并指定字段映射,es BodyString 实现
//  @param index 索引名
//  @param fieldsMapping 字段映射
//  @return bool 成功创建true 失败false
//  @return error
//
func (es *esdb) CreateIndex(index string,fieldsMapping string) (bool,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	var flag bool = false
	// 首先检测下索引是否存在
	exists, err := es.esClient.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		es.log.Error(ctx,"检测索引为：“"+index+"“是否存在时出错",log.String("错误：", err.Error()))
	}
	if !exists {
		var startTime = time.Now()//起始调用时间
		// 传入索引不存在，则创建一个
		_, err := es.esClient.CreateIndex(index).BodyString(fieldsMapping).Do(ctx)
		usedTime := time.Since(startTime)//结束调用时间
		if err != nil {
			// Handle error
			es.log.Error(ctx,
				"创建索引为：“"+index+"“时出错",
				log.String("错误：", err.Error()),
				log.Duration("time",usedTime),
			)
		}else{
			flag = true
			es.log.Info(ctx,"创建索引为：“"+index+"“成功",log.Duration("time",usedTime))
		}
	}
	return flag,err
}

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
func (es *esdb) InsertRow(index string,typeStr string,id string,body interface{}) (string,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	var startTime = time.Now()//起始调用时间
	// 使用client创建一个新的文档
	put1, err := es.esClient.Index().
		Index(index). // 设置索引名称
		Type(typeStr).
		Id(id).
		BodyJson(body). // 指定前面声明的索引内容
		Do(ctx)         // 执行请求，需要传入一个上下文对象
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"添加索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
		)
	}else{
		es.log.Info(
			ctx,
			"添加索引文档为：“"+index+"“成功",
			log.Duration("time",usedTime),
			)
	}

    return put1.Id,err
}

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
func (es *esdb)UpdateRows(index string,typeStr string,query *elastic.TermQuery,script *elastic.Script) (int64,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	var startTime = time.Now()//起始调用时间
	res, err := es.esClient.UpdateByQuery(index).
		Type(typeStr).
		// 设置查询条件，TermQuery 结构体
		Query(
			query,
		).
		// 通过脚本更新内容，Script 结构体
		Script(
			script,
		).
		Do(ctx)
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"更新索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
		)
	}else{
		if res.Total>0{
			es.log.Info(
				ctx,
				"更新索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
			)
		}else{
			es.log.Info(
				ctx,
				"要更新的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
			)
		}
	}
    return res.Total,err
}

//
//  DelRows
//  @Description: 删除单个或者多个索引文档
//  @param index 索引名称
//  @param typeStr 类型名称
//  @param query 查询query接口类型，删除条件
//  @return int64 返回0 没有这条文档，大于0标识删除一个或者多个成功
//  @return error
//
func (es *esdb)DelRows(index string,typeStr string,query *elastic.TermQuery) (int64,error)  {
	defer es.esClient.Stop()//函数退出时关闭es连接
	var startTime = time.Now()//起始调用时间
	res, err := es.esClient.DeleteByQuery(index). // 设置索引名
		Type(typeStr).
		// 设置查询条件，TermQuery 结构体
		Query(
			query,
		).
		// 文档冲突也继续删除
		ProceedOnVersionConflict().
		Do(ctx)
	usedTime := time.Since(startTime)//结束调用时间
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"删除索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
		)
	}else{
		if res.Total>0{
			es.log.Info(
				ctx,
				"删除索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
			)
		}else{
			es.log.Info(
				ctx,
				"要删除的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
			)
		}
	}
	return res.Total,err
}

//
//  TermSearch
//  @Description: 精确查找，不分词，不带分析器，如中文词语是无法用term查到的
//  @param index 索引名
//  @param typeStr 类型名
//  @param query 查找条件体
//  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序
//  @param page 起始元素数
//  @param limit 偏移量
//  @return *elastic.SearchResult
//  @return error
//
func (es *esdb)TermSearch(index string,typeStr string,termQuery *elastic.TermQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
    //var ascending = true
    //if sort=="desc"{
	//	ascending = false
	//}
	var startTime = time.Now()//起始调用时间
	searchService := es.esClient.Search().
		Index(index).          // 设置索引名
		Type(typeStr).         //设置类型名
		Query(termQuery).      // 设置查询条件
		//				Sort("first_name", ascending). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(page).          // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(limit).          // 设置分页参数 - 每页大小
		Pretty(true)    // 查询结果返回可读性较好的JSON格式
	//设置排序，默认倒序DESC
	if sort!=""{
		if ascending == "ASC"{
			searchService.Sort(sort,true)
		}else{
			searchService.Sort(sort,false)
		}
	}
	searchRes, err := searchService.Do(ctx)                 // 执行请求
	usedTime := time.Since(startTime)//结束调用时间
	/*
	 * 如果报Error 400 (Bad Request): all shards failed [type=search_phase_execution_exception]
	 * 则sort的字段没设置"fielddata": true
	 */
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"term精确检索索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
//			log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
		)
	}else{
		if searchRes.TotalHits()>0{
			es.log.Info(
				ctx,
				"索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}else{
			es.log.Info(
				ctx,
				"要检索的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}
	}
	return searchRes, err
}

//
//  MatchSearch
//  @Description: 分词查找，带分析器，如中文词语是可查到的
//  @param index 索引名
//  @param typeStr 类型名
//  @param query 查找条件体
//  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序
//  @param page 起始元素数
//  @param limit 偏移量
//  @return *elastic.SearchResult
//  @return error
//
func (es *esdb)MatchSearch(index string,typeStr string,matchQuery *elastic.MatchQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	//var ascending = true
	//if sort=="desc"{
	//	ascending = false
	//}
	var startTime = time.Now()//起始调用时间
	searchService := es.esClient.Search().
		Index(index).          // 设置索引名
		Type(typeStr).         //设置类型名
		Query(matchQuery).      // 设置查询条件
		//				Sort("first_name", ascending). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(page).          // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(limit).          // 设置分页参数 - 每页大小
		Pretty(true)    // 查询结果返回可读性较好的JSON格式
	//设置排序，默认倒序DESC
	if sort!=""{
		if ascending == "ASC"{
			searchService.Sort(sort,true)
		}else{
			searchService.Sort(sort,false)
		}
	}
	searchRes, err := searchService.Do(ctx)                 // 执行请求
	usedTime := time.Since(startTime)//结束调用时间
	/*
	 * 如果报Error 400 (Bad Request): all shards failed [type=search_phase_execution_exception]
	 * 则sort的字段没设置"fielddata": true
	 */
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"term精确检索索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
			//			log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
		)
	}else{
		if searchRes.TotalHits()>0{
			es.log.Info(
				ctx,
				"索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}else{
			es.log.Info(
				ctx,
				"要检索的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}
	}
	return searchRes, err
}

//
//  RangeSearch
//  @Description: 范围查找，带分析器，如中文词语是可查到的
//  @param index 索引名
//  @param typeStr 类型名
//  @param query 范围查找结构体RangeQuery
//  @param sort 排序条件（前提是做过sort的字段设置"fielddata": true），否则不要排序
//  @param page 起始元素数
//  @param limit 偏移量
//  @return *elastic.SearchResult
//  @return error
//
func (es *esdb)RangeSearch(index string,typeStr string,rangeQuery *elastic.RangeQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	//var ascending = true
	//if sort=="desc"{
	//	ascending = false
	//}
	var startTime = time.Now()//起始调用时间
	searchService := es.esClient.Search().
		Index(index).          // 设置索引名
		Type(typeStr).         //设置类型名
		Query(rangeQuery).      // 设置查询条件
		//				Sort("first_name", ascending). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(page).          // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(limit).          // 设置分页参数 - 每页大小
		Pretty(true)    // 查询结果返回可读性较好的JSON格式
	//设置排序，默认倒序DESC
	if sort!=""{
		if ascending == "ASC"{
			searchService.Sort(sort,true)
		}else{
			searchService.Sort(sort,false)
		}
	}
	searchRes, err := searchService.Do(ctx)              // 执行请求
	usedTime := time.Since(startTime)//结束调用时间
	/*
	 * 如果报Error 400 (Bad Request): all shards failed [type=search_phase_execution_exception]
	 * 则sort的字段没设置"fielddata": true
	 */
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"term精确检索索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
			//			log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
		)
	}else{
		if searchRes.TotalHits()>0{
			es.log.Info(
				ctx,
				"索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}else{
			es.log.Info(
				ctx,
				"要检索的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}
	}
	return searchRes, err
}

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
func (es *esdb)BoolSearch(index string,typeStr string,query *elastic.BoolQuery,sort string,ascending string,page int,limit int) (*elastic.SearchResult,error){
	defer es.esClient.Stop()//函数退出时关闭es连接
	//var ascending = true
	//if sort=="desc"{
	//	ascending = false
	//}
	var startTime = time.Now()//起始调用时间
	searchService := es.esClient.Search().
		Index(index).          // 设置索引名
		Type(typeStr).         //设置类型名
		Query(query).      // 设置查询条件
		//				Sort("first_name", ascending). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(page).          // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(limit).          // 设置分页参数 - 每页大小
		Pretty(true)    // 查询结果返回可读性较好的JSON格式
		//设置排序，默认倒序DESC
		if sort!=""{
			if ascending == "ASC"{
				searchService.Sort(sort,true)
			}else{
				searchService.Sort(sort,false)
			}
		}
		searchRes, err := searchService.Do(ctx)
		              // 执行请求
	usedTime := time.Since(startTime)//结束调用时间
	/*
	 * 如果报Error 400 (Bad Request): all shards failed [type=search_phase_execution_exception]
	 * 则sort的字段没设置"fielddata": true
	 */
	if err != nil {
		// Handle error
		es.log.Error(ctx,
			"term精确检索索引文档为：“"+index+"“时出错",
			log.String("错误：", err.Error()),
			log.Duration("time",usedTime),
			//			log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
		)
	}else{
		if searchRes.TotalHits()>0{
			es.log.Info(
				ctx,
				"索引文档：“"+index+"“成功",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}else{
			es.log.Info(
				ctx,
				"要检索的索引文档：“"+index+"“不存在",
				log.Duration("time",usedTime),
				//				log.String("es内部搜索耗时:",strconv.FormatInt(searchRes.TookInMillis,64)),
			)
		}
	}
	return searchRes, err
}