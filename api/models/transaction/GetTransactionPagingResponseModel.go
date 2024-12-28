package transaction_model

import db "stt/database/postgres/sqlc"

type GetTransactionsPagingResponseModel struct {
	Transactions  []db.GetTransactionsPagingRow `json:"transactions"`
	TotalRow      int32                         `json:"total"`
	SumMatchValue int64                         `json:"sum_match_value"`
	SumFee        int64                         `json:"sum_fee"`
	SumTax        int64                         `json:"sum_tax"`
	SumReturn     int64                         `json:"sum_return"`
}
