package engine

import "fmt"

type nodeConfig struct {
	InputCount  int
	OutputCount int
}

var nodeRegistry map[string]nodeConfig = map[string]nodeConfig{
	"MathAdd": {
		InputCount:  2,
		OutputCount: 1,
	},
	"OutputLog": {
		InputCount:  1,
		OutputCount: 0,
	},
	"InputNumber": {
		InputCount:  0,
		OutputCount: 1,
	},
}

func GetNodeConfig(nodeType string) (nodeConfig, error) {
	nc, found := nodeRegistry[nodeType]
	if !found {
		return nodeConfig{}, fmt.Errorf("Node Config not found for %s", nodeType)
	}
	return nc, nil
}
