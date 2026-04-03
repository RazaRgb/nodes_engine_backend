package engine

//TODO: implement logging with different log levels

//const LOGLEVEL int = 2

func executionManager(state *e_State) (map[string]e_Socket, error) {
	engComm := e_Communication{
		interrupt:      make(chan error, 3),
		valuePropagate: make(chan e_workerValue, 5),
	}

	outPutSockMap := make(map[string]e_Socket)

	for nodeID, _ := range state.NodeMap {
		if state.DegMap[nodeID] == 0 {
			err := spawnWorker(state, &engComm, nodeID)
			if err != nil {
				return nil, err
			}
		}
	}

	defer close(engComm.interrupt)
	defer close(engComm.valuePropagate)

	for {
		select {
		case err := <-engComm.interrupt:
			if err != nil {
				return nil, err
			} else {
				return outPutSockMap, nil
			}

		case sockStruct := <-engComm.valuePropagate:
			for _, sock := range sockStruct.socket {
				tgtRef := state.AdjList[e_SocketReference{
					NodeID:   sockStruct.nodeID,
					SocketID: sock.ID,
				}]

				for _, tgtSockRef := range tgtRef {
					state.SockMap[tgtSockRef] = &sock
					state.DegMap[tgtSockRef.NodeID]--
					if state.DegMap[tgtSockRef.NodeID] == 0 {
						err := spawnWorker(state, &engComm, tgtSockRef.NodeID)
						if err != nil {
							return nil, err
						}
					}
				}

				outPutSockMap[(sockStruct.nodeID + ":" + sock.ID)] = sock
			}

			//update coutner
			state.nodeCounter--
			if state.nodeCounter == 0 {
				return outPutSockMap, nil
			}
		}
	}
}

func spawnWorker(
	state *e_State,
	engComm *e_Communication,
	nodeID string,
) error {

	nodePtr := state.NodeMap[nodeID]
	resolver, err := getNodeResolver(nodePtr.NodeType)
	if err != nil {
		return err
	}
	outSockArr := make([]e_Socket, 0)
	for _, sockRef := range state.NodeMap[nodeID].OutSockArr {
		outSockArr = append(outSockArr, *state.SockMap[sockRef])
	}
	inpSockArr := make([]e_Socket, 0)
	for _, sockRef := range state.NodeMap[nodeID].InpSockArr {
		inpSockArr = append(inpSockArr, *state.SockMap[sockRef])
	}
	go worker(inpSockArr, outSockArr, resolver, engComm, nodeID)
	return nil
}

func worker(
	inpSockrefArr []e_Socket,
	outSockrefArr []e_Socket,
	resolver nodeResolver,
	engComm *e_Communication,
	nodeID string,
) {

	//{ // LOgging
	//	if LOGLEVEL >= 2 {
	//		fmt.Printf("LOG2: working on %s\n", nodeID)
	//	}
	//}

	outputSock, err := resolver(inpSockrefArr, outSockrefArr)
	if err != nil {
		engComm.interrupt <- err
		return
	}
	engComm.valuePropagate <- struct {
		socket []e_Socket
		nodeID string
	}{
		socket: outputSock,
		nodeID: nodeID,
	}
}
