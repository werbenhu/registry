package registry

import (
	"log"
	"strconv"

	"github.com/werbenhu/chash"
)

const (
	// Registry service default group name
	registryName = "registry-group"

	// This indicates how many replicated elements of a service need to be virtualized
	DefaultReplicas = "10000"
)

// Registry is the registry server object
type Registry struct {
	opt  *Option
	serf Discovery
	api  Api
}

// New a registry object that can start a registry server when calling Serve().
func New(opts []IOption) *Registry {

	option := DefaultOption()
	for _, o := range opts {
		o(option)
	}
	s := &Registry{opt: option}

	if len(s.opt.Advertise) == 0 {
		s.opt.Advertise = s.opt.Addr
	}

	if len(s.opt.Advertise) == 0 {
		s.opt.Advertise = s.opt.Addr
	}

	s.serf = NewSerf(NewMember(
		s.opt.Id,
		s.opt.Bind,
		s.opt.BindAdvertise,
		s.opt.Registries,
		registryName,
		s.opt.Advertise,
	))

	s.api = NewRpcServer()
	s.serf.SetHandler(s)
	return s
}

// Serve run the registry server
func (s *Registry) Serve() error {
	if err := s.serf.Start(); err != nil {
		return err
	}
	go func() {
		if err := s.api.Start(s.opt.Addr); err != nil {
			log.Panic(err)
		}
	}()
	return nil
}

// Close will close the registry server
func (s *Registry) Close() {
	s.api.Stop()
	s.serf.Stop()
}

// OnMemberJoin triggered when a new service is registered
func (s *Registry) OnMemberJoin(m *Member) error {
	log.Printf("[INFO] a new member joined, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

// OnMemberLeave triggered when a new service is leaves
func (s *Registry) OnMemberLeave(m *Member) error {
	log.Printf("[INFO] a new member left, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.delete(m)
}

// OnMemberLeave triggered when a new service is updated
func (s *Registry) OnMemberUpdate(m *Member) error {
	log.Printf("[INFO] a new member updated, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

// delete a service from chash
func (s *Registry) delete(m *Member) error {
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

// insert a service to chash
func (s *Registry) insert(m *Member) error {
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
	if err := group.Upsert(m.Service.Id, payload); err != nil {
		return err
	}
	return nil
}

// Match assign a service to a key with consistent hashing algorithm
func (s *Registry) Match(groupName string, key string) (*Service, error) {
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

// Members get services list of a group
// groupName:
//
//	the group name of the services
func (s *Registry) Members(groupName string) []*Service {
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
