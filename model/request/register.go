package request

type Register struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}