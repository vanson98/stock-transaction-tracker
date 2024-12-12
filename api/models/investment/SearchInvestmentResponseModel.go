package investment_model

import db "stt/database/postgres/sqlc"

type SearchInvestmentResponseModel struct {
	Investments []db.SearchInvestmentPagingRow `json:"investments"`
	TotalItems  int64                          `json:"total_items"`
}
