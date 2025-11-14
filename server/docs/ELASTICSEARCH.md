# Elasticsearch 集成使用指南

gin-vue-admin 已经集成了 Elasticsearch 支持，本指南将帮助你快速上手使用。

## 1. 配置 Elasticsearch

### 1.1 修改配置文件

编辑 `server/config.yaml`，配置 Elasticsearch 连接信息：

```yaml
# system configuration
system:
    use-elasticsearch: true  # 启用 Elasticsearch

# elasticsearch configuration
elasticsearch:
    addresses:
        - http://127.0.0.1:9200  # ES集群地址，可以配置多个
    username: ""                 # 用户名（如果需要认证）
    password: ""                 # 密码（如果需要认证）
    index: gva_default          # 默认索引名称
    max-retries: 3              # 最大重试次数
    timeout: 30                 # 请求超时时间（秒）
    compression: ""             # 压缩方式，可选 "gzip"
    cert-file: ""               # HTTPS 证书文件路径（如果使用 HTTPS）
```

### 1.2 启动 Elasticsearch

确保 Elasticsearch 服务正在运行：

```bash
# Docker 方式启动（推荐用于开发环境）
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  elasticsearch:8.19.0
```

## 2. 使用 ES 工具类

### 2.1 初始化客户端

```go
import (
    "context"
    "github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"
)

// 创建 ES 客户端实例
client := elasticsearch.NewESClient()
ctx := context.Background()
```

### 2.2 索引管理

```go
// 创建索引
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
            "price": map[string]interface{}{
                "type": "float",
            },
        },
    },
}
err := client.CreateIndex(ctx, "products", mapping)

// 检查索引是否存在
exists, err := client.IndexExists(ctx, "products")

// 删除索引
err = client.DeleteIndex(ctx, "products")
```

### 2.3 文档操作

#### 创建/更新文档

```go
type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

product := Product{
    ID:    "1",
    Name:  "iPhone 15",
    Price: 6999.00,
}

// 索引文档
err := client.IndexDocument(ctx, "products", product.ID, product)
```

#### 获取文档

```go
var product Product
err := client.GetDocument(ctx, "products", "1", &product)
```

#### 更新文档

```go
update := map[string]interface{}{
    "price": 6499.00,
}
err := client.UpdateDocument(ctx, "products", "1", update)
```

#### 删除文档

```go
err := client.DeleteDocument(ctx, "products", "1")
```

### 2.4 批量操作

```go
documents := []map[string]interface{}{
    {
        "id":    "1",
        "name":  "iPhone 15",
        "price": 6999.00,
    },
    {
        "id":    "2",
        "name":  "MacBook Pro",
        "price": 12999.00,
    },
}

err := client.BulkIndexDocument(ctx, "products", documents)
```

### 2.5 搜索操作

#### 匹配所有文档

