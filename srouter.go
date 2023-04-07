package srouter

import (
	"log"
	"strconv"

	"github.com/werbenhu/chash"
)

const (
	TagGroup        = "group"
	TagService      = "service"
	TagReplicas     = "replicas"
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

	s.serf = NewSerf(NewMember(
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

func (s *SRouter) OnMemberJoin(m *Member) error {
	log.Printf("[INFO] a new member joined, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

func (s *SRouter) OnMemberLeave(m *Member) error {
	log.Printf("[INFO] a new member left, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.delete(m)
}

func (s *SRouter) OnMemberUpdate(m *Member) error {
	log.Printf("[INFO] a new member updated, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

func (s *SRouter) delete(m *Member) error {
	if len(m.Service.Group) == 0 {
		return ErrGroupNameEmpty
	}

	replicas, err := strconv.Atoi(m.Replicas)
	if err != nil {
		return ErrReplicasParam
	}

	group, _ := chash.CreateGroup(m.Service.Group, replicas)
	if err := group.Delete(m.Service.Id); err != nil {
		return err
	}
	return nil
}

func (s *SRouter) insert(m *Member) error {
	if len(m.Service.Group) == 0 {
		return ErrGroupNameEmpty
	}

	replicas, err := strconv.Atoi(m.Replicas)
	if err != nil {
		return ErrReplicasParam
	}

	payload, err := m.Marshal()
	if err != nil {
		return err
	}

	group, _ := chash.CreateGroup(m.Service.Group, replicas)
	if err := group.Insert(m.Service.Id, payload); err != nil {
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

	m := &Member{}
	if err := m.Unmarshal(payload); err != nil {
		return nil, err
	}
	return &m.Service, nil
}

func (s *SRouter) Members(groupName string) []*Service {
	services := make([]*Service, 0)
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return services
	}

	elements := group.GetElements()

	for _, element := range elements {
		m := &Member{}
		if err := m.Unmarshal(element.Payload); err != nil {
			log.Printf("[ERROR] element to member err:%s\n", err.Error())
			continue
		}
		services = append(services, &m.Service)
	}
	return services
}
