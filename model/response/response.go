package response

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data"`
}

type ResponseArray struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data"`
	Length  int         `json:"length"`
}

type ResponseTest struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ResponsePendapatan struct {
	PendapatanKotor  float64 `json:"pendapatan_kotor"`
	PendapatanBersih float64 `json:"pendapatan_bersih"`
}
