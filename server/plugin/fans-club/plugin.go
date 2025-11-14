package fans_club

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/initialize"
	"github.com/gin-gonic/gin"
)

type plugin struct{}

var Plugin = new(plugin)

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()

	// 注册API到系统
	initialize.Api(ctx)

	// 注册菜单到系统
	initialize.Menu(ctx)

	// 初始化数据库表
	initialize.Gorm(ctx)

	// 初始化路由
	initialize.Router(group)
}
