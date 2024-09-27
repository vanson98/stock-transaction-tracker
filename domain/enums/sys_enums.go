package enums

type InvestmentStatus int
type TransactionType int
type TranserType int

const (
	NotActive InvestmentStatus = iota
	Active
	BuyOut
)

const (
	Buy TransactionType = iota
	Sell
)

const (
	Add TranserType = iota
	Withdraw
)
