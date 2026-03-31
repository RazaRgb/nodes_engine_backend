package engine

import (
	"backend/src/models"
	"encoding/json"
	"fmt"
	"os"
)

func executionManager(state e_State) (models.Tree, error) {

	return models.Tree{}, nil
}

func createState(tree models.Tree) (e_State, error) {
	var state e_State
	state = e_State{
		NodeMap: make(map[string]e_Node),
		DegMap:  make(map[string]int),
		AdjList: make(map[e_SocketReference][]e_SocketReference),
	}

	for _, node := range tree.Nodes {
		nc, err := getNodeConfig(node.Type)
		if err != nil {
			return e_State{}, err
		}

		state.DegMap[node.ID] = nc.inputCount

		state.NodeMap[node.ID] = e_Node{
			ID:         node.ID,
			NodeType:   node.Type,
			InpSockArr: make([]e_SocketReference, nc.inputCount),
			OutSockArr: make([]e_SocketReference, nc.outputCount),
		}
	}

	for _, edge := range tree.Edges {
		sourceRef := e_SocketReference{
			NodeID:   edge.Source,
			SocketID: edge.SourceHandle,
		}
		targetRef := e_SocketReference{
			NodeID:   edge.Target,
			SocketID: edge.TargetHandle,
		}

		state.AdjList[sourceRef] = append(state.AdjList[sourceRef], targetRef)
	}
	fmt.Print("state: \n")
	err := json.NewEncoder(os.Stdout).Encode(state)
	if err != nil {
		fmt.Printf("Error encoding json:\n %v", err)
	}

	return state, nil
}
