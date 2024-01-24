package request

type AddWorker struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type EditWorker struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Username string `json:"username" form:"username"`
}