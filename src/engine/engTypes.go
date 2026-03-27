package engine

type E_Socket struct {
	ID    string
	Data  any
	Error error
}

type E_SocketReference struct {
	NodeID   string
	SocketID string
}

type E_Node struct {
	ID         string
	InpSockArr []E_SocketReference
	OutSockArr []E_SocketReference
	NodeType   string
}

type E_State struct {
	NodeMap map[string]E_Node
	FuncMap map[string](func(inputSocks []E_Socket, outputSocks []E_Socket) ([]E_Socket, error))

	AdjList map[E_SocketReference]([]E_SocketReference) //directional mapping of socket connections
}
