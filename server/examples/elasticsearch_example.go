package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"
	"go.uber.org/zap"
)

// ESExample ES 使用示例
type ESExample struct{}

// 示例文档结构
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

// ExampleBasicOperations 基础操作示例
func (e *ESExample) ExampleBasicOperations() {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		global.GVA_LOG.Error("ES客户端未初始化")
		return
	}

	indexName := "products"

	// 1. 创建索引
	global.GVA_LOG.Info("=== 创建索引 ===")
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "text",
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"price": map[string]interface{}{
					"type": "float",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"stock": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	if err := client.CreateIndex(ctx, indexName, mapping); err != nil {
		global.GVA_LOG.Error("创建索引失败", zap.Error(err))
	}

	// 2. 索引文档（创建）
	global.GVA_LOG.Info("=== 索引文档 ===")
	product := Product{
		ID:          "1",
		Name:        "iPhone 15 Pro",
		Description: "最新款苹果手机，A17 Pro芯片",
		Price:       7999.00,
		Category:    "手机",
		Stock:       100,
		CreatedAt:   time.Now(),
	}

	if err := client.IndexDocument(ctx, indexName, product.ID, product); err != nil {
		global.GVA_LOG.Error("索引文档失败", zap.Error(err))
	}

	// 3. 获取文档
	global.GVA_LOG.Info("=== 获取文档 ===")
	var retrievedProduct Product
	if err := client.GetDocument(ctx, indexName, "1", &retrievedProduct); err != nil {
		global.GVA_LOG.Error("获取文档失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("获取到的产品", zap.Any("product", retrievedProduct))
	}

	// 4. 更新文档
	global.GVA_LOG.Info("=== 更新文档 ===")
	update := map[string]interface{}{
		"price": 7499.00,
		"stock": 150,
	}
	if err := client.UpdateDocument(ctx, indexName, "1", update); err != nil {
		global.GVA_LOG.Error("更新文档失败", zap.Error(err))
	}

	// 5. 删除文档
	global.GVA_LOG.Info("=== 删除文档 ===")
	if err := client.DeleteDocument(ctx, indexName, "1"); err != nil {
		global.GVA_LOG.Error("删除文档失败", zap.Error(err))
	}

	// 6. 删除索引
	global.GVA_LOG.Info("=== 删除索引 ===")
	if err := client.DeleteIndex(ctx, indexName); err != nil {
		global.GVA_LOG.Error("删除索引失败", zap.Error(err))
	}
}

// ExampleBulkOperations 批量操作示例
func (e *ESExample) ExampleBulkOperations() {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		global.GVA_LOG.Error("ES客户端未初始化")
		return
	}

	indexName := "products"

	// 创建索引
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"name":     map[string]interface{}{"type": "text"},
				"price":    map[string]interface{}{"type": "float"},
				"category": map[string]interface{}{"type": "keyword"},
			},
		},
	}
	client.CreateIndex(ctx, indexName, mapping)

	// 批量索引文档
	global.GVA_LOG.Info("=== 批量索引文档 ===")
	documents := []map[string]interface{}{
		{
			"id":       "1",
			"name":     "iPhone 15",
			"price":    6999.00,
			"category": "手机",
		},
		{
			"id":       "2",
			"name":     "MacBook Pro",
			"price":    12999.00,
			"category": "笔记本",
		},
		{
			"id":       "3",
			"name":     "AirPods Pro",
			"price":    1999.00,
			"category": "耳机",
		},
	}

	if err := client.BulkIndexDocument(ctx, indexName, documents); err != nil {
		global.GVA_LOG.Error("批量索引失败", zap.Error(err))
	}

	// 统计文档数量
	count, err := client.Count(ctx, indexName, nil)
	if err != nil {
		global.GVA_LOG.Error("统计文档数量失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("文档总数", zap.Int64("count", count))
	}

	// 清理
	client.DeleteIndex(ctx, indexName)
}

