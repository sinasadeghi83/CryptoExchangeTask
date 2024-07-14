package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func addErrorToBody[T any](body *map[string]interface{}, status int, v T) {
	delete(*body, "data")
	(*body)["status"] = status
	(*body)["message"] = v
}

func EncodeResponse[T any](w http.ResponseWriter, r *http.Request, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	body := map[string]interface{}{
		"status": status,
		"data":   v,
	}

	if status != http.StatusOK {
		addErrorToBody(&body, status, v)
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		addErrorToBody(&body, http.StatusInternalServerError, "Internal Error")
		json.NewEncoder(w).Encode(body)
		fmt.Println(fmt.Errorf("encode json: %w", err))
	}
}

func DecodeJson[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
