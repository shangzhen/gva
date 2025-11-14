package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// Menu 初始化菜单注册
func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  24,
			Path:      "fansClub",
			Name:      "fansClub",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      10,
			Meta: model.Meta{
				Title: "粉丝团",
				Icon:  "user-group",
			},
		},
		{
			ParentId:  24,
			Path:      "fansClubList",
			Name:      "fansClubList",
			Hidden:    false,
			Component: "plugin/fans-club/view/club/list.vue",
			Sort:      11,
			Meta: model.Meta{
				Title:     "粉丝团列表",
				Icon:      "list",
				KeepAlive: true,
			},
		},
		{
			ParentId:  24,
			Path:      "myClubs",
			Name:      "myClubs",
			Hidden:    false,
			Component: "plugin/fans-club/view/club/my.vue",
			Sort:      12,
			Meta: model.Meta{
				Title:     "我的粉丝团",
				Icon:      "user",
				KeepAlive: true,
			},
		},
		{
			ParentId:  24,
			Path:      "clubDetail/:id",
			Name:      "clubDetail",
			Hidden:    true,
			Component: "plugin/fans-club/view/club/detail.vue",
			Sort:      13,
			Meta: model.Meta{
				Title: "粉丝团详情",
				Icon:  "detail",
			},
		},
	}
	utils.RegisterMenus(entities...)
}
