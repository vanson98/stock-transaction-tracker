package investment_model

type SearchInvestmentRequestModel struct {
	AccountIds []int64 `form:"account_ids[]" binding:"required"`
	SearchText string  `form:"search_text"`
	OrderBy    string  `form:"order_by"`
	SortType   string  `form:"sort_type"`
	Page       int32   `form:"page" binding:"min=1"`
	PageSize   int32   `form:"page_size" binding:"required"`
}
