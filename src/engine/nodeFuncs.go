package engine

import "fmt"

// func ResolveInputNumber([]engine.E_Socket)(engine.E_Socket){
//
// }

func ResolveMathAdd(inpSock []E_Socket, outSock []E_Socket) ([]E_Socket, error) {
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
