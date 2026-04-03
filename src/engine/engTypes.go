package engine

import ()

type e_Socket struct {
	ID    string `json:"id"`
	Data  any    `json:"data"`
	Error error  `json:"error"`
}

type e_SocketReference struct {
	NodeID   string `json:"nodeID"`
	SocketID string `json:"SocketID"`
}

type e_Node struct {
	ID         string              `json:"id"`
	InpSockArr []e_SocketReference `json:"inpSockArr"`
	OutSockArr []e_SocketReference `json:"outSockArr"`
	NodeType   string              `json:"nodeType"`
}

type e_State struct {
	NodeMap     map[string]*e_Node `json:"nodeMap"`
	SockMap     map[e_SocketReference]*e_Socket
	DegMap      map[string]int                              `json:"degMap"`
	AdjList     map[e_SocketReference]([]e_SocketReference) `json:"AdjList"` //directional mapping of socket connections
	nodeCounter int
}

type e_Communication struct {
	interrupt      chan error
	valuePropagate chan e_workerValue
}

type e_workerValue struct {
	socket []e_Socket
	nodeID string
}

func (e e_SocketReference) MarshalText() ([]byte, error) {
	res := e.NodeID + ":" + e.SocketID
	return []byte(res), nil
}
