package apimodels

type TransferMoneyRequest struct {
	AccountID int64  `json:"account"`
	Amount    int64  `json:"amount"`
	EntryType string `json:"entry_type"`
}
