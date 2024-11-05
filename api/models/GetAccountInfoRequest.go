package apimodels

type GetAccountInfoRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}
