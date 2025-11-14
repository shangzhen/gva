package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/router"
	"github.com/gin-gonic/gin"
)

// Router 初始化路由
func Router(engine *gin.Engine) {
	// 获取私有路由组（需要JWT认证）
	privateGroup := engine.Group("")
	router.RouterGroupApp.InitFansClubRouter(privateGroup)
}
