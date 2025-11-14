package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FansClubPostApi struct{}

var postService = service.ServiceGroupApp.FansClubPostService

// CreatePost 创建动态
// @Tags FansClubPost
// @Summary 创建动态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreatePostReq true "创建动态"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /fansClubPost/createPost [post]
func (a *FansClubPostApi) CreatePost(c *gin.Context) {
	var req request.CreatePostReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	post, err := postService.CreatePost(req, userId)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	response.OkWithData(post, c)
}

// UpdatePost 更新动态
// @Tags FansClubPost
// @Summary 更新动态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdatePostReq true "更新动态"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /fansClubPost/updatePost [put]
func (a *FansClubPostApi) UpdatePost(c *gin.Context) {
	var req request.UpdatePostReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = postService.UpdatePost(req, userId)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeletePost 删除动态
// @Tags FansClubPost
// @Summary 删除动态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "删除动态"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /fansClubPost/deletePost [delete]
func (a *FansClubPostApi) DeletePost(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = postService.DeletePost(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetPostList 获取动态列表
// @Tags FansClubPost
// @Summary 获取动态列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PostSearch true "获取动态列表"
// @Success 200 {object} response.Response{data=response.PostListResponse} "获取成功"
// @Router /fansClubPost/getPostList [get]
func (a *FansClubPostApi) GetPostList(c *gin.Context) {
	var req request.PostSearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	userId := utils.GetUserID(c)
	list, err := postService.GetPostList(req, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}

// GetPost 获取动态详情
// @Tags FansClubPost
// @Summary 获取动态详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.GetById true "获取动态详情"
// @Success 200 {object} response.Response{data=model.FansClubPost} "获取成功"
// @Router /fansClubPost/getPost [get]
func (a *FansClubPostApi) GetPost(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	post, err := postService.GetPost(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(post, c)
}

// LikePost 点赞动态
// @Tags FansClubPost
// @Summary 点赞动态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "点赞动态"
// @Success 200 {object} response.Response{msg=string} "点赞成功"
// @Router /fansClubPost/likePost [post]
func (a *FansClubPostApi) LikePost(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = postService.LikePost(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("点赞失败!", zap.Error(err))
		response.FailWithMessage("点赞失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("点赞成功", c)
}
