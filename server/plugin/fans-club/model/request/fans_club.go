package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
)

// FansClubSearch 查询条件
type FansClubSearch struct {
	model.FansClub
	PageInfo
}

// CreateFansClubReq 创建粉丝团请求
type CreateFansClubReq struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

// UpdateFansClubReq 更新粉丝团请求
type UpdateFansClubReq struct {
	ID          uint   `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" form:"description"`
	Avatar      string `json:"avatar" form:"avatar"`
	Status      *int   `json:"status" form:"status"`
}

// JoinClubReq 加入粉丝团请求
type JoinClubReq struct {
	ClubID uint `json:"clubId" form:"clubId" binding:"required"`
}

// MemberSearch 成员查询条件
type MemberSearch struct {
	ClubID uint   `json:"clubId" form:"clubId"`
	UserID uint   `json:"userId" form:"userId"`
	Role   string `json:"role" form:"role"`
	PageInfo
}

// UpdateMemberRoleReq 更新成员角色
type UpdateMemberRoleReq struct {
	ID   uint   `json:"id" form:"id" binding:"required"`
	Role string `json:"role" form:"role" binding:"required,oneof=admin member"`
}
