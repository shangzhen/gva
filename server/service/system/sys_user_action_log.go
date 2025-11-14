package system

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserActionLogService struct{}

// InitIndex 初始化索引（创建mapping）
func (s *UserActionLogService) InitIndex() error {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	// 检查索引是否存在
	exists, err := client.IndexExists(ctx, indexName)
	if err != nil {
		global.GVA_LOG.Error("检查索引失败", zap.Error(err))
		return err
	}

	if exists {
		global.GVA_LOG.Info("索引已存在", zap.String("index", indexName))
		return nil
	}

	// 创建索引和mapping
	mapping := log.GetESMapping()
	if err := client.CreateIndex(ctx, indexName, mapping); err != nil {
		global.GVA_LOG.Error("创建索引失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("索引创建成功", zap.String("index", indexName))
	return nil
}

// CreateLog 创建日志
func (s *UserActionLogService) CreateLog(req *request.UserActionLogCreate) error {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return fmt.Errorf("ES客户端未初始化")
	}

	// 构建日志对象
	log := system.UserActionLog{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Username:   req.Username,
		Action:     req.Action,
		Module:     req.Module,
		Method:     req.Method,
		Path:       req.Path,
		IP:         req.IP,
		UserAgent:  req.UserAgent,
		Status:     req.Status,
		Latency:    req.Latency,
		Request:    req.Request,
		Response:   req.Response,
		ErrorMsg:   req.ErrorMsg,
		CreateTime: time.Now(),
	}

	// 索引文档
	indexName := log.GetESIndexName()
	if err := client.IndexDocument(ctx, indexName, log.ID, log); err != nil {
		global.GVA_LOG.Error("创建日志失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("日志创建成功", zap.String("id", log.ID))
	return nil
}

// GetLog 获取单条日志
func (s *UserActionLogService) GetLog(id string) (*system.UserActionLog, error) {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return nil, fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	var result system.UserActionLog
	if err := client.GetDocument(ctx, indexName, id, &result); err != nil {
		global.GVA_LOG.Error("获取日志失败", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// SearchLogs 搜索日志
func (s *UserActionLogService) SearchLogs(req *request.UserActionLogSearch) (*response.UserActionLogListResponse, error) {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return nil, fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	// 构建查询条件
	boolQuery := &elasticsearch.BoolQuery{
		Must: []map[string]interface{}{},
	}

	// 用户ID精确匹配
	if req.UserID != nil {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("user_id", *req.UserID))
	}

	// 用户名模糊搜索
	if req.Username != "" {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.MatchQuery("username", req.Username))
	}

	// 操作动作
	if req.Action != "" {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("action", req.Action))
	}

	// 操作模块
	if req.Module != "" {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("module", req.Module))
	}

	// 请求方法
	if req.Method != "" {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("method", req.Method))
	}

	// IP地址
	if req.IP != "" {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("ip", req.IP))
	}

	// 状态码
	if req.Status != nil {
		boolQuery.Must = append(boolQuery.Must, elasticsearch.TermQuery("status", *req.Status))
	}

	// 时间范围
	if req.StartTime != "" || req.EndTime != "" {
		var gte, lte interface{}
		if req.StartTime != "" {
			gte = req.StartTime
		}
		if req.EndTime != "" {
			lte = req.EndTime
		}
		boolQuery.Must = append(boolQuery.Must, elasticsearch.RangeQuery("create_time", gte, lte))
	}

	// 关键词搜索（搜索path和error_msg）
	if req.Keyword != "" {
		keywordQuery := &elasticsearch.BoolQuery{
			Should: []map[string]interface{}{
				elasticsearch.MatchQuery("path", req.Keyword),
				elasticsearch.MatchQuery("error_msg", req.Keyword),
			},
			MinimumShouldMatch: 1,
		}
		boolQuery.Must = append(boolQuery.Must, keywordQuery.ToMap())
	}

	// 如果没有任何查询条件，使用match_all
	var query map[string]interface{}
	if len(boolQuery.Must) == 0 {
		query = elasticsearch.MatchAllQuery()
	} else {
		query = boolQuery.ToMap()
	}

	// 构建排序
	sort := []map[string]interface{}{}
	if req.OrderField != "" {
		orderType := "desc"
		if req.OrderType == "asc" {
			orderType = "asc"
		}
		sort = append(sort, map[string]interface{}{
			req.OrderField: map[string]interface{}{
				"order": orderType,
			},
		})
	} else {
		// 默认按创建时间降序
		sort = append(sort, map[string]interface{}{
			"create_time": map[string]interface{}{
				"order": "desc",
			},
		})
	}

	// 构建搜索请求
	searchReq := &elasticsearch.SearchRequest{
		Query: query,
		From:  (req.Page - 1) * req.PageSize,
		Size:  req.PageSize,
		Sort:  sort,
	}

	// 执行搜索
	result, err := client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索日志失败", zap.Error(err))
		return nil, err
	}

	// 解析结果
	logs := make([]system.UserActionLog, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var log system.UserActionLog
		// 将map转为json再转为struct
		jsonData, _ := json.Marshal(hit.Source)
		if err := json.Unmarshal(jsonData, &log); err != nil {
			global.GVA_LOG.Warn("解析日志失败", zap.Error(err))
			continue
		}
		logs = append(logs, log)
	}

	return &response.UserActionLogListResponse{
		List:     logs,
		Total:    result.Hits.Total.Value,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteLog 删除日志
func (s *UserActionLogService) DeleteLog(id string) error {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	if err := client.DeleteDocument(ctx, indexName, id); err != nil {
		global.GVA_LOG.Error("删除日志失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("日志删除成功", zap.String("id", id))
	return nil
}

// BatchCreateLogs 批量创建日志
func (s *UserActionLogService) BatchCreateLogs(logs []system.UserActionLog) error {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return fmt.Errorf("ES客户端未初始化")
	}

	if len(logs) == 0 {
		return nil
	}

	// 转换为map格式，添加id字段
	documents := make([]map[string]interface{}, 0, len(logs))
	for _, log := range logs {
		if log.ID == "" {
			log.ID = uuid.New().String()
		}
		if log.CreateTime.IsZero() {
			log.CreateTime = time.Now()
		}

		// 转换为map
		jsonData, _ := json.Marshal(log)
		var doc map[string]interface{}
		json.Unmarshal(jsonData, &doc)
		doc["id"] = log.ID // 添加id字段用于批量操作

		documents = append(documents, doc)
	}

	// 批量索引
	indexName := logs[0].GetESIndexName()
	if err := client.BulkIndexDocument(ctx, indexName, documents); err != nil {
		global.GVA_LOG.Error("批量创建日志失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("批量创建日志成功", zap.Int("count", len(logs)))
	return nil
}

// GetStats 获取统计数据
func (s *UserActionLogService) GetStats(req *request.UserActionLogStats) (*response.UserActionLogStatsResponse, error) {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return nil, fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	// 设置默认分组字段
	groupByField := "action"
	if req.GroupBy != "" {
		groupByField = req.GroupBy
	}

	// 构建聚合查询
	aggs := map[string]interface{}{
		"group_stats": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": groupByField,
				"size":  20,
			},
		},
	}

	// 执行聚合
	aggResult, err := client.Aggregate(ctx, indexName, aggs)
	if err != nil {
		global.GVA_LOG.Error("统计查询失败", zap.Error(err))
		return nil, err
	}

	// 解析聚合结果
	stats := make([]map[string]interface{}, 0)
	if groupStats, ok := aggResult["group_stats"].(map[string]interface{}); ok {
		if buckets, ok := groupStats["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					stats = append(stats, b)
				}
			}
		}
	}

	// 获取总数
	var gte, lte interface{}
	if req.StartTime != "" {
		gte = req.StartTime
	}
	if req.EndTime != "" {
		lte = req.EndTime
	}
	query := elasticsearch.RangeQuery("create_time", gte, lte)

	total, err := client.Count(ctx, indexName, query)
	if err != nil {
		global.GVA_LOG.Warn("统计总数失败", zap.Error(err))
		total = 0
	}

	return &response.UserActionLogStatsResponse{
		Total: total,
		Stats: stats,
	}, nil
}

// DeleteIndex 删除索引（危险操作，谨慎使用）
func (s *UserActionLogService) DeleteIndex() error {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		return fmt.Errorf("ES客户端未初始化")
	}

	log := system.UserActionLog{}
	indexName := log.GetESIndexName()

	if err := client.DeleteIndex(ctx, indexName); err != nil {
		global.GVA_LOG.Error("删除索引失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("索引删除成功", zap.String("index", indexName))
	return nil
}
