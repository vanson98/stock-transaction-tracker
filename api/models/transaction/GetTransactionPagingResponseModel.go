package transaction_model

import db "stt/database/postgres/sqlc"

type GetTransactionsPagingResponseModel struct {
	Transactions []db.GetTransactionsPagingRow `json:"transactions"`
	Total        int32                         `json:"total"`
}
