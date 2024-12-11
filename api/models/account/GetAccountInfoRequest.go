package account_model

type GetAccountInfoRequest struct {
	Ids []int64 `form:"ids" binding:"required"`
}
