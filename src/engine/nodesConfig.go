package engine

import "fmt"

type nodeConfig struct {
	inputCount  int
	outputCount int
}

type nodeResolver func(inputSocks []e_Socket, outputSocks []e_Socket) ([]e_Socket, error)
type sockPopulate func(*e_State, *e_Node, string) error

type nodeReg struct {
	config          map[string]nodeConfig
	resolvers       map[string](nodeResolver)
	populateSockets map[string](sockPopulate)
}

var nodeRegistry nodeReg = nodeReg{
	config: map[string]nodeConfig{

		"mathAdd": {
			inputCount:  2,
			outputCount: 1,
		},
		"mathMultiply": {
			inputCount:  2,
			outputCount: 1,
		},

		"outputLog": {
			inputCount:  1,
			outputCount: 0,
		},

		"inputNumber": {
			inputCount:  0,
			outputCount: 1,
		},
		"inputString": {
			inputCount:  0,
			outputCount: 1,
		},

		"stringConcat": {
			inputCount:  2,
			outputCount: 1,
		},

		"aiLLM": {
			inputCount:  3,
			outputCount: 1,
		},

		"codeExecute": {
			inputCount:  -1,
			outputCount: -1,
		},
	},
	resolvers: map[string]nodeResolver{
		"mathAdd":      resolveMathAdd,
		"inputNumber":  resolveInputNumber,
		"outputLog":    resolveOutputLog,
		"mathMultiply": resolveMathMultiply,
		"inputString":  resolveInputString,
		"stringConcat": resolveStringConcat,
		"aiLLM":        resolveAiLLM,
		"codeExecute":  resolveCodeExecute,
	},
	populateSockets: map[string]sockPopulate{
		"inputNumber": popInputNumber,
		"inputString": popInputString,
		"codeExecute": popCodeExecute,
	},
}

func getNodeConfig(nodeType string) (nodeConfig, error) {
	nc, found := nodeRegistry.config[nodeType]
	if !found {
		return nodeConfig{}, fmt.Errorf("Node Config not found for %s", nodeType)
	}
	return nc, nil
}

func getNodeResolver(nodeType string) (nodeResolver, error) {
	res, found := nodeRegistry.resolvers[nodeType]
	if !found {
		return nil, fmt.Errorf("Node resolver not found for %s", nodeType)
	}
	return res, nil
}

func getNodePopulateFunc(nodeType string) (sockPopulate, bool) {
	f, found := nodeRegistry.populateSockets[nodeType]
	if !found {
		return nil, false
	}
	return f, true
}
