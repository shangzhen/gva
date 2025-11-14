# 用户操作日志 ES 集成示例

这是一个完整的 Elasticsearch 集成示例，实现了用户操作日志的记录、查询和统计功能。

## 功能特性

- ✅ 完整的 ES Mapping 定义
- ✅ 日志的增删改查操作
- ✅ 多条件组合搜索
- ✅ 时间范围查询
- ✅ 分页和排序
- ✅ 聚合统计
- ✅ RESTful API 接口

## 目录结构

```
server/
├── model/system/
│   ├── sys_user_action_log.go           # 数据模型和 Mapping 定义
│   ├── request/
│   │   └── sys_user_action_log.go       # 请求参数定义
│   └── response/
│       └── sys_user_action_log.go       # 响应数据定义
├── service/system/
│   └── sys_user_action_log.go           # 业务逻辑层
├── api/v1/system/
│   └── sys_user_action_log.go           # API 处理器
├── router/system/
│   └── sys_user_action_log.go           # 路由定义
└── test/
    └── elasticsearch_test.go            # 测试脚本
```

## 数据模型

### ES Mapping 定义

用户操作日志的 Mapping 包含以下字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | keyword | 日志唯一ID |
| user_id | long | 用户ID |
| username | keyword | 用户名 |
| action | keyword | 操作动作（login, logout, create等） |
| module | keyword | 操作模块（user, role, menu等） |
| method | keyword | HTTP方法（GET, POST, PUT, DELETE） |
| path | text | 请求路径 |
| ip | ip | IP地址 |
| user_agent | text | 用户代理 |
| status | integer | 响应状态码 |
| latency | long | 响应时间（毫秒） |
| request | text | 请求参数（不索引） |
| response | text | 响应数据（不索引） |
| error_msg | text | 错误信息 |
| create_time | date | 创建时间 |

### 索引设置

- **分片数**: 3
- **副本数**: 1
- **刷新间隔**: 5秒
- **分析器**: standard

## API 接口

### 1. 初始化索引

创建 ES 索引和 Mapping。

**请求**

```http
POST /userActionLog/initIndex
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "msg": "初始化索引成功"
}
```

### 2. 创建日志

创建一条用户操作日志。

**请求**

```http
POST /userActionLog/createLog
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1,
  "username": "admin",
  "action": "login",
  "module": "user",
  "method": "POST",
  "path": "/base/login",
  "ip": "127.0.0.1",
  "user_agent": "Mozilla/5.0",
  "status": 200,
  "latency": 123
}
```

**响应**

```json
{
  "code": 0,
  "msg": "创建成功"
}
```

### 3. 获取单条日志

根据 ID 获取日志详情。

**请求**

```http
GET /userActionLog/getLog/:id
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "data": {
    "id": "xxx-xxx-xxx",
    "user_id": 1,
    "username": "admin",
    "action": "login",
    "module": "user",
    "method": "POST",
    "path": "/base/login",
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0",
    "status": 200,
    "latency": 123,
    "create_time": "2024-01-01T00:00:00Z"
  },
  "msg": "查询成功"
}
```

### 4. 搜索日志

多条件搜索用户操作日志。

**请求**

```http
POST /userActionLog/searchLogs
Authorization: Bearer <token>
Content-Type: application/json

{
  "page": 1,
  "pageSize": 10,
  "user_id": 1,
  "action": "login",
  "module": "user",
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-12-31T23:59:59Z",
  "order_field": "create_time",
  "order_type": "desc"
}
```

**支持的搜索条件**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码（必填） |
| pageSize | int | 每页数量（必填） |
| user_id | int | 用户ID（精确匹配） |
| username | string | 用户名（模糊搜索） |
| action | string | 操作动作（精确匹配） |
| module | string | 操作模块（精确匹配） |
| method | string | 请求方法（精确匹配） |
| ip | string | IP地址（精确匹配） |
| status | int | 响应状态码（精确匹配） |
| start_time | string | 开始时间（ISO8601格式） |
| end_time | string | 结束时间（ISO8601格式） |
| keyword | string | 关键词（搜索path和error_msg） |
| order_field | string | 排序字段 |
| order_type | string | 排序类型（asc/desc） |

**响应**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": "xxx-xxx-xxx",
        "user_id": 1,
        "username": "admin",
        "action": "login",
        "module": "user",
        "create_time": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10
  },
  "msg": "查询成功"
}
```

### 5. 删除日志

根据 ID 删除日志。

**请求**

```http
DELETE /userActionLog/deleteLog/:id
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "msg": "删除成功"
}
```

### 6. 获取统计数据

获取日志的聚合统计数据。

**请求**

```http
POST /userActionLog/getStats
Authorization: Bearer <token>
Content-Type: application/json

