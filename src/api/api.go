package api

import (
	"net/http"
	"encoding/json"
)

func Error(w http.ResponseWriter, errorMessage string) {
	result, _ := json.Marshal(map[string]interface{}{
		"error": errorMessage,
	})

	w.Write(result)
}