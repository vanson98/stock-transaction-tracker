package transaction_model

type AddTransactionRequest struct {
	AccountId    int64  `json:"account_id" binding:"required"`
	Ticker       string `json:"ticker" binding:"required"`
	TradingDate  string `json:"trading_date" binding:"required"`
	Trade        string `json:"trade" binding:"required"`
	Volume       int64  `json:"volume" binding:"required"`
	OrderPrice   int64  `json:"order_price" binding:"required"`
	MatchVolume  int64  `json:"match_volume" binding:"required"`
	MatchPrice   int64  `json:"match_price" binding:"required"`
	MatchValue   int64  `json:"match_value" binding:"required"`
	Fee          int64  `json:"fee" binding:"min=0"`
	Tax          int64  `json:"tax" binding:"min=0" `
	Cost         int64  `json:"cost" binding:"required"`
	Return       int64  `json:"return"`
	UploadStatus string `json:"upload_status" binding:"required"`
	CheckCost    bool   `json:"check_cost" binding:"required"`
}
