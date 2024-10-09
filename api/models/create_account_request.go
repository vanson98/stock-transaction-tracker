package apimodels

type CreateAccountRequest struct {
	ChannelName string `json:"channel_name" binding:"required"`
	Owner       string `json:"owner" binding:"required"`
	Currency    string `json:"currency" binding:"required,oneof=USD VND EUR"`
}