{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-12-31T23:59:59Z",
  "group_by": "action"
}
```

**支持的分组字段**

- `action` - 按操作动作分组
- `module` - 按操作模块分组
- `user_id` - 按用户ID分组
- `username` - 按用户名分组
- `method` - 按请求方法分组
- `status` - 按响应状态码分组

**响应**

```json
{
  "code": 0,
  "data": {
    "total": 1000,
    "stats": [
      {
        "key": "login",
        "doc_count": 500
      },
      {
        "key": "logout",
        "doc_count": 300
      },
      {
        "key": "create",
        "doc_count": 200
      }
    ]
  },
  "msg": "统计成功"
}
```

### 7. 删除索引（危险操作）

删除整个 ES 索引（谨慎使用）。

**请求**

```http
DELETE /userActionLog/deleteIndex
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "msg": "删除索引成功"
}
```

## 使用示例

### 1. 配置 Elasticsearch

编辑 `server/config.yaml`：

```yaml
system:
    use-elasticsearch: true  # 启用 ES

elasticsearch:
    addresses:
        - http://127.0.0.1:9200
    index: user_action_logs
    max-retries: 3
    timeout: 30
```

### 2. 启动 Elasticsearch

```bash
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  elasticsearch:8.19.0
```

### 3. 启动项目

```bash
cd server
go run main.go
```

### 4. 初始化索引

使用 Postman 或 curl：

```bash
curl -X POST http://localhost:8888/userActionLog/initIndex \
  -H "x-token: your-jwt-token"
```

### 5. 创建测试数据

```bash
curl -X POST http://localhost:8888/userActionLog/createLog \
  -H "x-token: your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "username": "admin",
    "action": "login",
    "module": "user",
    "method": "POST",
    "path": "/base/login",
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0",
    "status": 200,
    "latency": 123
  }'
```

### 6. 搜索日志

```bash
curl -X POST http://localhost:8888/userActionLog/searchLogs \
  -H "x-token: your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "pageSize": 10,
    "action": "login"
  }'
```

### 7. 运行测试脚本

```bash
cd server/test
go run elasticsearch_test.go
```

**注意**：需要先修改测试脚本中的 `token` 变量为实际的 JWT token。

## 代码示例

### 在代码中记录日志

```go
import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
)

// 记录用户操作日志
func recordUserAction(c *gin.Context) {
    userActionLogService := service.ServiceGroupApp.SystemServiceGroup.UserActionLogService

    logReq := &request.UserActionLogCreate{
        UserID:    1,
        Username:  "admin",
        Action:    "create",
        Module:    "user",
        Method:    c.Request.Method,
        Path:      c.Request.URL.Path,
        IP:        c.ClientIP(),
        UserAgent: c.Request.UserAgent(),
        Status:    200,
        Latency:   123,
    }

    if err := userActionLogService.CreateLog(logReq); err != nil {
        // 记录失败处理
    }
}
```

### 在中间件中自动记录

```go
func UserActionLogMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()

        // 处理请求
        c.Next()

        // 计算响应时间
        latency := time.Since(startTime).Milliseconds()

        // 获取用户信息
        userID, _ := c.Get("userId")
        username, _ := c.Get("username")

        // 记录日志
        userActionLogService.CreateLog(&request.UserActionLogCreate{
            UserID:    userID.(uint),
            Username:  username.(string),
            Action:    getAction(c.Request.Method, c.Request.URL.Path),
            Module:    getModule(c.Request.URL.Path),
            Method:    c.Request.Method,
            Path:      c.Request.URL.Path,
            IP:        c.ClientIP(),
            UserAgent: c.Request.UserAgent(),
            Status:    c.Writer.Status(),
            Latency:   latency,
        })
    }
}
```

## 性能优化建议

1. **批量写入**：使用 `BatchCreateLogs` 批量写入日志，提高性能
2. **异步记录**：使用 channel 和 goroutine 异步记录日志，避免阻塞主流程
3. **索引分片**：根据数据量调整分片数量（当前设置为3）
4. **定期归档**：使用 ILM 策略定期归档旧数据
5. **合理分词**：对于不需要全文搜索的字段使用 keyword 类型

## 常见问题

### Q: 搜索结果为空？
A: 确保索引已刷新，ES 默认 5 秒刷新一次，可以等待几秒或手动刷新索引。

### Q: 如何查看 ES 中的数据？
A: 可以使用 Kibana 或直接访问 ES API：
```bash
curl http://localhost:9200/user_action_logs/_search?pretty
```

### Q: 如何修改 Mapping？
A: 修改 `model/system/sys_user_action_log.go` 中的 `GetESMapping` 方法，然后删除旧索引重新创建。

### Q: 日志量太大怎么办？
A:
1. 使用时间分区索引（如 user_action_logs-2024.01）
2. 设置 ILM 策略自动删除旧数据
3. 使用冷热分离架构

## 扩展功能

可以基于此示例扩展以下功能：

1. **日志审计**：记录敏感操作的详细信息
2. **异常告警**：监控错误日志并发送告警
3. **用户行为分析**：分析用户操作习惯和热点功能
4. **性能监控**：统计接口响应时间
5. **安全审计**：检测异常登录和可疑操作

## 参考资源

- [Elasticsearch 官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [gin-vue-admin 文档](https://www.gin-vue-admin.com/)
- [ES Go 客户端文档](https://github.com/elastic/go-elasticsearch)
