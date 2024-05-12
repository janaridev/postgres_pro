package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMsg string `json:"errorMsg,omitempty"`
}

func NewErrorResponse(errorMsg string) *ErrorResponse {
	return &ErrorResponse{
		ErrorMsg: errorMsg,
	}
}

func SendJSONResponse(w http.ResponseWriter, contentType string, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
