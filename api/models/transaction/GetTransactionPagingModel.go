package transaction_model

type GetTransactionsPagingModel struct {
	AccountIds []int64 `form:"account_ids" binding:"required"`
	Ticker     string  `form:"ticker"`
	OrderBy    string  `form:"order_by"`
	OrderType  string  `form:"order_type"`
	Page       int32   `form:"page"`
	PageSize   int32   `form:"page_size"`
}
