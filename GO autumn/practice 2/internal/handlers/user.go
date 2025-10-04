package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"user_id": id})

	case http.MethodPost:
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data["name"] == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid name"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"created": data["name"]})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
