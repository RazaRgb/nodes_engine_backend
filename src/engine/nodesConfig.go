package engine

import "fmt"

type nodeConfig struct {
	inputCount  int
	outputCount int
}

type resolver func(inputSocks []e_Socket, outputSocks []e_Socket) ([]e_Socket, error)

type nodeReg struct {
	config    map[string]nodeConfig
	resolvers map[string](resolver)
}

var nodeRegistry nodeReg = nodeReg{
	config: map[string]nodeConfig{
		"mathAdd": {
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
	},
	resolvers: map[string]resolver{
		"mathAdd":     resolveMathAdd,
		"inputNumber": resolveInputNumber,
		"outputLog":   resolveMathAdd,
	},
}

func getNodeConfig(nodeType string) (nodeConfig, error) {
	nc, found := nodeRegistry.config[nodeType]
	if !found {
		return nodeConfig{}, fmt.Errorf("Node Config not found for %s", nodeType)
	}
	return nc, nil
}

func getNodeResolver(nodeType string) (resolver, error) {
	res, found := nodeRegistry.resolvers[nodeType]
	if !found {
		return nil, fmt.Errorf("Node resolver not found for %s", nodeType)
	}
	return res, nil
}
