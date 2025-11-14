package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

type ElasticsearchService struct{}

var ElasticsearchServiceApp = new(ElasticsearchService)

// 搜索结果
type EsSearchResult struct {
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Data     []interface{} `json:"data"`
}

// 索引单个文档
func (e *ElasticsearchService) IndexDocument(ctx context.Context,
	index string,
	docID string,
	doc interface{}) error {
	client := global.GVA_ES
	if client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	// 序列化文档
	data, err := json.Marshal(doc)
	if err != nil {
		global.GVA_LOG.Error("文档序列化失败", zap.Error(err))
		return err
	}

	// 索引文档
	res, err := client.Index(
		index,
		bytes.NewReader(data),
		client.Index.WithDocumentID(docID),
		client.Index.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("索引文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}

	return nil
}

// 批量索引文档
func (e *ElasticsearchService) BulkIndex(ctx context.Context,
	index string,
	docs map[string]interface{}) error {
	client := global.GVA_ES
	if client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	var buf bytes.Buffer

	for docID, doc := range docs {
		// 元数据行
		meta := []byte(fmt.Sprintf(`{"index":{"_id":"%s"}}%s`, docID, "\n"))
		buf.Write(meta)

		// 文档数据行
		data, err := json.Marshal(doc)
		if err != nil {
			global.GVA_LOG.Error("文档序列化失败", zap.Error(err))
			continue
		}
		buf.Write(data)
		buf.WriteString("\n")
	}

	// 批量索引
	res, err := client.Bulk(
		bytes.NewReader(buf.Bytes()),
		client.Bulk.WithIndex(index),
		client.Bulk.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("批量索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch bulk error: %s", res.String())
	}

	return nil
}

// 搜索文档
func (e *ElasticsearchService) Search(ctx context.Context,
	index string,
	query map[string]interface{},
	page int,
	pageSize int) (*EsSearchResult, error) {
	client := global.GVA_ES
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client not initialized")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建查询体
	queryBody := map[string]interface{}{
		"query": query,
		"from":  (page - 1) * pageSize,
		"size":  pageSize,
	}

	data, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	// 执行搜索
	res, err := client.Search(
		client.Search.WithIndex(index),
		client.Search.WithBody(bytes.NewReader(data)),
		client.Search.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	// 解析响应
	var searchResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		global.GVA_LOG.Error("解析搜索响应失败", zap.Error(err))
		return nil, err
	}

	// 提取总数和数据
	hits := searchResp["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))

	var documents []interface{}
	for _, hit := range hits["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		documents = append(documents, hitMap["_source"])
	}

	return &EsSearchResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     documents,
	}, nil
}

// 简单文本搜索
func (e *ElasticsearchService) SimpleSearch(ctx context.Context,
	index string,
	keyword string,
	fields []string,
	page int,
	pageSize int) (*EsSearchResult, error) {
	if len(fields) == 0 {
		fields = []string{"title", "content"}
	}

	query := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":  keyword,
			"fields": fields,
		},
	}

	return e.Search(ctx, index, query, page, pageSize)
}

// 删除文档
func (e *ElasticsearchService) DeleteDocument(ctx context.Context,
	index string,
	docID string) error {
	client := global.GVA_ES
	if client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	res, err := client.Delete(
		index,
		docID,
		client.Delete.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("删除文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch delete error: %s", res.String())
	}

	return nil
}

// 更新文档
func (e *ElasticsearchService) UpdateDocument(ctx context.Context,
	index string,
	docID string,
	updates map[string]interface{}) error {
	client := global.GVA_ES
	if client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	updateBody := map[string]interface{}{
		"doc": updates,
	}

	data, err := json.Marshal(updateBody)
	if err != nil {
		return err
	}

	res, err := client.Update(
		index,
		docID,
		bytes.NewReader(data),
		client.Update.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("更新文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch update error: %s", res.String())
	}

	return nil
}

// 获取单个文档
func (e *ElasticsearchService) GetDocument(ctx context.Context,
	index string,
	docID string) (map[string]interface{}, error) {
	client := global.GVA_ES
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client not initialized")
	}

	res, err := client.Get(
		index,
		docID,
		client.Get.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("获取文档失败", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("elasticsearch get error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result["_source"].(map[string]interface{}), nil
}

// 删除索引
func (e *ElasticsearchService) DeleteIndex(ctx context.Context, index string) error {
	client := global.GVA_ES
	if client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	res, err := client.Indices.Delete(
		[]string{index},
		client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		global.GVA_LOG.Error("删除索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch delete index error: %s", res.String())
	}

	return nil
}

// 获取索引信息
func (e *ElasticsearchService) GetIndexInfo(ctx context.Context, index string) (map[string]interface{}, error) {
	client := global.GVA_ES
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client not initialized")
	}

	res, err := client.Indices.Get(
		[]string{index},
		client.Indices.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
