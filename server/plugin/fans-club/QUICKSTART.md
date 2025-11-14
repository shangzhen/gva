# 粉丝团插件 - 快速开始指南

## 1. 启用插件

### 方法一：在主程序中注册（推荐）

在 `server/main.go` 中添加插件注册：

```go
import (
    fansclub "github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club"
)

// 在初始化插件的位置添加
func initPlugin(engine *gin.Engine) {
    // ... 其他插件 ...

    // 注册粉丝团插件
    fansclub.Plugin.Register(engine)
}
```

### 方法二：使用插件管理系统

如果你的项目支持动态插件管理，可以通过插件管理界面安装。

## 2. 初始化数据库

插件会自动创建所需的数据库表，包括：
- `gva_fans_club` - 粉丝团表
- `gva_fans_club_member` - 成员表
- `gva_fans_club_post` - 动态表

启动服务后，GORM会自动迁移表结构。

可选：执行 `initialize/init.sql` 来创建索引优化查询性能。

## 3. 访问前端页面

启动前端项目后，在菜单中找到"粉丝团"模块，包含：
- **粉丝团列表**：浏览所有粉丝团
- **我的粉丝团**：查看我加入的粉丝团

## 4. 基本使用流程

### 4.1 创建粉丝团

1. 进入"粉丝团列表"页面
2. 点击"创建粉丝团"按钮
3. 填写粉丝团信息：
   - 名称（必填）
   - 描述
   - 头像URL
4. 点击"确定"创建

创建成功后，你将自动成为该粉丝团的团长。

### 4.2 加入粉丝团

1. 在"粉丝团列表"中浏览粉丝团
2. 点击"加入"按钮
3. 加入成功后，可以在"我的粉丝团"中查看

### 4.3 管理粉丝团

作为团长，你可以：
1. 编辑粉丝团信息
2. 管理成员角色（设置/取消管理员）
3. 移除成员
4. 删除粉丝团

进入粉丝团详情页面，切换到"成员管理"标签进行操作。

### 4.4 发布动态

1. 进入粉丝团详情页
2. 在"粉丝动态"标签下点击"发布动态"
3. 输入动态内容
4. 点击"发布"

### 4.5 点赞动态

在动态列表中，点击"点赞"按钮即可。

## 5. API调用示例

### 创建粉丝团

```javascript
import { createFansClub } from '@/plugin/fans-club/api/fansClub'

const data = {
  name: '我的粉丝团',
  description: '欢迎加入我的粉丝团',
  avatar: 'https://example.com/avatar.jpg'
}

const res = await createFansClub(data)
```

### 加入粉丝团

```javascript
import { joinClub } from '@/plugin/fans-club/api/fansClubMember'

const res = await joinClub({ clubId: 1 })
```

### 发布动态

```javascript
import { createPost } from '@/plugin/fans-club/api/fansClubPost'

const data = {
  clubId: 1,
  content: '大家好，这是我的第一条动态！',
  images: []
}

const res = await createPost(data)
```

## 6. 权限说明

### 团长权限
- 编辑粉丝团信息
- 删除粉丝团
- 设置/取消管理员
- 移除任何成员
- 删除任何动态

### 管理员权限
- 移除普通成员
- 删除任何动态

### 普通成员权限
- 发布动态
- 编辑/删除自己的动态
- 点赞动态
- 退出粉丝团

## 7. 常见问题

### Q: 如何退出粉丝团？
A: 进入"我的粉丝团"页面，点击对应粉丝团的"退出"按钮。注意：团长不能退出，需要先删除粉丝团或转让团长。

### Q: 如何删除粉丝团？
A: 只有团长可以删除粉丝团。在粉丝团列表中，点击"删除"按钮。删除粉丝团会同时删除所有成员关系和动态。

### Q: 如何查看粉丝团成员？
A: 进入粉丝团详情页，切换到"成员管理"标签即可查看所有成员。

### Q: 动态支持图片吗？
A: 当前版本已预留图片字段，但需要配置文件上传服务才能使用。

## 8. 下一步

- 查看 [README.md](./README.md) 了解完整功能
- 查看 API 文档了解所有接口
- 根据需求扩展功能（评论、标签、积分商城等）

## 9. 技术支持

如遇到问题，请：
1. 查看服务端日志
2. 检查数据库连接
3. 确认JWT认证正常
4. 查看浏览器控制台错误信息

## 10. 更新日志

### v1.0.0
- 初始版本发布
- 基础粉丝团管理功能
- 成员管理功能
- 动态发布功能
