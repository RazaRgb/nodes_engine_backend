package handlers

import (
	"backend/src/db"
	"backend/src/models"
	"backend/src/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"net/http"
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

	err := utils.JsonRead(r, &newProj)
	if err != nil {
		http.Error(w, "unable to create project", http.StatusInternalServerError)
		return
	}

	newProj.Owner = email
	newProj.ID = uuid.NewString()
	// //change tree id when implementing multi-tree projects
	// newTree.ID = newProj.ID
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
	json.NewEncoder(w).Encode(
		struct{
		Project models.Project `json:"project"`
		TreeID string `json:"tree_id"`
	}{
			Project: newProj,
			TreeID: newTree.ID,
		})
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

func HandleGETProjectData(w http.ResponseWriter, r *http.Request) {
	projID := r.PathValue("id")

	treeListany, err := db.RunInTransactionWithReturn(func(tx pgx.Tx) (any, error) {
		var treeList []models.Tree
		treeIDs, err := db.GetTreeIDsForProject(projID)
		if err != nil {
			return nil, err
		}

		for _, treeID := range treeIDs {
			tree, err := db.GetTreeFromDB(treeID, tx)
			fmt.Printf("tree: \n %+v \n", tree)
			if err != nil {
				return tree, err
			}
			treeList = append(treeList, tree)
		}
		if len(treeList) == 0 {
			return treeList, fmt.Errorf("no trees found for the projectID")
		}
		return treeList, nil
	})
	if err != nil {
		fmt.Printf("error occured while running transaction \n %+v\n", err)
		http.Error(w, "unable to get project data", http.StatusInternalServerError)
		return
	}

	treeList, ok := treeListany.([]models.Tree)
	if !ok {
		fmt.Println("Tree list isnt a list of trees?")
		fmt.Printf("treelist : \n %+v\n", treeList)
		http.Error(w, "unable to get project data ", http.StatusInternalServerError)
		return
	}

	responseStruct := struct {
		ProjID string        `json:"project_id"`
		Trees  []models.Tree `json:"tree_list"`
	}{
		ProjID: projID,
		Trees:  treeList,
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JsonWrite(w, responseStruct, http.StatusOK)
}

func HandlePUTProjectData(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to get projects", http.StatusInternalServerError)
		return
	}

	requestStruct := struct {
		ProjID string        `json:"project_id"`
		Trees  []models.Tree `json:"tree_list"`
	}{}

	err := utils.JsonRead(r, &requestStruct)
	if err != nil {
		http.Error(w, "unable to parse request", http.StatusBadRequest)
		return
	}
	fmt.Printf("\nprinting req struct  %+v \n\n", requestStruct)

	found, err := db.MatchProjectWithEmail(requestStruct.ProjID, email)
	if err != nil {
		http.Error(w, "an error occured", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "unauthorized transaction", http.StatusUnauthorized)
		return
	}

	err = db.RunInTransaction(func(tx pgx.Tx) error {

		treeIDList, err := db.GetTreeIDsForProject(requestStruct.ProjID, tx)
		if err != nil {
			return err
		}

		for _, treeID := range treeIDList {
			err := db.ClearTreeContent(treeID, tx)
			if err != nil {
				return err
			}
		}
		for _, tree := range requestStruct.Trees {
			err := db.InsertTreeContentInDB(tree, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, "unable to save project", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
