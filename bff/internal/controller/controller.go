package controller

import "github.com/cibeiwanjia/microTemp/bff/internal/controller/hi"

type HttpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Controller struct{}

func NewHiController() *hi.HiController {
	return &hi.HiController{}
}
