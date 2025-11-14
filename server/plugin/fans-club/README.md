# 粉丝团插件 (Fans Club Plugin)

## 功能简介

粉丝团插件是一个完整的粉丝社群管理系统，支持创建粉丝团、成员管理、动态发布等功能。

## 主要功能

### 1. 粉丝团管理
- 创建粉丝团
- 编辑粉丝团信息
- 删除粉丝团
- 查看粉丝团列表
- 查看粉丝团详情
- 查看我的粉丝团

### 2. 成员管理
- 加入粉丝团
- 退出粉丝团
- 查看成员列表
- 角色管理（团长、管理员、成员）
- 移除成员
- 成员等级和积分系统

### 3. 粉丝动态
- 发布动态
- 编辑动态
- 删除动态
- 查看动态列表
- 点赞动态
- 支持图片上传（预留功能）

## 数据库表结构

### gva_fans_club (粉丝团表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | varchar(100) | 粉丝团名称 |
| description | text | 粉丝团描述 |
| avatar | varchar(255) | 头像URL |
| owner_id | uint | 创建者ID |
| member_count | int | 成员数量 |
| level | int | 粉丝团等级 |
| status | int | 状态（0-待审核，1-正常，2-禁用） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间 |

### gva_fans_club_member (成员表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| club_id | uint | 粉丝团ID |
| user_id | uint | 用户ID |
| role | varchar(20) | 角色（owner-团长，admin-管理员，member-成员） |
| level | int | 成员等级 |
| points | int | 积分 |
| joined_at | timestamp | 加入时间 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间 |

### gva_fans_club_post (动态表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| club_id | uint | 粉丝团ID |
| user_id | uint | 发布者ID |
| content | text | 动态内容 |
| images | text | 图片JSON数组 |
| like_count | int | 点赞数 |
| comment_count | int | 评论数 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间 |

## API接口

### 粉丝团相关
- `POST /fansClub/createFansClub` - 创建粉丝团
- `PUT /fansClub/updateFansClub` - 更新粉丝团
- `DELETE /fansClub/deleteFansClub` - 删除粉丝团
- `GET /fansClub/getFansClub` - 获取粉丝团详情
- `GET /fansClub/getFansClubList` - 获取粉丝团列表
- `GET /fansClub/getMyClubs` - 获取我的粉丝团

### 成员相关
- `POST /fansClubMember/joinClub` - 加入粉丝团
- `POST /fansClubMember/quitClub` - 退出粉丝团
- `GET /fansClubMember/getMemberList` - 获取成员列表
- `PUT /fansClubMember/updateMemberRole` - 更新成员角色
- `DELETE /fansClubMember/removeMember` - 移除成员

### 动态相关
- `POST /fansClubPost/createPost` - 创建动态
- `PUT /fansClubPost/updatePost` - 更新动态
- `DELETE /fansClubPost/deletePost` - 删除动态
- `GET /fansClubPost/getPostList` - 获取动态列表
- `GET /fansClubPost/getPost` - 获取动态详情
- `POST /fansClubPost/likePost` - 点赞动态

## 安装使用

### 1. 后端安装

将插件放置在 `server/plugin/fans-club` 目录下。

在 `server/plugin/plugin.go` 或相应的插件管理文件中注册插件：

```go
import (
    fansclub "github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club"
)

// 在插件注册函数中添加
fansclub.Plugin.Register(engine)
```

### 2. 前端安装

前端文件已放置在 `web/src/plugin/fans-club` 目录下，包含：
- `api/` - API调用
- `view/` - 页面组件

### 3. 数据库初始化

插件会自动创建数据库表。如需手动执行SQL，可以运行 `initialize/init.sql`。

### 4. 菜单配置

插件已自动注册菜单，包含：
- 粉丝团列表
- 我的粉丝团
- 粉丝团详情（隐藏菜单）

## 权限说明

### 角色权限
1. **团长 (owner)**
   - 编辑粉丝团信息
   - 删除粉丝团
   - 管理所有成员
   - 设置/取消管理员
   - 移除成员
   - 删除任何动态

2. **管理员 (admin)**
   - 移除普通成员
   - 删除任何动态

3. **成员 (member)**
   - 发布动态
   - 编辑/删除自己的动态
   - 点赞动态
   - 退出粉丝团

## 目录结构

```
server/plugin/fans-club/
├── api/                    # API处理层
│   ├── fans_club.go
│   ├── fans_club_member.go
│   ├── fans_club_post.go
│   └── enter.go
├── initialize/             # 初始化
│   ├── api.go
│   ├── gorm.go
│   ├── menu.go
│   ├── router.go
│   └── init.sql
├── model/                  # 数据模型
│   ├── request/
│   ├── response/
│   ├── fans_club.go
│   ├── fans_club_member.go
│   └── fans_club_post.go
├── router/                 # 路由
│   ├── fans_club.go
│   └── enter.go
├── service/                # 业务逻辑层
│   ├── fans_club.go
│   ├── fans_club_member.go
│   ├── fans_club_post.go
│   └── enter.go
├── plugin.go               # 插件注册
└── README.md

web/src/plugin/fans-club/
├── api/                    # API调用
│   ├── fansClub.js
│   ├── fansClubMember.js
│   └── fansClubPost.js
└── view/                   # 页面组件
    └── club/
        ├── list.vue        # 粉丝团列表
        ├── my.vue          # 我的粉丝团
        └── detail.vue      # 粉丝团详情
```

## 扩展功能建议

1. **评论系统**：为动态添加评论功能
2. **标签系统**：为粉丝团添加标签分类
3. **积分商城**：使用积分兑换奖励
4. **任务系统**：完成任务获得积分
5. **勋章系统**：成员成就徽章
6. **图片上传**：支持动态图片上传
7. **消息通知**：粉丝团动态通知
8. **数据统计**：粉丝团数据分析

## 注意事项

1. 所有API接口都需要JWT认证
2. 操作记录中间件已集成，会自动记录操作日志
3. 软删除已启用，删除的数据可以恢复
4. 建议为大型粉丝团启用缓存优化查询性能
5. 图片上传功能需要配置文件上传服务

## 技术栈

- **后端**: Go + Gin + GORM
- **前端**: Vue 3 + Element Plus
- **数据库**: MySQL / PostgreSQL / SQLite

## 版本历史

### v1.0.0 (2024-01-14)
- 初始版本
- 基础粉丝团管理功能
- 成员管理功能
- 动态发布功能

## 许可证

MIT License
