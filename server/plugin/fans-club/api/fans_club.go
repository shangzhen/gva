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

type FansClubApi struct{}

var fansClubService = service.ServiceGroupApp.FansClubService

// CreateFansClub 创建粉丝团
// @Tags FansClub
// @Summary 创建粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateFansClubReq true "创建粉丝团"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /fansClub/createFansClub [post]
func (a *FansClubApi) CreateFansClub(c *gin.Context) {
	var req request.CreateFansClubReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	club, err := fansClubService.CreateFansClub(req, userId)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	response.OkWithData(club, c)
}

// UpdateFansClub 更新粉丝团
// @Tags FansClub
// @Summary 更新粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateFansClubReq true "更新粉丝团"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /fansClub/updateFansClub [put]
func (a *FansClubApi) UpdateFansClub(c *gin.Context) {
	var req request.UpdateFansClubReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = fansClubService.UpdateFansClub(req, userId)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteFansClub 删除粉丝团
// @Tags FansClub
// @Summary 删除粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "删除粉丝团"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /fansClub/deleteFansClub [delete]
func (a *FansClubApi) DeleteFansClub(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = fansClubService.DeleteFansClub(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetFansClub 获取粉丝团详情
// @Tags FansClub
// @Summary 获取粉丝团详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.GetById true "获取粉丝团详情"
// @Success 200 {object} response.Response{data=response.FansClubResponse} "获取成功"
// @Router /fansClub/getFansClub [get]
func (a *FansClubApi) GetFansClub(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	club, err := fansClubService.GetFansClub(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(club, c)
}

// GetFansClubList 获取粉丝团列表
// @Tags FansClub
// @Summary 获取粉丝团列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.FansClubSearch true "获取粉丝团列表"
// @Success 200 {object} response.Response{data=response.FansClubListResponse} "获取成功"
// @Router /fansClub/getFansClubList [get]
func (a *FansClubApi) GetFansClubList(c *gin.Context) {
	var req request.FansClubSearch
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
	list, err := fansClubService.GetFansClubList(req, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}

// GetMyClubs 获取我的粉丝团
// @Tags FansClub
// @Summary 获取我的粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取我的粉丝团"
// @Success 200 {object} response.Response{data=response.FansClubListResponse} "获取成功"
// @Router /fansClub/getMyClubs [get]
func (a *FansClubApi) GetMyClubs(c *gin.Context) {
	var req request.PageInfo
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
	list, err := fansClubService.GetMyClubs(userId, req)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}
