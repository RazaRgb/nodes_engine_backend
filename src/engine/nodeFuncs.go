package engine

import (
	"encoding/json"
	"fmt"
)

// func ResolveInputNumber([]engine.E_Socket)(engine.E_Socket){
//
// }

func resolveMathAdd(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	fmt.Printf("input sockets in resolve mathadd: %+v \n", inpSock)
	fmt.Printf("output sockets in resolve mathadd: %+v \n", outSock)
	fmt.Println("------------------")

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathAdd requires exactly 2 inputs and 1 output")
	}
	inp0, ok := inpSock[0].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}
	inp1, ok := inpSock[1].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}

	outSock[0].Data = inp1 + inp0
	return outSock, nil
}

func resolveInputNumber(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {
	if len(inpSock) != 0 || len(outSock) != 1 {
		return outSock, fmt.Errorf("inputNumber requires exactly 0 inputs and 1 output")
	}
	return outSock, nil
}

func resolveOutputLog(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {
	if len(inpSock) != 1 || len(outSock) != 0 {
		return outSock, fmt.Errorf("outputLog requires exactly 1 input and 0 output")
	}
	return []e_Socket{}, nil
}

func resolveMathMultiply(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {
	{ // logging
		fmt.Printf("input sockets in resolve mathadd: %+v \n", inpSock)
		fmt.Printf("output sockets in resolve mathadd: %+v \n", outSock)
		fmt.Println("------------------")
	}

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathMultiply requires exactly 2 inputs and 1 output")
	}
	inp0, ok := inpSock[0].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}
	inp1, ok := inpSock[1].Data.(float64)
	if !ok {
		return outSock, fmt.Errorf("Incorrect DataType as input in MathAdd")
	}

	outSock[0].Data = inp1 * inp0
	return outSock, nil
}

// -------------------------------
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

	//fmt.Printf("state.SockMap[e_SocketReference{: %+v \n", state.SockMap[e_SocketReference{
	//	NodeID:   nodePtr.ID,
	//	SocketID: "o1",
	//}])

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
