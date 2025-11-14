package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
)

// PostSearch 动态查询条件
type PostSearch struct {
	ClubID uint `json:"clubId" form:"clubId"`
	UserID uint `json:"userId" form:"userId"`
	PageInfo
}

// CreatePostReq 创建动态请求
type CreatePostReq struct {
	ClubID  uint     `json:"clubId" binding:"required" form:"clubId"`
	Content string   `json:"content" binding:"required,min=1,max=1000" form:"content"`
	Images  []string `json:"images" form:"images"`
}

// UpdatePostReq 更新动态请求
type UpdatePostReq struct {
	ID      uint     `json:"id" binding:"required" form:"id"`
	Content string   `json:"content" binding:"required,min=1,max=1000" form:"content"`
	Images  []string `json:"images" form:"images"`
}

// PostWithMember 动态及成员信息
type PostWithMember struct {
	model.FansClubPost
	IsMember bool `json:"isMember" form:"isMember"` // 当前用户是否为成员
}
