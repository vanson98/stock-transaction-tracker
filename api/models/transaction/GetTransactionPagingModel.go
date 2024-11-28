package transaction_model

type GetTransactionsPagingModel struct {
	AccountId int64  `form:"account_id"`
	Ticker    string `form:"ticker"`
	Page      int32  `form:"page"`
	PageSize  int32  `form:"page_size"`
}
