package handlers

import (
	"backend/src/db"
	"backend/src/utils"
	"encoding/json"
	"net/http"
)

func HandleGETProjects(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}
	projList, err := db.GetProjectsInfo(email)
	if err != nil {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(projList)

}

func HandlePOSTProject(w http.ResponseWriter, r *http.Request) {
	// email, ok := utils.GetEmailFromContext(r.Context())
	// if !ok {
	// 	http.Error(w, "Unable to get projects", http.StatusInternalServerError)
	// 	return
	// }

	// json.NewDecoder

}

func HandleDELETEProject(w http.ResponseWriter, r *http.Request) {

}
