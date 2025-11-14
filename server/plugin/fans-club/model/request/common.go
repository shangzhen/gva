package request

// PageInfo 分页参数
type PageInfo struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// GetById 通过ID获取
type GetById struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// IdsReq ID列表请求
type IdsReq struct {
	Ids []uint `json:"ids" binding:"required"`
}
