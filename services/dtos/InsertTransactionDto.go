package dtos

import (
	db "stt/database/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type InsertTransactionDto struct {
	Ticker      string
	TradingDate pgtype.Timestamp
	Trade       db.TradeType
	Volume      int64
	MatchVolume int64
	OrderPrice  int64
	MatchPrice  int64
	MatchValue  int64
	Fee         int64
	Tax         int64
}
