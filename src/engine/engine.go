package engine

import (
	"backend/src/models"

	"github.com/google/uuid"
	//"fmt"
)

func ExecuteTree(tree models.Tree, userID uuid.UUID) (struct {
	Sockmap map[string]e_Socket `json:"SockMap"`
}, error) {

	//fmt.Printf("tree to exec : \n %+v \n", tree)

	State, err := createState(tree, userID)
	if err != nil {
		return struct {
			Sockmap map[string]e_Socket `json:"SockMap"`
		}{}, err
	}

	resultSockMap, err := executionManager(&State)
	if err != nil {
		return struct {
			Sockmap map[string]e_Socket `json:"SockMap"`
		}{}, err
	}

	result := struct {
		Sockmap map[string]e_Socket `json:"SockMap"`
	}{
		Sockmap: resultSockMap,
	}

	return result, nil
}
