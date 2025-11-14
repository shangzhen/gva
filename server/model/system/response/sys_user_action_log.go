package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// UserActionLogResponse 用户操作日志响应
type UserActionLogResponse struct {
	system.UserActionLog
}

// UserActionLogListResponse 用户操作日志列表响应
type UserActionLogListResponse struct {
	List  []system.UserActionLog `json:"list"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	PageSize int                 `json:"pageSize"`
}

// UserActionLogStatsResponse 用户操作日志统计响应
type UserActionLogStatsResponse struct {
	Total int64                  `json:"total"`       // 总数
	Stats []map[string]interface{} `json:"stats"`     // 统计数据
}
