package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// ESClient Elasticsearch 客户端封装
type ESClient struct {
	client *elasticsearch.Client
}

// NewESClient 创建ES客户端实例
func NewESClient() *ESClient {
	if global.GVA_ES == nil {
		global.GVA_LOG.Error("Elasticsearch客户端未初始化")
		return nil
	}
	return &ESClient{client: global.GVA_ES}
}

// IndexDocument 创建或更新文档
// index: 索引名称
// docID: 文档ID，如果为空则自动生成
// document: 要索引的文档数据
func (e *ESClient) IndexDocument(ctx context.Context, index, docID string, document interface{}) error {
	data, err := json.Marshal(document)
	if err != nil {
		global.GVA_LOG.Error("序列化文档失败", zap.Error(err))
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("索引文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("索引文档错误: %s", res.String())
	}

	global.GVA_LOG.Info("文档索引成功", zap.String("index", index), zap.String("docID", docID))
	return nil
}

// GetDocument 获取文档
// index: 索引名称
// docID: 文档ID
// result: 用于接收结果的指针
func (e *ESClient) GetDocument(ctx context.Context, index, docID string, result interface{}) error {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("获取文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return fmt.Errorf("文档不存在")
		}
		return fmt.Errorf("获取文档错误: %s", res.String())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(body, &doc); err != nil {
		return err
	}

	// 提取 _source 字段
	if source, ok := doc["_source"]; ok {
		sourceBytes, _ := json.Marshal(source)
		return json.Unmarshal(sourceBytes, result)
	}

	return fmt.Errorf("文档格式错误")
}

// DeleteDocument 删除文档
// index: 索引名称
// docID: 文档ID
func (e *ESClient) DeleteDocument(ctx context.Context, index, docID string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("删除文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return fmt.Errorf("文档不存在")
		}
		return fmt.Errorf("删除文档错误: %s", res.String())
	}

	global.GVA_LOG.Info("文档删除成功", zap.String("index", index), zap.String("docID", docID))
	return nil
}

// UpdateDocument 更新文档
// index: 索引名称
// docID: 文档ID
// doc: 要更新的字段
func (e *ESClient) UpdateDocument(ctx context.Context, index, docID string, doc interface{}) error {
	update := map[string]interface{}{
		"doc": doc,
	}

	data, err := json.Marshal(update)
	if err != nil {
		global.GVA_LOG.Error("序列化更新文档失败", zap.Error(err))
		return err
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("更新文档失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("更新文档错误: %s", res.String())
	}

	global.GVA_LOG.Info("文档更新成功", zap.String("index", index), zap.String("docID", docID))
	return nil
}

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Query map[string]interface{}   `json:"query"`
	From  int                      `json:"from,omitempty"`
	Size  int                      `json:"size,omitempty"`
	Sort  []map[string]interface{} `json:"sort,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Took int64 `json:"took"`
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string                 `json:"_index"`
			ID     string                 `json:"_id"`
			Score  float64                `json:"_score"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Search 搜索文档
// index: 索引名称，多个索引用逗号分隔
// searchReq: 搜索请求参数
func (e *ESClient) Search(ctx context.Context, index string, searchReq *SearchRequest) (*SearchResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchReq); err != nil {
		global.GVA_LOG.Error("编码搜索请求失败", zap.Error(err))
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("解析错误响应失败: %s", err)
		}
		return nil, fmt.Errorf("搜索错误: %v", e)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// MatchAllQuery 创建匹配所有文档的查询
func MatchAllQuery() map[string]interface{} {
	return map[string]interface{}{
		"match_all": map[string]interface{}{},
	}
}

// MatchQuery 创建匹配查询
func MatchQuery(field, value string) map[string]interface{} {
	return map[string]interface{}{
		"match": map[string]interface{}{
			field: value,
		},
	}
}

// TermQuery 创建精确匹配查询
func TermQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	}
}

