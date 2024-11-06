package account_model

type TransferMoneyRequest struct {
	AccountID int64  `json:"accountId" binding:"required,min=1"`
	Amount    int64  `json:"amount" binding:"required"`
	EntryType string `json:"entryType" binding:"required,oneof=IT TM"`
	Currency  string `json:"currency" binding:"required,currency"`
}
