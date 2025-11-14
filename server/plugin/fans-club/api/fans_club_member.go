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

type FansClubMemberApi struct{}

var memberService = service.ServiceGroupApp.FansClubMemberService

// JoinClub 加入粉丝团
// @Tags FansClubMember
// @Summary 加入粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.JoinClubReq true "加入粉丝团"
// @Success 200 {object} response.Response{msg=string} "加入成功"
// @Router /fansClubMember/joinClub [post]
func (a *FansClubMemberApi) JoinClub(c *gin.Context) {
	var req request.JoinClubReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = memberService.JoinClub(req.ClubID, userId)
	if err != nil {
		global.GVA_LOG.Error("加入失败!", zap.Error(err))
		response.FailWithMessage("加入失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("加入成功", c)
}

// QuitClub 退出粉丝团
// @Tags FansClubMember
// @Summary 退出粉丝团
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "退出粉丝团"
// @Success 200 {object} response.Response{msg=string} "退出成功"
// @Router /fansClubMember/quitClub [post]
func (a *FansClubMemberApi) QuitClub(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = memberService.QuitClub(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("退出失败!", zap.Error(err))
		response.FailWithMessage("退出失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("退出成功", c)
}

// GetMemberList 获取成员列表
// @Tags FansClubMember
// @Summary 获取成员列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.MemberSearch true "获取成员列表"
// @Success 200 {object} response.Response{data=response.MemberListResponse} "获取成功"
// @Router /fansClubMember/getMemberList [get]
func (a *FansClubMemberApi) GetMemberList(c *gin.Context) {
	var req request.MemberSearch
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

	list, err := memberService.GetMemberList(req)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}

// UpdateMemberRole 更新成员角色
// @Tags FansClubMember
// @Summary 更新成员角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateMemberRoleReq true "更新成员角色"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /fansClubMember/updateMemberRole [put]
func (a *FansClubMemberApi) UpdateMemberRole(c *gin.Context) {
	var req request.UpdateMemberRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = memberService.UpdateMemberRole(req, userId)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// RemoveMember 移除成员
// @Tags FansClubMember
// @Summary 移除成员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "移除成员"
// @Success 200 {object} response.Response{msg=string} "移除成功"
// @Router /fansClubMember/removeMember [delete]
func (a *FansClubMemberApi) RemoveMember(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	err = memberService.RemoveMember(req.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("移除失败!", zap.Error(err))
		response.FailWithMessage("移除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("移除成功", c)
}
