package core_http_respone

type ErrorResponse struct {
	Error   string `json:"error" example:"ful error text"`
	Message string `json:"message" example:"short human-readble message"`
}
