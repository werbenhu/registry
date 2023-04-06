package srouter

import (
	"log"
	"strconv"

	"github.com/werbenhu/chash"
)

const (
	TagGroup    = "group"
	TagService  = "service"
	TagReplicas = "replicas"

	SRouterName     = "srouter-group"
	DefaultReplicas = "10000"
)

type SRouter struct {
	opt  *Option
	serf Discovery
	api  Api
}

func New(opts []IOption) *SRouter {

	option := DefaultOption()
	for _, o := range opts {
		o(option)
	}
	s := &SRouter{opt: option}

	s.serf = NewSerf(NewAgent(
		s.opt.Id,
		s.opt.Addr,
		s.opt.Advertise,
		s.opt.Routers,
		SRouterName,
		s.opt.Service,
	))

	s.api = NewRpcServer()
	s.serf.SetHandler(s)
	return s
}

func (s *SRouter) Serve() error {
	go func() {
		if err := s.serf.Start(); err != nil {
			log.Panic(err)
		}
		if err := s.api.Start(s.opt.ApiPort); err != nil {
			log.Panic(err)
		}
	}()
	return nil
}

func (s *SRouter) Close() {
	s.api.Stop()
	s.serf.Stop()
}

func (s *SRouter) OnAgentJoin(agent *Agent) error {
	log.Printf("[INFO] a new agent joined, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.insert(agent)
}

func (s *SRouter) OnAgentLeave(agent *Agent) error {
	log.Printf("[INFO] a new agent left, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.delete(agent)
}

func (s *SRouter) OnAgentUpdate(agent *Agent) error {
	log.Printf("[INFO] a new agent updated, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.insert(agent)
}

func (s *SRouter) delete(agent *Agent) error {
	if len(agent.Service.Group) == 0 {
		return ErrGroupNameEmpty
	}

	replicas, err := strconv.Atoi(agent.Replicas)
	if err != nil {
		return ErrReplicasParam
	}

	group, _ := chash.CreateGroup(agent.Service.Group, replicas)
	if err := group.Delete(agent.Service.Id); err != nil {
		return err
	}
	return nil
}

func (s *SRouter) insert(agent *Agent) error {
	if len(agent.Service.Group) == 0 {
		return ErrGroupNameEmpty
	}

	replicas, err := strconv.Atoi(agent.Replicas)
	if err != nil {
		return ErrReplicasParam
	}

	payload, err := agent.Marshal()
	if err != nil {
		return err
	}

	group, _ := chash.CreateGroup(agent.Service.Group, replicas)
	if err := group.Insert(agent.Service.Id, payload); err != nil {
		return err
	}
	return nil
}

func (s *SRouter) Match(groupName string, key string) (*Service, error) {
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return nil, err
	}
	_, payload, err := group.Match(key)
	if err != nil {
		return nil, err
	}

	agent := &Agent{}
	if err := agent.Unmarshal(payload); err != nil {
		return nil, err
	}
	return &agent.Service, nil
}

func (s *SRouter) Members(groupName string) []*Service {
	services := make([]*Service, 0)
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return services
	}

	elements := group.GetElements()

	for _, element := range elements {
		agent := &Agent{}
		if err := agent.Unmarshal(element.Payload); err != nil {
			log.Printf("[ERROR] element to agent err:%s\n", err.Error())
			continue
		}
		services = append(services, &agent.Service)
	}
	return services
}
