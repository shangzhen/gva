package system

import (
	"github.com/gin-gonic/gin"
)

type UserActionLogRouter struct{}

func (s *UserActionLogRouter) InitUserActionLogRouter(Router *gin.RouterGroup) {
	userActionLogRouter := Router.Group("userActionLog")
	{
		userActionLogRouter.POST("initIndex", userActionLogApi.InitIndex)             // 初始化索引
		userActionLogRouter.POST("createLog", userActionLogApi.CreateLog)             // 创建日志
		userActionLogRouter.GET("getLog/:id", userActionLogApi.GetLog)                // 获取单条日志
		userActionLogRouter.POST("searchLogs", userActionLogApi.SearchLogs)           // 搜索日志
		userActionLogRouter.DELETE("deleteLog/:id", userActionLogApi.DeleteLog)       // 删除日志
		userActionLogRouter.POST("getStats", userActionLogApi.GetStats)               // 获取统计
		userActionLogRouter.DELETE("deleteIndex", userActionLogApi.DeleteIndex)       // 删除索引（危险操作）
		userActionLogRouter.POST("batchCreateTestData", userActionLogApi.BatchCreateTestData) // 批量创建测试数据
	}
}
