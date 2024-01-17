package response

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data"`
}

type ResponseTest struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
