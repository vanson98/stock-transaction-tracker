package investment_model

type SearchInvestmentModel struct {
	AccountId  int64  `form:"account_id" binding:"required,min=1"`
	SearchText string `form:"search_text"`
	OrderBy    string `form:"order_by"`
	SortType   string `form:"sort_type"`
	Page       int32  `form:"page" binding:"min=1"`
	PageSize   int32  `form:"page_size" binding:"required"`
}
