package transaction_model

type GetTransactionsPagingModel struct {
	AccountId int64  `form:"account_id" binding:"required"`
	Ticker    string `form:"ticker"`
	OrderBy   string `form:"order_by"`
	OrderType string `form:"order_type"`
	Page      int32  `form:"page"`
	PageSize  int32  `form:"page_size"`
}
