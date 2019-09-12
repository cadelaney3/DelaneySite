package utils

import (
	"net/http"
	"encoding/json"
)

func Message(status uint, message string) map[string]interface{} {
	return map[string]interface{} {"status": status, "message": message}
}

func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}