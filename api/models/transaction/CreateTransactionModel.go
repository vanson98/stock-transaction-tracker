package transaction_model

import (
	db "stt/database/postgres/sqlc"
)

type CreateTransactionModel struct {
	InvestmentID    int64                `json:"investment_id" binding:"required"`
	Ticker          string               `json:"ticker" binding:"required"`
	TradingDate     string               `json:"trading_date" binding:"required"`
	Trade           db.TradeType         `json:"trade" binding:"required"`
	Volume          int32                `json:"volume" binding:"required"`
	OrderPrice      int64                `json:"order_price" binding:"required"`
	MatchVolume     int32                `json:"match_volume" binding:"required"`
	MatchPrice      int64                `json:"match_price" binding:"required"`
	MatchValue      int64                `json:"match_value" binding:"required"`
	Fee             int32                `json:"fee" binding:"required"`
	Tax             int32                `json:"tax" binding:"required"`
	Cost            int64                `json:"cost" binding:"required"`
	CostOfGoodsSold int64                `json:"cost_of_goods_sold" binding:"required"`
	Return          int64                `json:"return" binding:"required"`
	Status          db.TransactionStatus `json:"status" binding:"required"`
}