// ExampleSearch 搜索示例
func (e *ESExample) ExampleSearch() {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		global.GVA_LOG.Error("ES客户端未初始化")
		return
	}

	indexName := "products"

	// 准备测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"name":        map[string]interface{}{"type": "text"},
				"description": map[string]interface{}{"type": "text"},
				"price":       map[string]interface{}{"type": "float"},
				"category":    map[string]interface{}{"type": "keyword"},
			},
		},
	}
	client.CreateIndex(ctx, indexName, mapping)

	documents := []map[string]interface{}{
		{"id": "1", "name": "iPhone 15 Pro", "description": "最新款苹果手机", "price": 7999.00, "category": "手机"},
		{"id": "2", "name": "iPhone 14", "description": "上一代苹果手机", "price": 5999.00, "category": "手机"},
		{"id": "3", "name": "MacBook Pro", "description": "苹果笔记本电脑", "price": 12999.00, "category": "笔记本"},
		{"id": "4", "name": "iPad Pro", "description": "苹果平板电脑", "price": 6499.00, "category": "平板"},
	}
	client.BulkIndexDocument(ctx, indexName, documents)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 1. 匹配所有文档
	global.GVA_LOG.Info("=== 匹配所有文档 ===")
	searchReq := &elasticsearch.SearchRequest{
		Query: elasticsearch.MatchAllQuery(),
		Size:  10,
	}
	result, err := client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果", zap.Int64("总数", result.Hits.Total.Value))
		for _, hit := range result.Hits.Hits {
			global.GVA_LOG.Info("文档", zap.Any("source", hit.Source))
		}
	}

	// 2. 精确匹配
	global.GVA_LOG.Info("=== 精确匹配 category=手机 ===")
	searchReq = &elasticsearch.SearchRequest{
		Query: elasticsearch.TermQuery("category", "手机"),
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 3. 全文搜索
	global.GVA_LOG.Info("=== 全文搜索 name 包含 iPhone ===")
	searchReq = &elasticsearch.SearchRequest{
		Query: elasticsearch.MatchQuery("name", "iPhone"),
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 4. 范围查询
	global.GVA_LOG.Info("=== 范围查询 price >= 6000 and price <= 8000 ===")
	searchReq = &elasticsearch.SearchRequest{
		Query: elasticsearch.RangeQuery("price", 6000, 8000),
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 5. 布尔查询（组合查询）
	global.GVA_LOG.Info("=== 布尔查询：name包含苹果 AND category=手机 ===")
	boolQuery := &elasticsearch.BoolQuery{
		Must: []map[string]interface{}{
			elasticsearch.MatchQuery("name", "苹果"),
			elasticsearch.TermQuery("category", "手机"),
		},
	}
	searchReq = &elasticsearch.SearchRequest{
		Query: boolQuery.ToMap(),
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 6. 分页和排序
	global.GVA_LOG.Info("=== 分页和排序 ===")
	searchReq = &elasticsearch.SearchRequest{
		Query: elasticsearch.MatchAllQuery(),
		From:  0,
		Size:  2,
		Sort: []map[string]interface{}{
			{
				"price": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("搜索结果（按价格降序）", zap.Int64("总数", result.Hits.Total.Value))
		for _, hit := range result.Hits.Hits {
			global.GVA_LOG.Info("文档", zap.Any("source", hit.Source))
		}
	}

	// 清理
	client.DeleteIndex(ctx, indexName)
}

// ExampleAdvancedQueries 高级查询示例
func (e *ESExample) ExampleAdvancedQueries() {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		global.GVA_LOG.Error("ES客户端未初始化")
		return
	}

	indexName := "products"

	// 准备测试数据
	client.CreateIndex(ctx, indexName, nil)
	documents := []map[string]interface{}{
		{"id": "1", "name": "iPhone 15 Pro Max", "category": "手机", "price": 9999.00},
		{"id": "2", "name": "iPhone 15 Pro", "category": "手机", "price": 7999.00},
		{"id": "3", "name": "iPhone 15", "category": "手机", "price": 5999.00},
		{"id": "4", "name": "MacBook Pro 16", "category": "笔记本", "price": 19999.00},
	}
	client.BulkIndexDocument(ctx, indexName, documents)
	time.Sleep(1 * time.Second)

	// 1. 多字段匹配
	global.GVA_LOG.Info("=== 多字段匹配 ===")
	query := elasticsearch.BuildMultiMatchQuery([]string{"name", "category"}, "iPhone")
	searchReq := &elasticsearch.SearchRequest{
		Query: query,
		Size:  10,
	}
	result, err := client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("多字段匹配结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 2. 通配符查询
	global.GVA_LOG.Info("=== 通配符查询 ===")
	query = elasticsearch.BuildWildcardQuery("name", "iPhone*Pro")
	searchReq = &elasticsearch.SearchRequest{
		Query: query,
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("通配符查询结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 3. 前缀查询
	global.GVA_LOG.Info("=== 前缀查询 ===")
	query = elasticsearch.BuildPrefixQuery("name", "iPhone")
	searchReq = &elasticsearch.SearchRequest{
		Query: query,
		Size:  10,
	}
	result, err = client.Search(ctx, indexName, searchReq)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("前缀查询结果", zap.Int64("总数", result.Hits.Total.Value))
	}

	// 4. 聚合查询
	global.GVA_LOG.Info("=== 聚合查询：按分类统计 ===")
	aggs := map[string]interface{}{
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category",
			},
		},
	}
	aggResult, err := client.Aggregate(ctx, indexName, aggs)
	if err != nil {
		global.GVA_LOG.Error("聚合查询失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("聚合结果", zap.Any("aggregations", aggResult))
	}

	// 清理
	client.DeleteIndex(ctx, indexName)
}

// ExampleScrollSearch 游标搜索示例（大数据量）
func (e *ESExample) ExampleScrollSearch() {
	ctx := context.Background()
	client := elasticsearch.NewESClient()
	if client == nil {
		global.GVA_LOG.Error("ES客户端未初始化")
		return
	}

	indexName := "products"

	// 准备大量测试数据
	client.CreateIndex(ctx, indexName, nil)
	documents := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		documents[i] = map[string]interface{}{
			"id":       fmt.Sprintf("%d", i+1),
			"name":     fmt.Sprintf("Product %d", i+1),
			"category": "test",
			"price":    float64(i+1) * 100,
		}
	}
	client.BulkIndexDocument(ctx, indexName, documents)
	time.Sleep(1 * time.Second)

	// 使用游标搜索
	global.GVA_LOG.Info("=== 游标搜索（每次获取10条）===")
	count := 0
	err := client.ScrollSearch(ctx, indexName, elasticsearch.MatchAllQuery(), 10, func(hits []map[string]interface{}) error {
		count += len(hits)
		global.GVA_LOG.Info("处理一批数据", zap.Int("数量", len(hits)))
		return nil
	})
	if err != nil {
		global.GVA_LOG.Error("游标搜索失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("游标搜索完成", zap.Int("总处理数量", count))
	}

	// 清理
	client.DeleteIndex(ctx, indexName)
}

// RunAllExamples 运行所有示例
func RunAllExamples() {
	example := &ESExample{}

	global.GVA_LOG.Info("======= ES 基础操作示例 =======")
	example.ExampleBasicOperations()

	global.GVA_LOG.Info("\n\n======= ES 批量操作示例 =======")
	example.ExampleBulkOperations()

	global.GVA_LOG.Info("\n\n======= ES 搜索示例 =======")
	example.ExampleSearch()

	global.GVA_LOG.Info("\n\n======= ES 高级查询示例 =======")
	example.ExampleAdvancedQueries()

	global.GVA_LOG.Info("\n\n======= ES 游标搜索示例 =======")
	example.ExampleScrollSearch()
}