// RangeQuery 创建范围查询
func RangeQuery(field string, gte, lte interface{}) map[string]interface{} {
	rangeMap := make(map[string]interface{})
	if gte != nil {
		rangeMap["gte"] = gte
	}
	if lte != nil {
		rangeMap["lte"] = lte
	}
	return map[string]interface{}{
		"range": map[string]interface{}{
			field: rangeMap,
		},
	}
}

// BoolQuery 创建布尔查询
type BoolQuery struct {
	Must               []map[string]interface{} `json:"must,omitempty"`
	Should             []map[string]interface{} `json:"should,omitempty"`
	MustNot            []map[string]interface{} `json:"must_not,omitempty"`
	Filter             []map[string]interface{} `json:"filter,omitempty"`
	MinimumShouldMatch int                      `json:"minimum_should_match,omitempty"`
}

func (b *BoolQuery) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"bool": b,
	}
}

// CreateIndex 创建索引
// index: 索引名称
// mapping: 索引映射配置（可选）
func (e *ESClient) CreateIndex(ctx context.Context, index string, mapping map[string]interface{}) error {
	var body io.Reader
	if mapping != nil {
		data, err := json.Marshal(mapping)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  body,
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("创建索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("创建索引错误: %s", res.String())
	}

	global.GVA_LOG.Info("索引创建成功", zap.String("index", index))
	return nil
}

// DeleteIndex 删除索引
func (e *ESClient) DeleteIndex(ctx context.Context, index string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		global.GVA_LOG.Error("删除索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("删除索引错误: %s", res.String())
	}

	global.GVA_LOG.Info("索引删除成功", zap.String("index", index))
	return nil
}

// IndexExists 检查索引是否存在
func (e *ESClient) IndexExists(ctx context.Context, index string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// BulkIndexDocument 批量索引文档
// index: 索引名称
// documents: 文档列表，每个元素是 map，必须包含 "id" 字段
func (e *ESClient) BulkIndexDocument(ctx context.Context, index string, documents []map[string]interface{}) error {
	if len(documents) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, doc := range documents {
		docID, ok := doc["id"]
		if !ok {
			return fmt.Errorf("文档缺少 id 字段")
		}

		// 删除 id 字段，避免重复
		delete(doc, "id")

		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
				"_id":    docID,
			},
		}

		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')

		docBytes, _ := json.Marshal(doc)
		buf.Write(docBytes)
		buf.WriteByte('\n')
	}

	res, err := e.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		e.client.Bulk.WithContext(ctx),
		e.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		global.GVA_LOG.Error("批量索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("批量索引错误: %s", res.String())
	}

	// 检查是否有错误
	var bulkResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		return err
	}

	if bulkResp["errors"].(bool) {
		return fmt.Errorf("批量索引存在错误")
	}

	global.GVA_LOG.Info("批量索引成功", zap.String("index", index), zap.Int("count", len(documents)))
	return nil
}

// Refresh 刷新索引
func (e *ESClient) Refresh(ctx context.Context, indices ...string) error {
	req := esapi.IndicesRefreshRequest{
		Index: indices,
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("刷新索引错误: %s", res.String())
	}

	return nil
}

// Count 获取文档数量
func (e *ESClient) Count(ctx context.Context, index string, query map[string]interface{}) (int64, error) {
	var buf bytes.Buffer
	if query != nil {
		queryMap := map[string]interface{}{
			"query": query,
		}
		if err := json.NewEncoder(&buf).Encode(queryMap); err != nil {
			return 0, err
		}
	}

	var body io.Reader
	if buf.Len() > 0 {
		body = &buf
	}

	res, err := e.client.Count(
		e.client.Count.WithContext(ctx),
		e.client.Count.WithIndex(index),
		e.client.Count.WithBody(body),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("统计文档错误: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	count, ok := result["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("无法解析文档数量")
	}

	return int64(count), nil
}

// MultiSearch 多索引搜索
func (e *ESClient) MultiSearch(ctx context.Context, searches []struct {
	Index string
	Query map[string]interface{}
}) ([]SearchResponse, error) {
	var buf bytes.Buffer

	for _, search := range searches {
		// 写入 header
		header := map[string]interface{}{
			"index": search.Index,
		}
		headerBytes, _ := json.Marshal(header)
		buf.Write(headerBytes)
		buf.WriteByte('\n')

		// 写入 query
		queryMap := map[string]interface{}{
			"query": search.Query,
		}
		queryBytes, _ := json.Marshal(queryMap)
		buf.Write(queryBytes)
		buf.WriteByte('\n')
	}

	res, err := e.client.Msearch(
		bytes.NewReader(buf.Bytes()),
		e.client.Msearch.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("多索引搜索错误: %s", res.String())
	}

	var msearchResp struct {
		Responses []SearchResponse `json:"responses"`
	}
	if err := json.NewDecoder(res.Body).Decode(&msearchResp); err != nil {
		return nil, err
	}

	return msearchResp.Responses, nil
}

// Aggregate 聚合查询
func (e *ESClient) Aggregate(ctx context.Context, index string, aggs map[string]interface{}) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": aggs,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("聚合查询错误: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	aggregations, ok := result["aggregations"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无法解析聚合结果")
	}

	return aggregations, nil
}

