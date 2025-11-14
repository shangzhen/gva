package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
)

// PostListResponse 动态列表响应
type PostListResponse struct {
	List     []model.FansClubPost `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}
