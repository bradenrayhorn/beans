package response

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
