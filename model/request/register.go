package request

type Register struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	ConfirmPassword string `json:"confirm_password"`
}