package system

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserActionLogApi struct{}

// InitIndex
// @Tags      UserActionLog
// @Summary   初始化用户操作日志索引
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "初始化成功"
// @Router    /userActionLog/initIndex [post]
func (u *UserActionLogApi) InitIndex(c *gin.Context) {
	if err := userActionLogService.InitIndex(); err != nil {
		global.GVA_LOG.Error("初始化索引失败!", zap.Error(err))
		response.FailWithMessage("初始化索引失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("初始化索引成功", c)
}

// CreateLog
// @Tags      UserActionLog
// @Summary   创建用户操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UserActionLogCreate           true  "创建日志"
// @Success   200   {object}  response.Response{msg=string}         "创建成功"
// @Router    /userActionLog/createLog [post]
func (u *UserActionLogApi) CreateLog(c *gin.Context) {
	var req request.UserActionLogCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := userActionLogService.CreateLog(&req); err != nil {
		global.GVA_LOG.Error("创建日志失败!", zap.Error(err))
		response.FailWithMessage("创建日志失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetLog
// @Tags      UserActionLog
// @Summary   获取单条日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      string                                              true   "日志ID"
// @Success   200  {object}  response.Response{data=system.UserActionLog,msg=string}  "获取成功"
// @Router    /userActionLog/getLog/:id [get]
func (u *UserActionLogApi) GetLog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("日志ID不能为空", c)
		return
	}

	log, err := userActionLogService.GetLog(id)
	if err != nil {
		global.GVA_LOG.Error("获取日志失败!", zap.Error(err))
		response.FailWithMessage("获取日志失败:"+err.Error(), c)
		return
	}
	response.OkWithData(log, c)
}

// SearchLogs
// @Tags      UserActionLog
// @Summary   搜索用户操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UserActionLogSearch                                true  "搜索条件"
// @Success   200   {object}  response.Response{data=response.UserActionLogListResponse,msg=string}  "搜索成功"
// @Router    /userActionLog/searchLogs [post]
func (u *UserActionLogApi) SearchLogs(c *gin.Context) {
	var req request.UserActionLogSearch
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	result, err := userActionLogService.SearchLogs(&req)
	if err != nil {
		global.GVA_LOG.Error("搜索日志失败!", zap.Error(err))
		response.FailWithMessage("搜索日志失败:"+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// DeleteLog
// @Tags      UserActionLog
// @Summary   删除日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      string                         true  "日志ID"
// @Success   200  {object}  response.Response{msg=string}  "删除成功"
// @Router    /userActionLog/deleteLog/:id [delete]
func (u *UserActionLogApi) DeleteLog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("日志ID不能为空", c)
		return
	}

	if err := userActionLogService.DeleteLog(id); err != nil {
		global.GVA_LOG.Error("删除日志失败!", zap.Error(err))
		response.FailWithMessage("删除日志失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetStats
// @Tags      UserActionLog
// @Summary   获取日志统计
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UserActionLogStats                                  true  "统计条件"
// @Success   200   {object}  response.Response{data=response.UserActionLogStatsResponse,msg=string}  "统计成功"
// @Router    /userActionLog/getStats [post]
func (u *UserActionLogApi) GetStats(c *gin.Context) {
	var req request.UserActionLogStats
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := userActionLogService.GetStats(&req)
	if err != nil {
		global.GVA_LOG.Error("统计失败!", zap.Error(err))
		response.FailWithMessage("统计失败:"+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// DeleteIndex
// @Tags      UserActionLog
// @Summary   删除索引（危险操作）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "删除成功"
// @Router    /userActionLog/deleteIndex [delete]
func (u *UserActionLogApi) DeleteIndex(c *gin.Context) {
	// 可以添加额外的权限验证
	if err := userActionLogService.DeleteIndex(); err != nil {
		global.GVA_LOG.Error("删除索引失败!", zap.Error(err))
		response.FailWithMessage("删除索引失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除索引成功", c)
}

// BatchCreateTestData
// @Tags      UserActionLog
// @Summary   批量创建测试数据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     count  query     int                            false  "生成数量"  default(100)
// @Success   200    {object}  response.Response{msg=string}  "创建成功"
// @Router    /userActionLog/batchCreateTestData [post]
func (u *UserActionLogApi) BatchCreateTestData(c *gin.Context) {
	// 获取生成数量
	countStr := c.DefaultQuery("count", "100")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 || count > 1000 {
		response.FailWithMessage("数量必须在1-1000之间", c)
		return
	}

	// 生成测试数据（这里简化处理，实际应该由前端或测试脚本调用）
	response.OkWithMessage("请使用测试脚本生成测试数据", c)
}
