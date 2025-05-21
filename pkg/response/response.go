package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, success bool, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Success: success,
		Data:    data,
		Error:   errMsg,
	}

	json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, true, data, "")
}

func Fail(w http.ResponseWriter, status int, errMsg string) {
	JSON(w, status, false, nil, errMsg)
}