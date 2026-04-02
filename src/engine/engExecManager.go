package engine

import (
	"backend/src/models"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func executionManager(state e_State) (models.Tree, error) {

	return models.Tree{}, nil
}

func createState(tree models.Tree) (e_State, error) {

	var state e_State
	state = e_State{
		NodeMap: make(map[string]*e_Node),
		DegMap:  make(map[string]int),
		AdjList: make(map[e_SocketReference][]e_SocketReference),
		SockMap: make(map[e_SocketReference]*e_Socket, 100),
	}

	defer func() {
		printEngineState(&state)
		fmt.Printf("sockRefMap dump : \n %+v \n", state.SockMap)
	}()

	for _, node := range tree.Nodes {
		nc, err := getNodeConfig(node.Type)
		if err != nil {
			return e_State{}, err
		}

		state.NodeMap[node.ID] = &e_Node{
			ID:         node.ID,
			NodeType:   node.Type,
			InpSockArr: make([]e_SocketReference, nc.inputCount),
			OutSockArr: make([]e_SocketReference, nc.outputCount),
		}

		// Create sockets
		for i := range nc.outputCount {
			sockID := "o" + strconv.Itoa(i+1)
			state.NodeMap[node.ID].OutSockArr[i].SocketID = (sockID)
			state.NodeMap[node.ID].OutSockArr[i].NodeID = node.ID
			state.SockMap[e_SocketReference{
				NodeID:   node.ID,
				SocketID: sockID,
			}] = &e_Socket{
				ID: sockID,
			}
		}
		for i := range nc.inputCount {
			sockID := "i" + strconv.Itoa(i+1)
			state.NodeMap[node.ID].InpSockArr[i].SocketID = (sockID)
			state.NodeMap[node.ID].InpSockArr[i].NodeID = node.ID
			state.SockMap[e_SocketReference{
				NodeID:   node.ID,
				SocketID: sockID,
			}] = &e_Socket{
				ID: sockID,
			}
		}
		state.DegMap[node.ID] = 0

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
		state.DegMap[edge.Target]++

		state.AdjList[sourceRef] = append(state.AdjList[sourceRef], targetRef)
	}

	for _, node := range tree.Nodes {
		if node.Type == "inputNumber" {
			fmt.Printf("state :\n %+v \n", state)
			err := popInputNumber(&state, state.NodeMap[node.ID], node.Data.Value)
			if err != nil {
				fmt.Printf("error while inserting values in inputNodeSocket")
				return e_State{}, err
			}
		}
	}

	return state, nil
}

func printEngineState(state *e_State) {
	fmt.Print("----Engine state---- \n")
	err := json.NewEncoder(os.Stdout).Encode(*state)
	fmt.Println()
	if err != nil {
		fmt.Printf("Error encoding json:\n %v", err)
	}
}
