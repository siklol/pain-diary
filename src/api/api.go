package api

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, errorMessage string) {
	result, _ := json.Marshal(map[string]interface{}{
		"error": errorMessage,
	})

	w.Write(result)
}

func Success(w http.ResponseWriter, result map[string]interface{}) {
	result["error"] = ""

	resultJson, _ := json.Marshal(result)

	w.Write(resultJson)
}
