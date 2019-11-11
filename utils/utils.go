package utils

import (
	"encoding/json"
	"net/http"
)

func BuildResponse(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message(false, "Invalid request"))
}

func RespondUnhautorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message(false, "Unhautorized"))
}

func RespondBadRequestWithMessage(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message(false, message))
}

func RespondBasicSuccess(w http.ResponseWriter, r *http.Request) {
	resp := Message(true, "Success")
	Respond(w, resp)
}
