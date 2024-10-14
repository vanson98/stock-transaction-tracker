package apimodels

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	FullName string `json:"fullname"`
	Password string `json:"password" binding:"required"`
}
