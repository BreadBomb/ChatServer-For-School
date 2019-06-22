package model

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Data    interface{} `json:"data"`
}
