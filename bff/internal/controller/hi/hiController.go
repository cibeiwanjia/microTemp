package hi

type HttpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HiController struct{}
