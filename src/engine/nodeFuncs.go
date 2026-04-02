package engine

import (
	"encoding/json"
	"fmt"
)

// func ResolveInputNumber([]engine.E_Socket)(engine.E_Socket){
//
// }

func resolveMathAdd(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathAdd requires exactly 2 inputs and 1 output")
	}
	_, ok := inpSock[0].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}
	_, ok = inpSock[1].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}

	outSock[0].Data = inpSock[0].Data.(float64) + inpSock[1].Data.(float64)
	return outSock, nil
}

func resolveInputNumber(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	if len(inpSock) != 0 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathAdd requires exactly 2 inputs and 1 output")
	}

	//val,ok :=

	outSock[0].Data = inpSock[0].Data.(float64) + inpSock[1].Data.(float64)
	return outSock, nil
}

// --
func popInputNumber(state *e_State, nodePtr *e_Node, jsonString string) error {
	jsonBytes := []byte(jsonString)
	valArr := make([]any, 1)

	err := json.Unmarshal(jsonBytes, &valArr)
	if err != nil {
		fmt.Printf("Error decoding json:\n %v \n", err)
		return err
	}

	val, ok := valArr[0].(float64)
	if !ok {
		fmt.Printf("incorrect datatype in inputNumber %+v \n", valArr[0])
		return fmt.Errorf("Incorrect Datatype in inputNumber")
	}
	fmt.Printf("state.SockMap[e_SocketReference{: %+v \n", state.SockMap[e_SocketReference{
		NodeID:   nodePtr.ID,
		SocketID: "o1",
	}])

	sockPtr, exists := state.SockMap[e_SocketReference{
		NodeID:   nodePtr.ID,
		SocketID: "o1",
	}]
	if !exists {
		return fmt.Errorf("socket reference lookup failed %+v\n",
			e_SocketReference{
				NodeID:   nodePtr.ID,
				SocketID: "o1",
			})
	}

	sockPtr.Data = val
	sockPtr.ID = "o1"
	return nil
}


