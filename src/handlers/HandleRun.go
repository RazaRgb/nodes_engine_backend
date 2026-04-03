package handlers

import (
	"backend/src/db"
	"backend/src/engine"
	"backend/src/models"
	"backend/src/utils"
	"fmt"
	"net/http"
)

func HandleRun(w http.ResponseWriter, r *http.Request) {

	requestStruct := struct {
		ProjID string        `json:"project_id"`
		Trees  []models.Tree `json:"tree_list"`
	}{}

	{ // validation and all
		email, ok := utils.GetEmailFromContext(r.Context())
		if !ok {
			http.Error(w, "Unable to get projects", http.StatusInternalServerError)
			return
		}

		err := utils.JsonRead(r, &requestStruct)
		if err != nil {
			http.Error(w, "unable to parse request", http.StatusBadRequest)
			return
		}

		found, err := db.MatchProjectWithEmail(requestStruct.ProjID, email)
		if err != nil {
			http.Error(w, "error in matching email and proj", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(w, "unauthorized transaction", http.StatusForbidden)
			return
		}
	}

	tree := requestStruct.Trees[0]

	result, err := engine.ExecuteTree(tree)
	if err != nil {
		http.Error(w, "error while execution", http.StatusBadRequest)
		fmt.Printf("unable to execute tree: \n %+v \n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JsonWrite(w, result, http.StatusOK)
}

type runResponse struct {
}
