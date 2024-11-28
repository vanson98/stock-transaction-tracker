package transaction_model

import (
	db "stt/database/postgres/sqlc"
)

type CreateTransactionModel struct {
	AccountId    int64                `json:"account_id" binding:"required,min=1"`
	InvestmentID int64                `json:"investment_id" binding:"required"`
	TradingDate  string               `json:"trading_date" binding:"required"`
	Trade        string               `json:"trade" binding:"trade"`
	Volume       int64                `json:"volume" binding:"required"`
	OrderPrice   int64                `json:"order_price" binding:"required"`
	MatchVolume  int64                `json:"match_volume" binding:"required"`
	MatchPrice   int64                `json:"match_price" binding:"required"`
	Fee          int64                `json:"fee" binding:"min=0"`
	Tax          int64                `json:"tax" binding:"min=0"`
	Status       db.TransactionStatus `json:"status" binding:"required"`
}
