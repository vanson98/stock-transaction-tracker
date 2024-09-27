package domain

import (
	"stt/domain/enums"
	"time"
)

type Transaction struct {
	Id                   int32
	Price                float32
	Amount               int32
	TransactionTime      time.Time
	UserId               int32
	CreatedTime          time.Time
	LastModifierUserId   int32
	LastModificationTime time.Time
	TransactionType      enums.TransactionType
	InvestmentId         int32
	TotalFee             float32
	CapitalCost          float32
}
