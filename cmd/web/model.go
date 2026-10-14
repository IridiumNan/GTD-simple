package web

import (
	"github.com/IridiumNan/GTD-simple/internal/models"
)

type RequestCapture struct {
	Data RequstData `json:"data"`
}

type RequstData struct {
	S models.Stuff `json:"stuff"`
}

type ResponseCapture struct {
	Msg string `json:"msg"`

	Data any `json:"data"`
}

func NewResponse(msg string, data any) *ResponseCapture {
	return &ResponseCapture{
		Msg:  msg,
		Data: data,
	}
}
