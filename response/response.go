package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Body       any    `json:"body,omitempty"`
}

func (r Response) Dump() []byte {
	marshal, _ := json.Marshal(r)
	return marshal
}

func (r Response) IsError() (string, bool) {
	if r.StatusCode == http.StatusInternalServerError {
		return r.Message, true
	}
	return "", false
}

func (r Response) IsNormal() bool {
	return r.StatusCode <= http.StatusNoContent
}

func OK(body any) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func Warn(msg string) *Response {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
	}
}

func Build(code int, msg string, body any) *Response {
	return &Response{
		StatusCode: code,
		Message:    msg,
		Body:       body,
	}
}

func Error(err error) *Response {
	return Warn(err.Error())
}
