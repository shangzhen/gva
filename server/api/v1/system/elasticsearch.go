package system

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ElasticsearchApi struct{}

// 搜索请求结构
type EsSearchRequest struct {
	Index    string   `json:"index" binding:"required"`
	Keyword  string   `json:"keyword" binding:"required"`
	Fields   []string `json:"fields"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

// 索引文档请求
type EsIndexRequest struct {
	Index string      `json:"index" binding:"required"`
	DocID string      `json:"doc_id" binding:"required"`
	Data  interface{} `json:"data" binding:"required"`
}

// 删除文档请求
type EsDeleteRequest struct {
	Index string `json:"index" binding:"required"`
	DocID string `json:"doc_id" binding:"required"`
}

// 获取文档请求
type EsGetRequest struct {
	Index string `json:"index" binding:"required"`
	DocID string `json:"doc_id" binding:"required"`
}

// 更新文档请求
type EsUpdateRequest struct {
	Index   string                 `json:"index" binding:"required"`
	DocID   string                 `json:"doc_id" binding:"required"`
	Updates map[string]interface{} `json:"updates" binding:"required"`
}

// 搜索文档
// @Tags      Elasticsearch
// @Summary   全文搜索
// @Produce   application/json
// @Param     data  body      EsSearchRequest  true  "搜索条件"
// @Success   200   {object}  response.Response
// @Router    /elasticsearch/search [post]
func (e *ElasticsearchApi) Search(c *gin.Context) {
	var req EsSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if global.GVA_ES == nil {
		response.FailWithMessage("Elasticsearch未初始化", c)
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	result, err := system.ElasticsearchServiceApp.SimpleSearch(
		ctx,
		req.Index,
		req.Keyword,
		req.Fields,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		global.GVA_LOG.Error("搜索失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// 索引文档
// @Tags      Elasticsearch
// @Summary   索引单个文档
// @Produce   application/json
// @Param     data  body      EsIndexRequest  true  "文档数据"
// @Success   200   {object}  response.Response
// @Router    /elasticsearch/index [post]
func (e *ElasticsearchApi) IndexDocument(c *gin.Context) {
	var req EsIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if global.GVA_ES == nil {
		response.FailWithMessage("Elasticsearch未初始化", c)
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	err := system.ElasticsearchServiceApp.IndexDocument(ctx, req.Index, req.DocID, req.Data)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}

// 删除文档
// @Tags      Elasticsearch
// @Summary   删除文档
// @Produce   application/json
// @Param     data  body      EsDeleteRequest  true  "删除参数"
// @Success   200   {object}  response.Response
// @Router    /elasticsearch/delete [post]
func (e *ElasticsearchApi) DeleteDocument(c *gin.Context) {
	var req EsDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if global.GVA_ES == nil {
		response.FailWithMessage("Elasticsearch未初始化", c)
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	err := system.ElasticsearchServiceApp.DeleteDocument(ctx, req.Index, req.DocID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}

// 获取文档
// @Tags      Elasticsearch
// @Summary   获取单个文档
// @Produce   application/json
// @Param     data  body      EsGetRequest  true  "获取参数"
// @Success   200   {object}  response.Response
// @Router    /elasticsearch/get [post]
func (e *ElasticsearchApi) GetDocument(c *gin.Context) {
	var req EsGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if global.GVA_ES == nil {
		response.FailWithMessage("Elasticsearch未初始化", c)
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	doc, err := system.ElasticsearchServiceApp.GetDocument(ctx, req.Index, req.DocID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(doc, c)
}

// 更新文档
// @Tags      Elasticsearch
// @Summary   更新文档
// @Produce   application/json
// @Param     data  body      EsUpdateRequest  true  "更新数据"
// @Success   200   {object}  response.Response
// @Router    /elasticsearch/update [post]
func (e *ElasticsearchApi) UpdateDocument(c *gin.Context) {
	var req EsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if global.GVA_ES == nil {
		response.FailWithMessage("Elasticsearch未初始化", c)
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	err := system.ElasticsearchServiceApp.UpdateDocument(ctx, req.Index, req.DocID, req.Updates)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}
