package system

import "time"

// UserActionLog 用户操作日志（存储在 ES 中）
type UserActionLog struct {
	ID         string    `json:"id"`                       // 日志ID
	UserID     uint      `json:"user_id"`                  // 用户ID
	Username   string    `json:"username"`                 // 用户名
	Action     string    `json:"action"`                   // 操作动作（login, logout, create, update, delete等）
	Module     string    `json:"module"`                   // 操作模块（user, role, menu等）
	Method     string    `json:"method"`                   // 请求方法（GET, POST, PUT, DELETE）
	Path       string    `json:"path"`                     // 请求路径
	IP         string    `json:"ip"`                       // IP地址
	UserAgent  string    `json:"user_agent"`               // 用户代理
	Status     int       `json:"status"`                   // 响应状态码
	Latency    int64     `json:"latency"`                  // 响应时间（毫秒）
	Request    string    `json:"request,omitempty"`        // 请求参数
	Response   string    `json:"response,omitempty"`       // 响应数据
	ErrorMsg   string    `json:"error_msg,omitempty"`      // 错误信息
	CreateTime time.Time `json:"create_time"`              // 创建时间
}

// GetESIndexName 获取ES索引名称
func (UserActionLog) GetESIndexName() string {
	return "user_action_logs"
}

// GetESMapping 获取ES映射定义
func (UserActionLog) GetESMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"refresh_interval":   "5s",
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"default": map[string]interface{}{
						"type": "standard",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "long",
				},
				"username": map[string]interface{}{
					"type": "keyword",
				},
				"action": map[string]interface{}{
					"type": "keyword",
				},
				"module": map[string]interface{}{
					"type": "keyword",
				},
				"method": map[string]interface{}{
					"type": "keyword",
				},
				"path": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"ip": map[string]interface{}{
					"type": "ip",
				},
				"user_agent": map[string]interface{}{
					"type": "text",
				},
				"status": map[string]interface{}{
					"type": "integer",
				},
				"latency": map[string]interface{}{
					"type": "long",
				},
				"request": map[string]interface{}{
					"type": "text",
					"index": false, // 不索引，只存储
				},
				"response": map[string]interface{}{
					"type": "text",
					"index": false, // 不索引，只存储
				},
				"error_msg": map[string]interface{}{
					"type": "text",
				},
				"create_time": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
			},
		},
	}
}
