package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/response"
	"gorm.io/gorm"
)

type FansClubService struct{}

// CreateFansClub 创建粉丝团
func (s *FansClubService) CreateFansClub(req request.CreateFansClubReq, userId uint) (club model.FansClub, err error) {
	club = model.FansClub{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		OwnerID:     userId,
		MemberCount: 1,
		Level:       1,
		Status:      1,
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建粉丝团
		if err := tx.Create(&club).Error; err != nil {
			return err
		}

		// 自动加入创建者为团长
		member := model.FansClubMember{
			ClubID:   club.ID,
			UserID:   userId,
			Role:     "owner",
			Level:    1,
			Points:   0,
			JoinedAt: time.Now(),
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		return nil
	})

	return
}

// UpdateFansClub 更新粉丝团
func (s *FansClubService) UpdateFansClub(req request.UpdateFansClubReq, userId uint) error {
	var club model.FansClub
	if err := global.GVA_DB.First(&club, req.ID).Error; err != nil {
		return errors.New("粉丝团不存在")
	}

	// 检查权限：只有团长可以修改
	if club.OwnerID != userId {
		return errors.New("无权限修改")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"avatar":      req.Avatar,
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	return global.GVA_DB.Model(&club).Updates(updates).Error
}

// DeleteFansClub 删除粉丝团
func (s *FansClubService) DeleteFansClub(id uint, userId uint) error {
	var club model.FansClub
	if err := global.GVA_DB.First(&club, id).Error; err != nil {
		return errors.New("粉丝团不存在")
	}

	// 检查权限：只有团长可以删除
	if club.OwnerID != userId {
		return errors.New("无权限删除")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除所有成员
		if err := tx.Where("club_id = ?", id).Delete(&model.FansClubMember{}).Error; err != nil {
			return err
		}

		// 删除所有动态
		if err := tx.Where("club_id = ?", id).Delete(&model.FansClubPost{}).Error; err != nil {
			return err
		}

		// 删除粉丝团
		if err := tx.Delete(&club).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetFansClub 获取粉丝团详情
func (s *FansClubService) GetFansClub(id uint, userId uint) (resp response.FansClubResponse, err error) {
	var club model.FansClub
	if err = global.GVA_DB.First(&club, id).Error; err != nil {
		return
	}

	resp.Club = club

	// 检查用户是否为成员
	var member model.FansClubMember
	if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", id, userId).First(&member).Error; err == nil {
		resp.IsMember = true
		resp.Role = member.Role
		if member.Role == "owner" {
			resp.IsOwner = true
		}
	}

	return
}

// GetFansClubList 获取粉丝团列表
func (s *FansClubService) GetFansClubList(req request.FansClubSearch, userId uint) (resp response.FansClubListResponse, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	cache, err := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("fans_club_search_%d_%s_%d", userId, req.Keyword, req.Status)).Result()
	if cache != "" {
		err = json.Unmarshal([]byte(cache), &resp.List)
		if err == nil {
			resp.PageSize = req.PageSize
			resp.Page = req.Page
			return
		}
	}

	db := global.GVA_DB.Model(&model.FansClub{})

	// 关键词搜索
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	var clubs []model.FansClub
	err = db.Count(&resp.Total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&clubs).Error
	if err != nil {
		return
	}

	// 查询用户的成员关系
	var memberMap = make(map[uint]model.FansClubMember)
	if userId > 0 {
		var members []model.FansClubMember
		clubIds := make([]uint, len(clubs))
		for i, club := range clubs {
			clubIds[i] = club.ID
		}
		global.GVA_DB.Where("user_id = ? AND club_id IN ?", userId, clubIds).Find(&members)
		for _, member := range members {
			memberMap[member.ClubID] = member
		}
	}

	// 组装响应
	resp.List = make([]response.FansClubResponse, len(clubs))
	for i, club := range clubs {
		resp.List[i].Club = club
		if member, ok := memberMap[club.ID]; ok {
			resp.List[i].IsMember = true
			resp.List[i].Role = member.Role
			if member.Role == "owner" {
				resp.List[i].IsOwner = true
			}
		}
	}

	listByte, _ := json.Marshal(resp.List)
	global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("fans_club_search_%d_%s_%d", userId, req.Keyword, req.Status), string(listByte), time.Minute)

	resp.Page = req.Page
	resp.PageSize = req.PageSize

	return
}

// GetMyClubs 获取我的粉丝团列表
func (s *FansClubService) GetMyClubs(userId uint, req request.PageInfo) (resp response.FansClubListResponse, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// 查询用户加入的粉丝团ID列表
	var members []model.FansClubMember
	db := global.GVA_DB.Where("user_id = ?", userId)

	err = db.Model(&model.FansClubMember{}).Count(&resp.Total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&members).Error
	if err != nil {
		return
	}

	// 查询粉丝团信息
	clubIds := make([]uint, len(members))
	memberMap := make(map[uint]model.FansClubMember)
	for i, member := range members {
		clubIds[i] = member.ClubID
		memberMap[member.ClubID] = member
	}

	var clubs []model.FansClub
	if len(clubIds) > 0 {
		err = global.GVA_DB.Where("id IN ?", clubIds).Find(&clubs).Error
		if err != nil {
			return
		}
	}

	// 组装响应
	resp.List = make([]response.FansClubResponse, len(clubs))
	for i, club := range clubs {
		resp.List[i].Club = club
		if member, ok := memberMap[club.ID]; ok {
			resp.List[i].IsMember = true
			resp.List[i].Role = member.Role
			if member.Role == "owner" {
				resp.List[i].IsOwner = true
			}
		}
	}

	resp.Page = req.Page
	resp.PageSize = req.PageSize

	return
}
