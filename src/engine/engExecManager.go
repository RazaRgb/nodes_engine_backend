package engine

import (
	"backend/src/models"
)

var State E_State

func ExecutionManager(tree models.Tree) (models.Tree, error) {
	err := CreateState(tree)
	if err != nil {
		return tree, err
	}

	return tree, nil
}

func CreateState(tree models.Tree) error {
	State = E_State{
		NodeMap: make(map[string]E_Node),
		AdjList: make(map[E_SocketReference][]E_SocketReference),
		FuncMap: map[string]func([]E_Socket, []E_Socket) ([]E_Socket, error){
			"mathAdd": ResolveMathAdd,
		},
	}

	for _, node := range tree.Nodes {
		State.NodeMap[node.ID] = E_Node{
			ID:         node.ID,
			NodeType:   node.Type,
			InpSockArr: make([]E_SocketReference, 0),
			OutSockArr: make([]E_SocketReference, 0),
		}
	}

	for _, edge := range tree.Edges {

		sourceRef := E_SocketReference{
			NodeID:   edge.Source,
			SocketID: edge.SourceHandle,
		}
		targetRef := E_SocketReference{
			NodeID:   edge.Target,
			SocketID: edge.TargetHandle,
		}

		State.AdjList[sourceRef] = append(State.AdjList[sourceRef], targetRef)
	}

	State.FuncMap = map[string]func([]E_Socket, []E_Socket) ([]E_Socket, error){
		"mathAdd": ResolveMathAdd,
	}

	return nil
}
