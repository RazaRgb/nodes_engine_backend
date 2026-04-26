package engine

import (
	"backend/src/models"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func createState(tree models.Tree) (e_State, error) {

	var state e_State
	state = e_State{
		NodeMap:     make(map[string]*e_Node),
		DegMap:      make(map[string]int),
		AdjList:     make(map[e_SocketReference][]e_SocketReference),
		SockMap:     make(map[e_SocketReference]*e_Socket, 100),
		nodeCounter: 0,
	}

	//defer func() {
	//	printEngineState(&state)
	//	fmt.Printf("sockRefMap dump : \n %+v \n", state.SockMap)
	//}()

	for _, node := range tree.Nodes {
		nc, err := getNodeConfig(node.Type)
		if err != nil {
			return e_State{}, err
		}

		if nc.inputCount != -1 && nc.outputCount != -1 {
			state.NodeMap[node.ID] = &e_Node{
				ID:         node.ID,
				NodeType:   node.Type,
				InpSockArr: make([]e_SocketReference, nc.inputCount),
				OutSockArr: make([]e_SocketReference, nc.outputCount),
			}
		} else {
			state.NodeMap[node.ID] = &e_Node{
				ID:         node.ID,
				NodeType:   node.Type,
				InpSockArr: make([]e_SocketReference, 0),
				OutSockArr: make([]e_SocketReference, 0),
			}
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
		state.nodeCounter++
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
		sockPopulate, found := getNodePopulateFunc(node.Type)
		if found {
			err := sockPopulate(&state, state.NodeMap[node.ID], node.Data.Value)
			if err != nil {
				fmt.Printf("error while inserting values in inputNodeSocket: %+v\n", err)
				return e_State{}, err
			}
			//fmt.Printf("Node %s populated : %+v\n", node.Type, *state.NodeMap[node.ID])
		}
		sockServiceCheck, found := getNodeServiceCheckFunc(node.Type)
		if found {
			err := sockServiceCheck(&state, state.NodeMap[node.ID])
			if err != nil {
				fmt.Printf("error while inserting values in inputNodeSocket: %+v\n", err)
				return e_State{}, err
			}
			//fmt.Printf("Node %s populated : %+v\n", node.Type, *state.NodeMap[node.ID])
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
