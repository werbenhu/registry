package srouter

type Handler interface {
	OnMemberJoin(*Member) error
	OnMemberLeave(*Member) error
	OnMemberUpdate(*Member) error
}

type Discovery interface {
	SetHandler(Handler)
	Members() []*Member
	LocalMember() *Member
	Start() error
	Stop()
}
