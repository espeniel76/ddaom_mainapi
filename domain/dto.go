package domain

import "net/http"

type CommonRequest struct {
	Vars       map[string]string
	Parameters map[string]interface{}
	HttpRquest *http.Request
	JWToken    string
}

type CommonResponse struct {
	ResultCode string      `json:"result_code"`
	ErrorDesc  string      `json:"error_desc"`
	Data       interface{} `json:"data" default:"[]"`
}
