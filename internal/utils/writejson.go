package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data any, statusCode int) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"Failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil{
		fmt.Println("Error writing response: ", err)
	}
} 