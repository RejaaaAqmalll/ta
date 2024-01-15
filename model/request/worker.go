package request

type AddWorker struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditWorker struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}