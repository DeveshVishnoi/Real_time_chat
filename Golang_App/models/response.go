package models

// APIResponseStruct is universal API template for sending API Response
type APIResponseStruct struct {
	Code     int         `json:"code"`
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}
type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	Data       any    `json:"data,omitempty"`
}
