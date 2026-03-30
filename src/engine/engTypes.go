package engine

type e_Socket struct {
	ID    string
	Data  any
	Error error
}

type e_SocketReference struct {
	NodeID   string
	SocketID string
}

type e_Node struct {
	ID         string
	InpSockArr []e_SocketReference
	OutSockArr []e_SocketReference
	NodeType   string
}

type e_State struct {
	NodeMap map[string]e_Node
	FuncMap map[string](func(inputSocks []e_Socket, outputSocks []e_Socket) ([]e_Socket, error))

	AdjList map[e_SocketReference]([]e_SocketReference) //directional mapping of socket connections
}