```go
searchReq := &elasticsearch.SearchRequest{
    Query: elasticsearch.MatchAllQuery(),
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 精确匹配（Term Query）

```go
searchReq := &elasticsearch.SearchRequest{
    Query: elasticsearch.TermQuery("category", "手机"),
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 全文搜索（Match Query）

```go
searchReq := &elasticsearch.SearchRequest{
    Query: elasticsearch.MatchQuery("name", "iPhone"),
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 范围查询

```go
searchReq := &elasticsearch.SearchRequest{
    Query: elasticsearch.RangeQuery("price", 5000, 8000),
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 布尔查询（组合查询）

```go
boolQuery := &elasticsearch.BoolQuery{
    Must: []map[string]interface{}{
        elasticsearch.MatchQuery("name", "iPhone"),
        elasticsearch.TermQuery("category", "手机"),
    },
    Filter: []map[string]interface{}{
        elasticsearch.RangeQuery("price", 5000, 10000),
    },
}

searchReq := &elasticsearch.SearchRequest{
    Query: boolQuery.ToMap(),
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 分页和排序

```go
searchReq := &elasticsearch.SearchRequest{
    Query: elasticsearch.MatchAllQuery(),
    From:  0,  // 起始位置
    Size:  20, // 每页数量
    Sort: []map[string]interface{}{
        {
            "price": map[string]interface{}{
                "order": "desc", // 按价格降序
            },
        },
    },
}
result, err := client.Search(ctx, "products", searchReq)
```

### 2.6 高级查询

#### 多字段匹配

```go
query := elasticsearch.BuildMultiMatchQuery([]string{"name", "description"}, "iPhone")
searchReq := &elasticsearch.SearchRequest{
    Query: query,
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 通配符查询

```go
query := elasticsearch.BuildWildcardQuery("name", "iPhone*Pro")
searchReq := &elasticsearch.SearchRequest{
    Query: query,
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 前缀查询

```go
query := elasticsearch.BuildPrefixQuery("name", "iPhone")
searchReq := &elasticsearch.SearchRequest{
    Query: query,
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

#### 模糊查询

```go
query := elasticsearch.BuildFuzzyQuery("name", "iPhne", 2) // 允许2个字符差异
searchReq := &elasticsearch.SearchRequest{
    Query: query,
    Size:  10,
}
result, err := client.Search(ctx, "products", searchReq)
```

### 2.7 聚合查询

```go
// 按分类统计
aggs := map[string]interface{}{
    "categories": map[string]interface{}{
        "terms": map[string]interface{}{
            "field": "category",
        },
    },
}
aggResult, err := client.Aggregate(ctx, "products", aggs)
```

### 2.8 游标搜索（大数据量）

```go
// 处理大量数据时使用游标搜索
err := client.ScrollSearch(ctx, "products", elasticsearch.MatchAllQuery(), 100,
    func(hits []map[string]interface{}) error {
        // 处理每一批数据
        for _, hit := range hits {
            // 处理逻辑
        }
        return nil
    })
```

### 2.9 统计文档数量

```go
// 统计所有文档
count, err := client.Count(ctx, "products", nil)

// 带条件统计
count, err := client.Count(ctx, "products", elasticsearch.TermQuery("category", "手机"))
```

## 3. 在 Service 中使用

```go
package service

import (
    "context"
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"
    "go.uber.org/zap"
)

type ProductService struct{}

// SearchProducts 搜索产品
func (s *ProductService) SearchProducts(keyword string, page, pageSize int) ([]map[string]interface{}, int64, error) {
    ctx := context.Background()
    client := elasticsearch.NewESClient()
    if client == nil {
        return nil, 0, fmt.Errorf("ES客户端未初始化")
    }

    // 构建搜索请求
    searchReq := &elasticsearch.SearchRequest{
        Query: elasticsearch.MatchQuery("name", keyword),
        From:  (page - 1) * pageSize,
        Size:  pageSize,
        Sort: []map[string]interface{}{
            {
                "price": map[string]interface{}{
                    "order": "desc",
                },
            },
        },
    }

    // 执行搜索
    result, err := client.Search(ctx, "products", searchReq)
    if err != nil {
        global.GVA_LOG.Error("搜索失败", zap.Error(err))
        return nil, 0, err
    }

    // 提取结果
    products := make([]map[string]interface{}, 0, len(result.Hits.Hits))
    for _, hit := range result.Hits.Hits {
        products = append(products, hit.Source)
    }

    return products, result.Hits.Total.Value, nil
}
```

## 4. 运行示例代码

项目中提供了完整的示例代码，可以通过以下方式运行：

```go
import "github.com/flipped-aurora/gin-vue-admin/server/examples"

// 运行所有示例
examples.RunAllExamples()
```

示例代码位置：`server/examples/elasticsearch_example.go`

## 5. 常用查询构建器

| 方法 | 说明 | 用途 |
|------|------|------|
| `MatchAllQuery()` | 匹配所有文档 | 获取所有数据 |
| `MatchQuery(field, value)` | 全文搜索 | 模糊匹配文本字段 |
| `TermQuery(field, value)` | 精确匹配 | 精确匹配 keyword 字段 |
| `RangeQuery(field, gte, lte)` | 范围查询 | 数值、日期范围查询 |
| `BoolQuery` | 布尔查询 | 组合多个查询条件 |
| `BuildMultiMatchQuery(fields, text)` | 多字段匹配 | 在多个字段中搜索 |
| `BuildWildcardQuery(field, value)` | 通配符查询 | 支持 * 和 ? 通配符 |
| `BuildPrefixQuery(field, value)` | 前缀查询 | 匹配指定前缀 |
| `BuildFuzzyQuery(field, value, fuzziness)` | 模糊查询 | 容错搜索 |
| `BuildExistsQuery(field)` | 字段存在查询 | 查找包含指定字段的文档 |
| `BuildIdsQuery(ids)` | ID 查询 | 通过文档 ID 批量查询 |

## 6. 注意事项

1. **Go 版本要求**：项目需要 Go 1.23+ 版本，如果你的 Go 版本较低，需要升级：
   ```bash
   # 下载并安装最新版本的 Go
   # https://go.dev/dl/
   ```

2. **ES 连接失败**：确保 Elasticsearch 服务正常运行，并且配置文件中的地址正确。

3. **索引映射**：创建索引时建议明确指定字段类型，避免自动推断导致的类型不匹配。

4. **分页限制**：Elasticsearch 默认最多返回 10000 条数据，如需处理更多数据，使用游标搜索 `ScrollSearch`。

5. **性能优化**：
   - 合理设置分片数量
   - 使用批量操作提高性能
   - 避免使用 wildcard 和 fuzzy 查询在大数据集上
   - 合理设置超时时间

## 7. 常见问题

### Q: 如何启用 Elasticsearch？
A: 在 `config.yaml` 中设置 `system.use-elasticsearch: true`

### Q: 搜索结果为空？
A: 检查索引是否刷新，可以等待几秒或手动调用 `client.Refresh(ctx, indexName)`

### Q: 中文搜索不准确？
A: 需要配置中文分词器（如 IK 分词器）：
```bash
# 安装 IK 分词器
docker exec -it elasticsearch bin/elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.19.0/elasticsearch-analysis-ik-8.19.0.zip
```

### Q: 如何调试 ES 查询？
A: 可以查看日志输出，或者使用 Kibana Dev Tools 测试查询语句。

## 8. 相关资源

- [Elasticsearch 官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [Go Elasticsearch 客户端文档](https://github.com/elastic/go-elasticsearch)
- [gin-vue-admin 项目文档](https://www.gin-vue-admin.com/)

## 9. 技术支持

如有问题，请提交 Issue 或联系项目维护者。
