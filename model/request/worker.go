package request

type AddWorker struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type EditWorker struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
}