// ScrollSearch 游标搜索（用于大数据量查询）
func (e *ESClient) ScrollSearch(ctx context.Context, index string, query map[string]interface{}, size int, processFunc func(hits []map[string]interface{}) error) error {
	// 初始搜索
	searchReq := map[string]interface{}{
		"query": query,
		"size":  size,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchReq); err != nil {
		return err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithScroll(time.Minute*2),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("游标搜索错误: %s", res.String())
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return err
	}

	// 获取 scroll_id
	scrollID := res.Header.Get("X-Elastic-Scroll-Id")
	if scrollID == "" {
		// 尝试从响应体中获取
		body, _ := io.ReadAll(res.Body)
		var scrollResp map[string]interface{}
		if err := json.Unmarshal(body, &scrollResp); err == nil {
			if sid, ok := scrollResp["_scroll_id"].(string); ok {
				scrollID = sid
			}
		}
	}

	// 处理第一批结果
	hits := make([]map[string]interface{}, len(searchResp.Hits.Hits))
	for i, hit := range searchResp.Hits.Hits {
		hits[i] = hit.Source
	}
	if err := processFunc(hits); err != nil {
		return err
	}

	// 继续滚动获取数据
	for {
		scrollRes, err := e.client.Scroll(
			e.client.Scroll.WithContext(ctx),
			e.client.Scroll.WithScrollID(scrollID),
			e.client.Scroll.WithScroll(time.Minute*2),
		)
		if err != nil {
			return err
		}

		if scrollRes.IsError() {
			scrollRes.Body.Close()
			break
		}

		var scrollResp SearchResponse
		if err := json.NewDecoder(scrollRes.Body).Decode(&scrollResp); err != nil {
			scrollRes.Body.Close()
			return err
		}
		scrollRes.Body.Close()

		if len(scrollResp.Hits.Hits) == 0 {
			break
		}

		hits := make([]map[string]interface{}, len(scrollResp.Hits.Hits))
		for i, hit := range scrollResp.Hits.Hits {
			hits[i] = hit.Source
		}
		if err := processFunc(hits); err != nil {
			return err
		}
	}

	// 清除 scroll
	e.client.ClearScroll(
		e.client.ClearScroll.WithScrollID(scrollID),
	)

	return nil
}

