package controllers

type Response struct {
	Code    int         `json:"code" `
	Data    interface{} `json:"data" `
	Message string      `json:"message" `
}

func (r *Response) ReturnFailureString(str string) {
	r.Code = -1
	r.Message = str
