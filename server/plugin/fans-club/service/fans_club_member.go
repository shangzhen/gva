package service

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/response"
	"gorm.io/gorm"
)

type FansClubMemberService struct{}

// JoinClub 加入粉丝团
func (s *FansClubMemberService) JoinClub(clubId uint, userId uint) error {
	// 检查粉丝团是否存在
	var club model.FansClub
	if err := global.GVA_DB.First(&club, clubId).Error; err != nil {
		return errors.New("粉丝团不存在")
	}

	// 检查是否已加入
	var count int64
	global.GVA_DB.Model(&model.FansClubMember{}).Where("club_id = ? AND user_id = ?", clubId, userId).Count(&count)
	if count > 0 {
		return errors.New("已经是成员了")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建成员记录
		member := model.FansClubMember{
			ClubID:   clubId,
			UserID:   userId,
			Role:     "member",
			Level:    1,
			Points:   0,
			JoinedAt: time.Now(),
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		// 更新粉丝团成员数量
		if err := tx.Model(&club).UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// QuitClub 退出粉丝团
func (s *FansClubMemberService) QuitClub(clubId uint, userId uint) error {
	var member model.FansClubMember
	if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", clubId, userId).First(&member).Error; err != nil {
		return errors.New("不是该粉丝团成员")
	}

	// 团长不能退出
	if member.Role == "owner" {
		return errors.New("团长不能退出，请先转让团长或删除粉丝团")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除成员记录
		if err := tx.Delete(&member).Error; err != nil {
			return err
		}

		// 更新粉丝团成员数量
		var club model.FansClub
		if err := tx.First(&club, clubId).Error; err != nil {
			return err
		}
		if err := tx.Model(&club).UpdateColumn("member_count", gorm.Expr("member_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetMemberList 获取成员列表
func (s *FansClubMemberService) GetMemberList(req request.MemberSearch) (resp response.MemberListResponse, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.GVA_DB.Model(&model.FansClubMember{})

	if req.ClubID > 0 {
		db = db.Where("club_id = ?", req.ClubID)
	}

	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}

	if req.Role != "" {
		db = db.Where("role = ?", req.Role)
	}

	var members []model.FansClubMember
	err = db.Count(&resp.Total).Error
	if err != nil {
		return
	}

	err = db.Preload("Club").Limit(limit).Offset(offset).Order("id desc").Find(&members).Error
	if err != nil {
		return
	}

	// TODO: 可以关联查询用户表，获取用户名等信息
	// 这里暂时不实现，因为需要知道用户表的结构

	resp.List = members
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	return
}

// UpdateMemberRole 更新成员角色
func (s *FansClubMemberService) UpdateMemberRole(req request.UpdateMemberRoleReq, operatorUserId uint) error {
	var member model.FansClubMember
	if err := global.GVA_DB.First(&member, req.ID).Error; err != nil {
		return errors.New("成员不存在")
	}

	// 检查操作者权限：必须是团长
	var operatorMember model.FansClubMember
	if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", member.ClubID, operatorUserId).First(&operatorMember).Error; err != nil {
		return errors.New("无权限操作")
	}

	if operatorMember.Role != "owner" {
		return errors.New("只有团长可以修改成员角色")
	}

	// 不能修改团长角色
	if member.Role == "owner" {
		return errors.New("不能修改团长角色")
	}

	return global.GVA_DB.Model(&member).Update("role", req.Role).Error
}

// RemoveMember 移除成员
func (s *FansClubMemberService) RemoveMember(memberId uint, operatorUserId uint) error {
	var member model.FansClubMember
	if err := global.GVA_DB.First(&member, memberId).Error; err != nil {
		return errors.New("成员不存在")
	}

	// 检查操作者权限：必须是团长或管理员
	var operatorMember model.FansClubMember
	if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", member.ClubID, operatorUserId).First(&operatorMember).Error; err != nil {
		return errors.New("无权限操作")
	}

	if operatorMember.Role != "owner" && operatorMember.Role != "admin" {
		return errors.New("只有团长和管理员可以移除成员")
	}

	// 不能移除团长
	if member.Role == "owner" {
		return errors.New("不能移除团长")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除成员记录
		if err := tx.Delete(&member).Error; err != nil {
			return err
		}

		// 更新粉丝团成员数量
		var club model.FansClub
		if err := tx.First(&club, member.ClubID).Error; err != nil {
			return err
		}
		if err := tx.Model(&club).UpdateColumn("member_count", gorm.Expr("member_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckMembership 检查是否为成员
func (s *FansClubMemberService) CheckMembership(clubId uint, userId uint) (isMember bool, role string, err error) {
	var member model.FansClubMember
	if err = global.GVA_DB.Where("club_id = ? AND user_id = ?", clubId, userId).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", err
	}
	return true, member.Role, nil
}
