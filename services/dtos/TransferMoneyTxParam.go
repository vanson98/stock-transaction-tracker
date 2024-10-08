package dtos

import db "stt/database/postgres/sqlc"

type TransferMoneyTxParam struct {
	AccountID int64        `json:"account"`
	Amount    int64        `json:"amount"`
	EntryType db.EntryType `json:"entry_from_type"`
}

type TransferMoneyTxResult struct {
	UpdatedAccount db.Account `json:"account"`
	AccountEntry   db.Entry   `json:"account_entry"`
}
