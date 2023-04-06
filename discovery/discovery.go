package discovery

type Handler interface {
	OnAgentJoin(*Agent) error
	OnAgentLeave(*Agent) error
	OnAgentUpdate(*Agent) error
}

type Discovery interface {
	SetHandler(Handler)
	Agents() []*Agent
	LocalAgent() *Agent
	Start() error
	Stop()
}
