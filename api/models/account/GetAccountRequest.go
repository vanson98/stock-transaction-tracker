package account_model

type GetAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}