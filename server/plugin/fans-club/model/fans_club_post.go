package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FansClubPost 粉丝团动态表
type FansClubPost struct {
	global.GVA_MODEL
	ClubID       uint   `json:"clubId" gorm:"not null;index;comment:粉丝团ID" form:"clubId"`
	UserID       uint   `json:"userId" gorm:"not null;comment:发布者ID" form:"userId"`
	Content      string `json:"content" gorm:"type:text;not null;comment:动态内容" form:"content"`
	Images       string `json:"images" gorm:"type:text;comment:图片JSON数组" form:"images"`
	LikeCount    int    `json:"likeCount" gorm:"default:0;comment:点赞数" form:"likeCount"`
	CommentCount int    `json:"commentCount" gorm:"default:0;comment:评论数" form:"commentCount"`

	// 关联查询字段
	Club     *FansClub `json:"club,omitempty" gorm:"foreignKey:ClubID" form:"clubId"`
	Username string    `json:"username,omitempty" gorm:"-" form:"username"` // 发布者用户名
}

func (FansClubPost) TableName() string {
	return "gva_fans_club_post"
}