// GetIndexMapping 获取索引映射
func (e *ESClient) GetIndexMapping(ctx context.Context, index string) (map[string]interface{}, error) {
	req := esapi.IndicesGetMappingRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("获取索引映射错误: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateIndexMapping 更新索引映射
func (e *ESClient) UpdateIndexMapping(ctx context.Context, index string, mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	req := esapi.IndicesPutMappingRequest{
		Index: []string{index},
		Body:  bytes.NewReader(data),
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("更新索引映射错误: %s", res.String())
	}

	global.GVA_LOG.Info("索引映射更新成功", zap.String("index", index))
	return nil
}

// HighlightSearch 高亮搜索
func (e *ESClient) HighlightSearch(ctx context.Context, index string, query map[string]interface{}, fields []string) (*SearchResponse, error) {
	highlightFields := make(map[string]interface{})
	for _, field := range fields {
		highlightFields[field] = map[string]interface{}{}
	}

	searchReq := map[string]interface{}{
		"query": query,
		"highlight": map[string]interface{}{
			"fields": highlightFields,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchReq); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("高亮搜索错误: %s", res.String())
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// SuggestSearch 建议搜索（自动补全）
func (e *ESClient) SuggestSearch(ctx context.Context, index, field, text string) ([]string, error) {
	suggest := map[string]interface{}{
		"my-suggest": map[string]interface{}{
			"text": text,
			"term": map[string]interface{}{
				"field": field,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]interface{}{
		"suggest": suggest,
	}); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("建议搜索错误: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	suggestions := []string{}
	if suggestResult, ok := result["suggest"].(map[string]interface{}); ok {
		if mySuggest, ok := suggestResult["my-suggest"].([]interface{}); ok {
			for _, item := range mySuggest {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if options, ok := itemMap["options"].([]interface{}); ok {
						for _, option := range options {
							if optionMap, ok := option.(map[string]interface{}); ok {
								if textVal, ok := optionMap["text"].(string); ok {
									suggestions = append(suggestions, textVal)
								}
							}
						}
					}
				}
			}
		}
	}

	return suggestions, nil
}

// ReIndex 重新索引（从一个索引复制到另一个索引）
func (e *ESClient) ReIndex(ctx context.Context, sourceIndex, destIndex string) error {
	body := map[string]interface{}{
		"source": map[string]interface{}{
			"index": sourceIndex,
		},
		"dest": map[string]interface{}{
			"index": destIndex,
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.ReindexRequest{
		Body: bytes.NewReader(data),
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("重新索引错误: %s", res.String())
	}

	global.GVA_LOG.Info("重新索引成功", zap.String("source", sourceIndex), zap.String("dest", destIndex))
	return nil
}

// BuildMultiMatchQuery 构建多字段匹配查询
func BuildMultiMatchQuery(fields []string, text string) map[string]interface{} {
	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":  text,
			"fields": fields,
		},
	}
}

// BuildWildcardQuery 构建通配符查询
func BuildWildcardQuery(field, value string) map[string]interface{} {
	return map[string]interface{}{
		"wildcard": map[string]interface{}{
			field: value,
		},
	}
}

// BuildPrefixQuery 构建前缀查询
func BuildPrefixQuery(field, value string) map[string]interface{} {
	return map[string]interface{}{
		"prefix": map[string]interface{}{
			field: value,
		},
	}
}

// BuildFuzzyQuery 构建模糊查询
func BuildFuzzyQuery(field, value string, fuzziness int) map[string]interface{} {
	fuzzyMap := map[string]interface{}{
		"value": value,
	}
	if fuzziness > 0 {
		fuzzyMap["fuzziness"] = fuzziness
	}
	return map[string]interface{}{
		"fuzzy": map[string]interface{}{
			field: fuzzyMap,
		},
	}
}

// BuildExistsQuery 构建字段存在查询
func BuildExistsQuery(field string) map[string]interface{} {
	return map[string]interface{}{
		"exists": map[string]interface{}{
			"field": field,
		},
	}
}

// BuildIdsQuery 构建ID查询
func BuildIdsQuery(ids []string) map[string]interface{} {
	return map[string]interface{}{
		"ids": map[string]interface{}{
			"values": ids,
		},
	}
}

// ParseSortString 解析排序字符串（例如："field1:asc,field2:desc"）
func ParseSortString(sortStr string) []map[string]interface{} {
	if sortStr == "" {
		return nil
	}

	var sorts []map[string]interface{}
	sortFields := strings.Split(sortStr, ",")
	for _, sortField := range sortFields {
		parts := strings.Split(strings.TrimSpace(sortField), ":")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			order := strings.TrimSpace(parts[1])
			sorts = append(sorts, map[string]interface{}{
				field: map[string]interface{}{
					"order": order,
				},
			})
		}
	}

	return sorts
}
