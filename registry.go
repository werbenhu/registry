package registry

import (
	"log"
	"strconv"

	"github.com/werbenhu/chash"
)

const (
	//路由服务默认的组名
	registryName = "registry-group"

	//一致性哈希，每一个服务需要虚拟成多少份elemnts
	DefaultReplicas = "10000"
)

type registry struct {
	opt  *Option
	serf Discovery
	api  Api
}

func New(opts []IOption) *registry {

	option := DefaultOption()
	for _, o := range opts {
		o(option)
	}
	s := &registry{opt: option}

	if len(s.opt.Service) == 0 {
		s.opt.Service = s.opt.ApiAddr
	}

	if len(s.opt.Advertise) == 0 {
		s.opt.Advertise = s.opt.Addr
	}

	s.serf = NewSerf(NewMember(
		s.opt.Id,
		s.opt.Addr,
		s.opt.Advertise,
		s.opt.Registries,
		registryName,
		s.opt.Service,
	))

	s.api = NewRpcServer()
	s.serf.SetHandler(s)
	return s
}

func (s *registry) Serve() error {
	if err := s.serf.Start(); err != nil {
		return err
	}
	go func() {
		if err := s.api.Start(s.opt.ApiAddr); err != nil {
			log.Panic(err)
		}
	}()
	return nil
}

func (s *registry) Close() {
	s.api.Stop()
	s.serf.Stop()
}

func (s *registry) OnMemberJoin(m *Member) error {
	log.Printf("[INFO] a new member joined, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

func (s *registry) OnMemberLeave(m *Member) error {
	log.Printf("[INFO] a new member left, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.delete(m)
}

func (s *registry) OnMemberUpdate(m *Member) error {
	log.Printf("[INFO] a new member updated, id:%s, addr:%s, group:%s, service:%s\n",
		m.Id, m.Addr, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

func (s *registry) delete(m *Member) error {
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

func (s *registry) insert(m *Member) error {
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

func (s *registry) Match(groupName string, key string) (*Service, error) {
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

func (s *registry) Members(groupName string) []*Service {
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
