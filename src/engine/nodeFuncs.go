package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// func ResolveInputNumber([]engine.E_Socket)(engine.E_Socket){
//
// }

func resolveMathAdd(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	{ //logging off
		//		fmt.Printf("input sockets in resolve mathadd: %+v \n", inpSock)
		//		fmt.Printf("output sockets in resolve mathadd: %+v \n", outSock)
		//		fmt.Println("------------------")
	}

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathAdd requires exactly 2 inputs and 1 output")
	}
	inp0, ok := inpSock[0].Data.(float64)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in MathAdd")
		outSock[0].Error = err
		return outSock, err
	}
	inp1, ok := inpSock[1].Data.(float64)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in MathAdd")
		outSock[0].Error = err
		return outSock, err
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

	retSock := inpSock[0]
	retSock.ID = "i1"

	return []e_Socket{retSock}, nil
}

func resolveMathMultiply(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {
	{ // logging off
		//	fmt.Printf("input sockets in resolve mathadd: %+v \n", inpSock)
		//	fmt.Printf("output sockets in resolve mathadd: %+v \n", outSock)
		//	fmt.Println("------------------")
	}

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("MathMultiply requires exactly 2 inputs and 1 output")
	}
	inp0, ok := inpSock[0].Data.(float64)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in MathMultiply")
		outSock[0].Error = err
		return outSock, err

	}
	inp1, ok := inpSock[1].Data.(float64)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in MathMultiply")
		outSock[1].Error = err
		return outSock, err

	}

	outSock[0].Data = inp1 * inp0
	return outSock, nil
}

func resolveInputString(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	{ // logging off
		//	fmt.Printf("input sockets in resolve inputString: %+v \n", inpSock)
		//	fmt.Printf("output sockets in resolve inputString: %+v \n", outSock)
		//	fmt.Println("------------------")
	}

	if len(inpSock) != 0 || len(outSock) != 1 {
		return outSock, fmt.Errorf("inputNumber requires exactly 0 inputs and 1 output")
	}
	return outSock, nil
}

func resolveStringConcat(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	{ //logging off
		//	fmt.Printf("input sockets in resolve stringConcat: %+v \n", inpSock)
		//	fmt.Printf("output sockets in resolve stringConcat: %+v \n", outSock)
		//	fmt.Println("------------------")
	}

	if len(inpSock) != 2 || len(outSock) != 1 {
		return outSock, fmt.Errorf("stringConcat requires exactly 2 inputs and 1 output")
	}
	inp0 := fmt.Sprint(inpSock[0].Data)

	inp1 := fmt.Sprint(inpSock[1].Data)

	outSock[0].Data = inp0 + inp1

	return outSock, nil
}

func resolveAiLLM(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	fmt.Printf("input sockets in resolve aiLLM: %+v \n", inpSock)
	fmt.Printf("output sockets in resolve aiLLM: %+v \n", outSock)
	fmt.Println("------------------")

	if len(inpSock) != 3 || len(outSock) != 1 {
		fmt.Println("not correct len")
		return outSock, fmt.Errorf("Ai LLM requires exactly 3 inputs and 1 output")
	}

	systemprompt, ok := inpSock[0].Data.(string)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in ai LLM")
		outSock[0].Error = err
		return outSock, err
	}
	userprompt, ok := inpSock[1].Data.(string)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in ai LLM")
		outSock[0].Error = err
		return outSock, err
	}
	timeout, ok := inpSock[2].Data.(float64)
	if !ok {
		err := fmt.Errorf("Incorrect DataType as input in ai LLM")
		outSock[0].Error = err
		return outSock, err
	}

	result, err := geminiService(systemprompt, userprompt, timeout)
	if err != nil {
		outSock[0].Error = err
		return outSock, err
	}

	outSock[0].Data = result
	return outSock, nil
}

func resolveCodeExecute(inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {

	fmt.Printf("input sockets in resolve codeExecute: %+v \n", inpSock)
	fmt.Printf("output sockets in resolve codeExecute: %+v \n", outSock)

	//str := ""
	//for i, _ := range inpSock {
	//	inp := fmt.Sprint(inpSock[i].Data)
	//	str += inp
	//}
	//outSock[0].Data = str

	script, ok := inpSock[0].Data.(string)
	if !ok {
		return outSock, fmt.Errorf("Incorrect datatype in code socket")
	}

	resultArr, err := gojaService(script, inpSock[1:], outSock)
	if err != nil {
		return nil, err
	}

	fmt.Println("------------------")
	return resultArr, nil
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

func popInputString(state *e_State, nodePtr *e_Node, jsonString string) error {
	jsonBytes := []byte(jsonString)
	valArr := make([]any, 1)

	err := json.Unmarshal(jsonBytes, &valArr)
	if err != nil {
		fmt.Printf("Error decoding json:\n %v \n", err)
		return err
	}

	val, ok := valArr[0].(string)
	if !ok {
		fmt.Printf("incorrect datatype in inputNumber %+v \n", valArr[0])
		return fmt.Errorf("Incorrect Datatype in inputNumber")
	}

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

func popCodeExecute(state *e_State, nodePtr *e_Node, jsonString string) error {
	jsonBytes := []byte(jsonString)
	objVal := struct {
		Inputs  []string `json:"inputs"`
		Outputs []string `json:"outputs"`
		Code    string   `json:"code"`
	}{}

	err := json.Unmarshal(jsonBytes, &objVal)
	if err != nil {
		fmt.Printf("Error decoding json:\n %v \n", err)
		return err
	}
	fmt.Printf("objVal values : %+v \n", objVal)

	nodePtr.InpSockArr = make([]e_SocketReference, len(objVal.Inputs)+1)
	nodePtr.OutSockArr = make([]e_SocketReference, len(objVal.Outputs))

	for i, _ := range nodePtr.InpSockArr {
		var sock e_Socket

		if i == 0 {
			sock = e_Socket{
				Data: objVal.Code,
				ID:   "i" + strconv.Itoa(i+1),
			}
		} else {
			sock = e_Socket{
				Data: objVal.Inputs[i-1],
				ID:   "i" + strconv.Itoa(i+1),
			}
		}

		sockref := e_SocketReference{
			NodeID:   nodePtr.ID,
			SocketID: "i" + strconv.Itoa(i+1),
		}

		state.SockMap[sockref] = &sock

		nodePtr.InpSockArr[i] = sockref

		state.SockMap[nodePtr.InpSockArr[i]] = &sock

		//fmt.Printf("vall:::: %+v \n", *state.SockMap[nodePtr.InpSockArr[i]])
	}

	{ //logs off
		//	fmt.Printf("%+v: ", len(objVal.Inputs))
		//	fmt.Printf("%+v, ", len(objVal.Outputs))

		//	fmt.Printf("%+v: ", len(nodePtr.InpSockArr))
		//	fmt.Printf("%+v, \n", len(nodePtr.OutSockArr))
	}

	fmt.Println()

	for i, _ := range nodePtr.OutSockArr {
		sock := e_Socket{
			Data: objVal.Outputs[i],
			ID:   "o" + strconv.Itoa(i+1),
		}
		sockref := e_SocketReference{
			NodeID:   nodePtr.ID,
			SocketID: "o" + strconv.Itoa(i+1),
		}

		state.SockMap[sockref] = &sock

		nodePtr.OutSockArr[i] = sockref

		state.SockMap[nodePtr.OutSockArr[i]] = &sock

		//fmt.Printf("vall:::: %+v \n", *state.SockMap[nodePtr.OutSockArr[i]])
	}
	return nil
}
