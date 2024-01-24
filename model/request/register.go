package request

type Register struct {
	Username        string `json:"username" form:"username" binding:"required"`
	Email           string `json:"email" form:"email" binding:"required,email"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"`
}