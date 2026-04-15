package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dop251/goja"
)

func gojaService(inputScript string, inpSock []e_Socket, outSock []e_Socket) ([]e_Socket, error) {
	vm := goja.New()
	timeout := time.AfterFunc(3*time.Second, func() {
		vm.Interrupt(fmt.Errorf("script execution timed out"))
	})
	defer timeout.Stop()

	for i, sock := range inpSock {
		vm.Set("var"+strconv.Itoa(i+1), sock.Data)
	}

	script := fmt.Sprintf("var execute = () => { %s };", inputScript)
	fmt.Printf("\n[GOJA]script:- %+v \n", script)

	_, err := vm.RunString(script)
	if err != nil {
		return nil, err
	}

	//inputMap := make(map[string]any)

	//for _, insock := range inpSock {
	//	inputMap[insock.ID] = insock.Data
	//}

	fn, ok := goja.AssertFunction(vm.Get("execute"))
	if !ok {
		return nil, fmt.Errorf("unable to assert execute function ")
	}

	result, err := fn(goja.Undefined())
	if err != nil {
		return nil, err
	}

	//fmt.Printf("\n[GOJA]inputMap:- %+v \n", inputMap)

	var outputs map[string](any)

	err = vm.ExportTo(result, &outputs)
	if err != nil {
		return nil, fmt.Errorf("JS must return an object: %+v", err)
	}

	for i, outsock := range outSock {
		val, found := outputs[outsock.ID]
		if found {
			outSock[i].Data = val
		}
		fmt.Printf("\n[GOJA]socketID:- %+v \n", outsock)
	}

	fmt.Printf("\n[GOJA]outsock:- %+v \n", outSock)
	fmt.Printf("\n[GOJA]outputs:- %+v \n", outputs)
	return outSock, nil
}
