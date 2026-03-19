package handlers

import (
	"backend/src/db"
	"backend/src/models"
	"backend/src/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"github.com/jackc/pgx/v5"
)

func HandleGETProjects(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	projList, err := db.GetProjectsInfo(email)
	fmt.Printf("projlist un marshalled: \n%+v\n", projList)
	if err != nil {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(projList)
}

func HandlePOSTProject(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	newProj := models.Project{}
	newTree := models.Tree{}

	err := utils.JsonRead(w, r, &newProj)
	if err != nil {
		http.Error(w, "unable to create project", http.StatusInternalServerError)
		return
	}

	newProj.Owner = email
	newProj.ID = uuid.NewString()
	newTree.ID = uuid.NewString()

	// Run transaction
	err = db.RunInTransaction(func(tx pgx.Tx) error {
		newProj, err = db.InsertProject(newProj, tx)
		if err != nil {
			return err
		}
		err = db.CreateNewTree(newTree, newProj.ID, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		http.Error(w, "unable to create new project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newProj)
}

func HandleDELETEProject(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	projID := r.PathValue("id")

	err := db.DeleteProject(projID, email)
	if err != nil {
		http.Error(w, "Unable to delete project", http.StatusInternalServerError)
		fmt.Printf("Errored while deleting project %+v \n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
