package investment_model

type UpdateMarketPriceRequestModel struct {
	InvestmentId int64 `json:"investment_id" binding:"required"`
	MarketPrice  int64 `json:"market_price" binding:"gte=0"`
}
