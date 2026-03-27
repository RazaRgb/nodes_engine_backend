package engine

import (
	"backend/src/models"
)

func ExecuteTree(tree models.Tree) (models.Tree, error){
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
