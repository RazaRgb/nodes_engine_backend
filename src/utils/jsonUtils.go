package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonWrite(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Printf("Error encoding json:\n %v", err)
		http.Error(w, "error encoding json ", http.StatusInternalServerError)
	}
}

func JsonRead(w http.ResponseWriter, req *http.Request, updatedData any) error {

	err := json.NewDecoder(req.Body).Decode(updatedData)
	if err != nil {
		fmt.Printf("Error decoding json:\n %v \n", err)
		return err
	}
	return nil
}
