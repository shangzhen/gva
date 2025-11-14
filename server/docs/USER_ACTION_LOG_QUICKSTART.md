# Elasticsearch 用户操作日志 - 快速开始

## 📝 功能概述

这是一个完整的 Elasticsearch 集成示例，实现了用户操作日志系统，包括：

- ✅ 完整的 ES Mapping 定义（14个字段）
- ✅ 日志的增删改查操作
- ✅ 多条件组合搜索（用户、模块、动作、时间范围等）
- ✅ 分页和排序
- ✅ 聚合统计（按动作、模块、用户等分组）
- ✅ RESTful API 接口（8个接口）
- ✅ 完整的测试脚本

## 🚀 快速开始

### 1. 启动 Elasticsearch

```bash
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  elasticsearch:8.19.0
```

### 2. 配置项目

编辑 `server/config.yaml`：

```yaml
system:
    use-elasticsearch: true

elasticsearch:
    addresses:
        - http://127.0.0.1:9200
    index: user_action_logs
    max-retries: 3
    timeout: 30
```

### 3. 启动项目

```bash
cd server
go run main.go
```

### 4. 测试接口

#### 4.1 初始化索引

```bash
curl -X POST http://localhost:8888/userActionLog/initIndex \
  -H "x-token: your-jwt-token"
```

#### 4.2 创建日志

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

#### 4.3 搜索日志

```bash
curl -X POST http://localhost:8888/userActionLog/searchLogs \
  -H "x-token: your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "pageSize": 10
  }'
```

#### 4.4 获取统计

```bash
curl -X POST http://localhost:8888/userActionLog/getStats \
  -H "x-token: your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "start_time": "2024-01-01T00:00:00Z",
    "end_time": "2024-12-31T23:59:59Z",
    "group_by": "action"
  }'
```

## 📁 文件结构

```
server/
├── model/system/
│   ├── sys_user_action_log.go              # 数据模型 + ES Mapping
│   ├── request/sys_user_action_log.go      # 请求参数
│   └── response/sys_user_action_log.go     # 响应数据
├── service/system/
│   └── sys_user_action_log.go              # 业务逻辑（10个方法）
├── api/v1/system/
│   └── sys_user_action_log.go              # API处理器（8个接口）
├── router/system/
│   └── sys_user_action_log.go              # 路由定义
├── test/
│   └── elasticsearch_test.go               # 测试脚本
└── docs/
    └── USER_ACTION_LOG_ES.md               # 详细文档
```

## 🔌 API 接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /userActionLog/initIndex | 初始化索引 |
| POST | /userActionLog/createLog | 创建日志 |
| GET | /userActionLog/getLog/:id | 获取日志 |
| POST | /userActionLog/searchLogs | 搜索日志 |
| DELETE | /userActionLog/deleteLog/:id | 删除日志 |
| POST | /userActionLog/getStats | 获取统计 |
| DELETE | /userActionLog/deleteIndex | 删除索引 |
| POST | /userActionLog/batchCreateTestData | 批量创建测试数据 |

## 📊 ES Mapping 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| id | keyword | 日志ID |
| user_id | long | 用户ID |
| username | keyword | 用户名 |
| action | keyword | 操作动作 |
| module | keyword | 操作模块 |
| method | keyword | HTTP方法 |
| path | text | 请求路径 |
| ip | ip | IP地址 |
| user_agent | text | 用户代理 |
| status | integer | 状态码 |
| latency | long | 响应时间 |
| request | text | 请求参数 |
| response | text | 响应数据 |
| error_msg | text | 错误信息 |
| create_time | date | 创建时间 |

## 🔍 搜索功能

支持以下搜索条件：

- **精确匹配**: user_id, action, module, method, ip, status
- **模糊搜索**: username
- **范围查询**: start_time, end_time
- **关键词搜索**: keyword（搜索path和error_msg）
- **分页**: page, pageSize
- **排序**: order_field, order_type

## 📈 统计功能

支持按以下字段分组统计：

- action - 操作动作
- module - 操作模块
- user_id - 用户ID
- username - 用户名
- method - 请求方法
- status - 响应状态码

## 🧪 运行测试

```bash
cd server/test
# 修改 elasticsearch_test.go 中的 token 为实际的 JWT token
go run elasticsearch_test.go
```

测试脚本会自动执行：
1. 初始化索引
2. 创建单条日志
3. 批量创建50条日志
4. 获取单条日志
5. 搜索所有日志
6. 按用户ID搜索
7. 按操作动作搜索
8. 按模块搜索
9. 按时间范围搜索
10. 获取统计数据

## 💡 使用建议

1. **异步记录**: 使用 goroutine 异步记录日志，避免影响主流程
2. **批量写入**: 对于大量日志使用批量接口提高性能
3. **定期归档**: 设置 ILM 策略定期删除或归档旧数据
4. **监控告警**: 监控错误日志并设置告警
5. **索引分片**: 根据数据量调整分片数量

## 📖 详细文档

完整的使用文档请查看：[USER_ACTION_LOG_ES.md](./USER_ACTION_LOG_ES.md)

包含：
- 详细的 API 文档
- 代码示例
- 性能优化建议
- 常见问题解答
- 扩展功能建议

## 🎯 核心代码示例

### 创建日志

```go
import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
)

userActionLogService := service.ServiceGroupApp.SystemServiceGroup.UserActionLogService

logReq := &request.UserActionLogCreate{
    UserID:    1,
    Username:  "admin",
    Action:    "login",
    Module:    "user",
    Method:    "POST",
    Path:      "/base/login",
    IP:        "127.0.0.1",
    UserAgent: "Mozilla/5.0",
    Status:    200,
    Latency:   123,
}

err := userActionLogService.CreateLog(logReq)
```

### 搜索日志

```go
searchReq := &request.UserActionLogSearch{
    PageInfo: request.PageInfo{
        Page:     1,
        PageSize: 10,
    },
    Action:    "login",
    StartTime: "2024-01-01T00:00:00Z",
    EndTime:   "2024-12-31T23:59:59Z",
}

result, err := userActionLogService.SearchLogs(searchReq)
```

## ⚠️ 注意事项

1. **JWT Token**: 所有接口都需要 JWT 认证，请先登录获取 token
2. **ES 版本**: 推荐使用 Elasticsearch 8.x 版本
3. **索引刷新**: ES 默认 5 秒刷新一次，创建后需等待几秒才能搜索到
4. **数据备份**: 定期备份重要数据
5. **权限控制**: 生产环境请启用 ES 的安全认证

## 🤝 贡献

如有问题或建议，欢迎提交 Issue 或 PR！

## 📄 许可证

MIT License
