package srouter

import (
	"errors"
	"log"
	"strconv"

	"github.com/werbenhu/chash"
	"github.com/werbenhu/srouter/api"
	"github.com/werbenhu/srouter/discovery"
)

const (
	GroupName = "router-group"
)

type SRouter struct {
	opt  *Option
	serf discovery.Discovery
	api  api.Api
}

func New(opts []IOption) *SRouter {

	option := DefaultOption()
	for _, o := range opts {
		o(option)
	}
	s := &SRouter{opt: option}

	s.serf = discovery.NewSerf(discovery.NewAgent(
		s.opt.Id,
		s.opt.Addr,
		s.opt.Advertise,
		s.opt.Members,
		GroupName,
		s.opt.Service,
	))

	s.api = api.NewHttp()
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
	s.serf.Stop()
}

func (s *SRouter) OnAgentJoin(agent *discovery.Agent) error {
	log.Printf("[INFO] a new agent joined, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.insert(agent)
}

func (s *SRouter) OnAgentLeave(agent *discovery.Agent) error {
	log.Printf("[INFO] a new agent left, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.delete(agent)
}

func (s *SRouter) OnAgentUpdate(agent *discovery.Agent) error {
	log.Printf("[INFO] a new agent updated, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)
	return s.insert(agent)
}

func (s *SRouter) delete(agent *discovery.Agent) error {
	log.Printf("[INFO] srouter delete agent, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)

	if len(agent.Service.Group) == 0 {
		return errors.New("srouter delete agent's group name can't be empty")
	}

	replicas, err := strconv.Atoi(agent.Replicas)
	if err != nil {
		return errors.New("srouter agent replicas param error")
	}

	group, _ := chash.CreateGroup(agent.Service.Group, replicas)
	if err := group.Delete(agent.Service.Id); err != nil {
		return err
	}
	return nil
}

func (s *SRouter) insert(agent *discovery.Agent) error {
	log.Printf("[INFO] srouter insert agent, id:%s, addr:%s, group:%s, service:%s\n",
		agent.Id, agent.Addr, agent.Service.Group, agent.Service.Addr)

	if len(agent.Service.Group) == 0 {
		return errors.New("srouter insert agent's group name can't be empty")
	}

	replicas, err := strconv.Atoi(agent.Replicas)
	if err != nil {
		return errors.New("srouter agent replicas param error")
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

func (s *SRouter) Match(groupName string, key string) (*discovery.Service, error) {
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return nil, err
	}
	_, payload, err := group.Match(key)
	if err != nil {
		return nil, err
	}

	agent := &discovery.Agent{}
	if err := agent.Unmarshal(payload); err != nil {
		return nil, err
	}
	return &agent.Service, nil
}

func (s *SRouter) Members(groupName string) []*discovery.Service {
	services := make([]*discovery.Service, 0)
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return services
	}

	elements := group.GetElements()

	for _, element := range elements {
		agent := &discovery.Agent{}
		if err := agent.Unmarshal(element.Payload); err != nil {
			log.Printf("[ERROR] element to agent err:%s\n", err.Error())
			continue
		}
		services = append(services, &agent.Service)
	}
	return services
}
