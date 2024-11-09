package sv_interface

import (
	"context"
	"stt/services/dtos"
)

type ITransactionService interface {
	CreateBuyingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) error
}
