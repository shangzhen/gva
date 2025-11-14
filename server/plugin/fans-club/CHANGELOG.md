# 更新日志 (Changelog)

## [v1.0.1] - 2024-01-14

### 修复 (Fixed)
- 修复 `initialize/api.go` - 更正API注册方式
  - 改用 `model.SysApi` 结构体数组
  - 使用 `utils.RegisterApis()` 注册API
  - 添加 `context.Context` 参数
  - 为API添加分组（粉丝团、粉丝团成员、粉丝团动态）

- 修复 `initialize/menu.go` - 更正菜单注册方式
  - 改用 `model.SysBaseMenu` 结构体数组
  - 使用 `utils.RegisterMenus()` 注册菜单
  - 添加 `context.Context` 参数
  - 设置 ParentId 为 24（插件菜单组）
  - 简化菜单结构，改为平级菜单

- 修复 `initialize/gorm.go` - 添加 context 参数
  - 函数签名改为 `Gorm(ctx context.Context)`

- 修复 `plugin.go` - 更新插件注册逻辑
  - 简化结构体命名为 `plugin`
  - 更新初始化函数调用，传入 context
  - 移除不必要的接口方法

### 变更说明

#### 之前的错误实现
```go
// api.go - 错误
func Api() model.PluginApiGroup {
    return model.PluginApiGroup{
        {
            Method:      "POST",
            Path:        "fansClub/createFansClub",
            Description: "创建粉丝团",
            ApiFunc:     api.ApiGroupApp.CreateFansClub,
        },
        // ...
    }
}
```

#### 正确的实现
```go
// api.go - 正确
func Api(ctx context.Context) {
    entities := []model.SysApi{
        {
            Path:        "/fansClub/createFansClub",
            Description: "创建粉丝团",
            ApiGroup:    "粉丝团",
            Method:      "POST",
        },
        // ...
    }
    utils.RegisterApis(entities...)
}
```

#### 菜单变更
- 之前使用嵌套的菜单结构（Children）
- 现在改为平级菜单结构，所有菜单都设置 ParentId 为 24
- 这样更符合 gin-vue-admin 的插件规范

### 技术细节

1. **API 路径修正**
   - 所有路径添加前导斜杠 `/`
   - 例如：`fansClub/createFansClub` → `/fansClub/createFansClub`

2. **API 分组**
   - 粉丝团：6个API
   - 粉丝团成员：5个API
   - 粉丝团动态：6个API

3. **菜单排序**
   - 粉丝团（根菜单）：Sort = 10
   - 粉丝团列表：Sort = 11
   - 我的粉丝团：Sort = 12
   - 粉丝团详情：Sort = 13（隐藏）

### 兼容性
- ✅ 与 gin-vue-admin v2.x 完全兼容
- ✅ 遵循标准插件开发规范
- ✅ 与 announcement 和 email 插件实现一致

---

## [v1.0.0] - 2024-01-14

### 新增 (Added)
- 初始版本发布
- 粉丝团管理功能
- 成员管理功能
- 动态发布功能
- 完整的前后端实现
- 详细文档
