package client

import "github.com/werbenhu/srouter"

type Register struct {
	serf    srouter.Discovery
	handler srouter.Handler
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) SetHandler(h srouter.Handler) {
	r.handler = h
}

func (r *Register) Run(id string, addr string, advertise string, routers string, group string, service string) error {
	member := srouter.NewMember(id, addr, advertise, routers, group, service)
	r.serf = srouter.NewSerf(member)
	r.serf.SetHandler(r.handler)
	return r.serf.Start()
}

func (r *Register) Stop() {
	r.serf.Stop()
}
