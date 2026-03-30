package engine

import (
	"backend/src/models"
	"fmt"
)

func ExecuteTree(tree models.Tree) (models.Tree, error) {

	fmt.Printf("tree to exec : \n %+v \n", tree)

	State, err := createState(tree)
	if err != nil {
		return tree, err
	}

	resultTree, err := executionManager(State)
	if err != nil {
		return tree, err
	}

	return resultTree, nil
}
