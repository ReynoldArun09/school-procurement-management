package utils

type ApiResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data,omitempty"`
}
