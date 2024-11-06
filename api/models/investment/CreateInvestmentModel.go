package investment_model

type CreateInvestmentModel struct {
	AccountID   int64  `json:"account_id" binding:"required,min=1"`
	Ticker      string `json:"ticker" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	MarketPrice int64  `json:"market_price" binding:"required"`
	Description string `json:"description"`
}
