package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FansClub 粉丝团表
type FansClub struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"type:varchar(100);not null;comment:粉丝团名称"`
	Description string `json:"description" form:"description" gorm:"type:text;comment:粉丝团描述"`
	Avatar      string `json:"avatar" form:"avatar" gorm:"type:varchar(255);comment:头像URL"`
	OwnerID     uint   `json:"ownerId" form:"ownerId" gorm:"not null;comment:创建者ID"`
	MemberCount int    `json:"memberCount" form:"memberCount" gorm:"default:1;comment:成员数量"`
	Level       int    `json:"level" form:"level" gorm:"default:1;comment:粉丝团等级"`
	Status      int    `json:"status" form:"status" gorm:"default:1;comment:状态:0-待审核,1-正常,2-禁用"`
}

func (FansClub) TableName() string {
	return "gva_fans_club"
}
