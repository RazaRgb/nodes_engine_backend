package engine

import (
	"backend/src/models"
)


func executionManager(state e_State) (models.Tree, error) {

	return models.Tree{}, nil
}

func createState(tree models.Tree) (e_State ,error) {
	var state e_State
	state = e_State{
		NodeMap: make(map[string]e_Node),
		AdjList: make(map[e_SocketReference][]e_SocketReference),
		FuncMap: map[string]func([]e_Socket, []e_Socket) ([]e_Socket, error){
			//"mathAdd": ResolveMathAdd,
		},
	}

	for _, node := range tree.Nodes {
		state.NodeMap[node.ID] = e_Node{
			ID:         node.ID,
			NodeType:   node.Type,
			InpSockArr: make([]e_SocketReference, 0),
			OutSockArr: make([]e_SocketReference, 0),
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

	state.FuncMap = map[string]func([]e_Socket, []e_Socket) ([]e_Socket, error){
		//"mathAdd": ResolveMathAdd,
	}

	return state, nil
}
