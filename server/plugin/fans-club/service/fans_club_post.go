package service

import (
	"encoding/json"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/response"
)

type FansClubPostService struct{}

// CreatePost 创建动态
func (s *FansClubPostService) CreatePost(req request.CreatePostReq, userId uint) (post model.FansClubPost, err error) {
	// 检查是否为粉丝团成员
	var member model.FansClubMember
	if err = global.GVA_DB.Where("club_id = ? AND user_id = ?", req.ClubID, userId).First(&member).Error; err != nil {
		return post, errors.New("只有成员才能发布动态")
	}

	// 序列化图片数组
	imagesJSON, _ := json.Marshal(req.Images)

	post = model.FansClubPost{
		ClubID:       req.ClubID,
		UserID:       userId,
		Content:      req.Content,
		Images:       string(imagesJSON),
		LikeCount:    0,
		CommentCount: 0,
	}

	err = global.GVA_DB.Create(&post).Error
	return
}

// UpdatePost 更新动态
func (s *FansClubPostService) UpdatePost(req request.UpdatePostReq, userId uint) error {
	var post model.FansClubPost
	if err := global.GVA_DB.First(&post, req.ID).Error; err != nil {
		return errors.New("动态不存在")
	}

	// 只能修改自己的动态
	if post.UserID != userId {
		return errors.New("无权限修改")
	}

	// 序列化图片数组
	imagesJSON, _ := json.Marshal(req.Images)

	updates := map[string]interface{}{
		"content": req.Content,
		"images":  string(imagesJSON),
	}

	return global.GVA_DB.Model(&post).Updates(updates).Error
}

// DeletePost 删除动态
func (s *FansClubPostService) DeletePost(id uint, userId uint) error {
	var post model.FansClubPost
	if err := global.GVA_DB.First(&post, id).Error; err != nil {
		return errors.New("动态不存在")
	}

	// 检查权限：只能删除自己的动态，或者团长/管理员可以删除
	if post.UserID != userId {
		var member model.FansClubMember
		if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", post.ClubID, userId).First(&member).Error; err != nil {
			return errors.New("无权限删除")
		}
		if member.Role != "owner" && member.Role != "admin" {
			return errors.New("无权限删除")
		}
	}

	return global.GVA_DB.Delete(&post).Error
}

// GetPostList 获取动态列表
func (s *FansClubPostService) GetPostList(req request.PostSearch, userId uint) (resp response.PostListResponse, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.GVA_DB.Model(&model.FansClubPost{})

	if req.ClubID > 0 {
		db = db.Where("club_id = ?", req.ClubID)
	}

	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}

	var posts []model.FansClubPost
	err = db.Count(&resp.Total).Error
	if err != nil {
		return
	}

	err = db.Preload("Club").Limit(limit).Offset(offset).Order("id desc").Find(&posts).Error
	if err != nil {
		return
	}

	// TODO: 可以关联查询用户表，获取发布者信息
	// 这里暂时不实现，因为需要知道用户表的结构

	resp.List = posts
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	return
}

// GetPost 获取动态详情
func (s *FansClubPostService) GetPost(id uint) (post model.FansClubPost, err error) {
	err = global.GVA_DB.Preload("Club").First(&post, id).Error
	return
}

// LikePost 点赞动态
func (s *FansClubPostService) LikePost(id uint, userId uint) error {
	var post model.FansClubPost
	if err := global.GVA_DB.First(&post, id).Error; err != nil {
		return errors.New("动态不存在")
	}

	// 检查是否为粉丝团成员
	var member model.FansClubMember
	if err := global.GVA_DB.Where("club_id = ? AND user_id = ?", post.ClubID, userId).First(&member).Error; err != nil {
		return errors.New("只有成员才能点赞")
	}

	// 简化实现：直接增加点赞数，实际应该记录点赞关系避免重复点赞
	return global.GVA_DB.Model(&post).UpdateColumn("like_count", post.LikeCount+1).Error
}
