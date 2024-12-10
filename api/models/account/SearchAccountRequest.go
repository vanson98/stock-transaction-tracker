package account_model

type SearchAccountRequest struct {
	Onwer string `form:"owner" binding:"required"`
}
