# 粉丝团插件集成指南

本文档说明如何将粉丝团插件集成到 gin-vue-admin 项目中。

## 后端集成

### 方式一：通过插件系统自动注册（推荐）

如果项目支持插件自动发现，插件会自动注册。

### 方式二：手动注册

#### 1. 在主程序中注册插件

编辑 `server/initialize/plugin.go` 或在主初始化文件中添加：

```go
package initialize

import (
    fansclub "github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club"
    "github.com/gin-gonic/gin"
)

func InstallPlugin(Router *gin.Engine) {
    // ... 其他插件 ...

    // 注册粉丝团插件
    fansclub.Plugin.Register(Router)
}
```

#### 2. 在 main.go 中调用

确保在 `server/main.go` 中调用了插件初始化：

```go
package main

import (
    "github.com/flipped-aurora/gin-vue-admin/server/initialize"
    // ... 其他导入 ...
)

func main() {
    // ... 其他初始化代码 ...

    Router := initialize.Routers()

    // 安装插件
    initialize.InstallPlugin(Router)

    // ... 其他代码 ...
}
```

#### 3. 数据库自动迁移

插件会在注册时自动创建数据库表，无需额外操作。

如果需要手动创建表或添加索引，可以执行：

```bash
cd server/plugin/fans-club/initialize
# 根据你的数据库类型执行 init.sql
mysql -u root -p your_database < init.sql
```

## 前端集成

### 1. 路由配置

前端路由会通过菜单自动生成，无需手动配置。

如果需要手动配置，在 `web/src/router/index.js` 中添加：

```javascript
{
  path: '/fansClub',
  name: 'fansClub',
  component: () => import('@/view/routerHolder.vue'),
  meta: {
    title: '粉丝团',
    icon: 'user-group'
  },
  children: [
    {
      path: 'list',
      name: 'fansClubList',
      component: () => import('@/plugin/fans-club/view/club/list.vue'),
      meta: {
        title: '粉丝团列表',
        icon: 'list',
        keepAlive: true
      }
    },
    {
      path: 'my',
      name: 'myClubs',
      component: () => import('@/plugin/fans-club/view/club/my.vue'),
      meta: {
        title: '我的粉丝团',
        icon: 'user',
        keepAlive: true
      }
    },
    {
      path: 'detail/:id',
      name: 'clubDetail',
      component: () => import('@/plugin/fans-club/view/club/detail.vue'),
      meta: {
        title: '粉丝团详情',
        hidden: true
      }
    }
  ]
}
```

### 2. API Base URL 配置

确保 `web/src/utils/request.js` 中的 baseURL 配置正确：

```javascript
const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 99999
})
```

## 权限配置

### 1. API 权限

插件API已自动注册到权限系统。管理员需要在"系统管理 -> API管理"中为角色分配权限。

主要API权限：
- 粉丝团管理：创建、编辑、删除、查询
- 成员管理：加入、退出、成员列表、角色管理
- 动态管理：发布、编辑、删除、点赞

### 2. 菜单权限

在"系统管理 -> 菜单管理"中为角色分配菜单权限。

## 环境要求

### 后端
- Go 1.16+
- Gin 框架
- GORM v2
- MySQL 5.7+ / PostgreSQL 9.6+ / SQLite 3

### 前端
- Node.js 14+
- Vue 3
- Element Plus
- Vite

## 配置项

### 后端配置（可选）

如果需要自定义配置，可以在 `server/plugin/fans-club/config/` 目录下添加配置文件。

### 前端配置（可选）

可以在 `.env.development` 或 `.env.production` 中添加环境变量。

## 验证安装

### 1. 启动后端服务

```bash
cd server
go run main.go
```

检查日志输出：
```
fans-club plugin: register table success
```

### 2. 启动前端服务

```bash
cd web
npm run dev
```

### 3. 访问系统

1. 登录系统
2. 在菜单中找到"粉丝团"模块
3. 尝试创建一个粉丝团

## 常见问题

### 1. 插件未注册

**问题**：启动后看不到粉丝团菜单

**解决**：
- 检查插件是否正确注册
- 检查日志中是否有错误信息
- 确认角色有菜单权限

### 2. API 请求失败

**问题**：前端请求返回 404

**解决**：
- 检查路由是否正确注册
- 检查 JWT 认证是否正常
- 查看后端日志确认路由路径

### 3. 数据库表未创建

**问题**：操作时报表不存在

**解决**：
- 确认 GORM AutoMigrate 已执行
- 检查数据库连接是否正常
- 手动执行 init.sql

### 4. 权限错误

**问题**：操作时提示无权限

**解决**：
- 在 API 管理中为角色分配权限
- 在菜单管理中为角色分配菜单
- 确认 JWT token 有效

## 卸载插件

### 1. 移除代码

```bash
rm -rf server/plugin/fans-club
rm -rf web/src/plugin/fans-club
```

### 2. 移除注册

从初始化代码中移除插件注册调用。

### 3. 清理数据库（可选）

```sql
DROP TABLE IF EXISTS gva_fans_club;
DROP TABLE IF EXISTS gva_fans_club_member;
DROP TABLE IF EXISTS gva_fans_club_post;
```

**警告**：这将删除所有粉丝团数据，请谨慎操作！

## 技术支持

如有问题，请查看：
- [README.md](./README.md) - 功能说明
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始
- 项目 Issues

## 更新插件

### 1. 备份数据

```bash
mysqldump -u root -p database_name gva_fans_club gva_fans_club_member gva_fans_club_post > fans_club_backup.sql
```

### 2. 更新代码

替换插件目录下的文件。

### 3. 数据库迁移

GORM 会自动处理新增字段。如有结构性变更，请参考更新日志手动执行迁移脚本。

### 4. 重启服务

```bash
# 后端
cd server && go run main.go

# 前端
cd web && npm run dev
```

## 生产环境部署

### 1. 编译后端

```bash
cd server
go build -o gva main.go
```

### 2. 构建前端

```bash
cd web
npm run build
```

### 3. 配置 Nginx（可选）

```nginx
location /api/ {
    proxy_pass http://localhost:8888/;
}
```

### 4. 性能优化建议

- 启用 Redis 缓存粉丝团信息
- 为高频查询添加索引
- 配置 CDN 加速静态资源
- 启用 GZIP 压缩

## 安全建议

1. **JWT 认证**：所有API都需要有效的JWT token
2. **权限控制**：严格控制角色权限
3. **SQL 注入防护**：使用 GORM 参数化查询
4. **XSS 防护**：前端输出时进行转义
5. **CSRF 防护**：使用 CSRF token
6. **限流**：配置 API 限流防止滥用

## 监控建议

- 记录关键操作日志
- 监控 API 响应时间
- 监控数据库连接数
- 设置错误告警

---

集成完成后，你的系统将拥有完整的粉丝团管理功能！
