package dtos

import (
	db "stt/database/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTransactionDto struct {
	AccountId    int64
	InvestmentId int64
	Ticker       string
	TradingDate  pgtype.Timestamp
	Trade        db.TradeType
	Volume       int64
	OrderPrice   int64
	MatchVolume  int64
	MatchPrice   int64
	Fee          int64
	Tax          int64
	Status       string
}
