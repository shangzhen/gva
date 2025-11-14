package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// UserActionLogCreate 创建用户操作日志请求
type UserActionLogCreate struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Action    string `json:"action" binding:"required"`
	Module    string `json:"module" binding:"required"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    int    `json:"status"`
	Latency   int64  `json:"latency"`
	Request   string `json:"request"`
	Response  string `json:"response"`
	ErrorMsg  string `json:"error_msg"`
}

// UserActionLogSearch 搜索用户操作日志请求
type UserActionLogSearch struct {
	request.PageInfo
	UserID     *uint   `json:"user_id" form:"user_id"`         // 用户ID
	Username   string  `json:"username" form:"username"`       // 用户名（模糊搜索）
	Action     string  `json:"action" form:"action"`           // 操作动作
	Module     string  `json:"module" form:"module"`           // 操作模块
	Method     string  `json:"method" form:"method"`           // 请求方法
	IP         string  `json:"ip" form:"ip"`                   // IP地址
	Status     *int    `json:"status" form:"status"`           // 响应状态码
	StartTime  string  `json:"start_time" form:"start_time"`   // 开始时间
	EndTime    string  `json:"end_time" form:"end_time"`       // 结束时间
	Keyword    string  `json:"keyword" form:"keyword"`         // 关键词（搜索path和error_msg）
	OrderField string  `json:"order_field" form:"order_field"` // 排序字段
	OrderType  string  `json:"order_type" form:"order_type"`   // 排序类型 asc/desc
}

// UserActionLogStats 用户操作日志统计请求
type UserActionLogStats struct {
	StartTime string `json:"start_time" form:"start_time" binding:"required"` // 开始时间
	EndTime   string `json:"end_time" form:"end_time" binding:"required"`     // 结束时间
	GroupBy   string `json:"group_by" form:"group_by"`                        // 分组字段 action, module, user等
}
