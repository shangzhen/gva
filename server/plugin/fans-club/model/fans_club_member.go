package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FansClubMember 粉丝团成员表
type FansClubMember struct {
	global.GVA_MODEL
	ClubID   uint      `json:"clubId" gorm:"not null;index;comment:粉丝团ID"`
	UserID   uint      `json:"userId" gorm:"not null;index;comment:用户ID"`
	Role     string    `json:"role" gorm:"type:varchar(20);default:'member';comment:角色:owner-团长,admin-管理员,member-成员"`
	Level    int       `json:"level" gorm:"default:1;comment:成员等级"`
	Points   int       `json:"points" gorm:"default:0;comment:积分"`
	JoinedAt time.Time `json:"joinedAt" gorm:"comment:加入时间"`

	// 关联查询字段
	Club     *FansClub `json:"club,omitempty" gorm:"foreignKey:ClubID"`
	Username string    `json:"username,omitempty" gorm:"-"` // 用户名，不存数据库
}

func (FansClubMember) TableName() string {
	return "gva_fans_club_member"
}
