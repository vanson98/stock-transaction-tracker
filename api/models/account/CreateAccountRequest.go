package account_model

type CreateAccountRequest struct {
	ChannelName string `json:"channel_name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Currency    string `json:"currency" binding:"required,currency"`
}
