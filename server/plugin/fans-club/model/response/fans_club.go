package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
)

// FansClubResponse 粉丝团响应
type FansClubResponse struct {
	Club     model.FansClub `json:"club"`
	IsMember bool           `json:"isMember"` // 当前用户是否为成员
	IsOwner  bool           `json:"isOwner"`  // 当前用户是否为团长
	Role     string         `json:"role"`     // 当前用户在粉丝团中的角色
}

// FansClubListResponse 粉丝团列表响应
type FansClubListResponse struct {
	List     []FansClubResponse `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

// MemberResponse 成员响应
type MemberResponse struct {
	model.FansClubMember
}

// MemberListResponse 成员列表响应
type MemberListResponse struct {
	List     []model.FansClubMember `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
}
