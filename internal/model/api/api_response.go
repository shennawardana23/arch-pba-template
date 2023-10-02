package api

type WebResponse struct {
	Code   int
	Status string
	Data   interface{}
}

type ApiResponse struct {
	Header HeaderResponse `json:"header"`
	Data   interface{}    `json:"data"`
	Error  interface{}    `json:"error"`
}

type HeaderResponse struct {
	ServerTimeMs  int64  `json:"serverTimeMs"`
	ProcessTimeMs int64  `json:"processTimeMs"`
	RequestId     string `json:"requestId"`
}

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Errors  interface{} `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return e.Message.(string)
}

type ErrorValidate struct {
	Key     string `json:"key"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
