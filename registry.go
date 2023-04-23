package registry

import (
	"log"
	"strconv"

	"github.com/werbenhu/chash"
)

const (
	registryName    = "registry-group" // Default group name for registry service
	DefaultReplicas = "10000"          // Default number of replicas to virtualize a service
)

// Registry is the registry server object
type Registry struct {
	opt  *Option
	serf Discovery
	api  Api
}

// New creates a new registry object that can start a registry server when calling Serve().
func New(opts []IOption) *Registry {

	option := DefaultOption()
	for _, o := range opts {
		o(option)
	}

	s := &Registry{opt: option}

	// If Advertise is not set, set it to Addr
	if len(s.opt.Advertise) == 0 {
		s.opt.Advertise = s.opt.Addr
	}

	s.api = NewRpcServer()
	s.serf = NewSerf(NewMember(
		s.opt.Id,
		s.opt.Bind,
		s.opt.BindAdvertise,
		s.opt.Registries,
		registryName,
		s.opt.Advertise,
	))
	s.serf.SetHandler(s)
	return s
}

// Serve runs the registry server
func (s *Registry) Serve() {
	if err := s.serf.Start(); err != nil {
		panic(err)
	}
	if err := s.api.Start(s.opt.Addr); err != nil {
		panic(err)
	}
}

// Close closes the registry server
func (s *Registry) Close() {
	if s.serf != nil {
		s.serf.Stop()
	}
	if s.api != nil {
		s.api.Stop()
	}
	chash.RemoveAllGroup() // Remove all groups from chash
	log.Printf("[DEBUG] registry server is closed.\n")
}

// OnMemberJoin is triggered when a new service is registered
func (s *Registry) OnMemberJoin(m *Member) error {
	log.Printf("[INFO] a new member joined, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

// OnMemberLeave is triggered when a service leaves
func (s *Registry) OnMemberLeave(m *Member) error {
	log.Printf("[INFO] a new member left, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.delete(m)
}

// OnMemberUpdate is triggered when a service is updated
func (s *Registry) OnMemberUpdate(m *Member) error {
	log.Printf("[INFO] a new member updated, id:%s, bind:%s, group:%s, service:%s\n",
		m.Id, m.Bind, m.Service.Group, m.Service.Addr)
	return s.insert(m)
}

// delete removes a service from chash
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

// insert adds a service to chash
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

// Match uses a consistent hashing algorithm to assign a service to a key.
func (s *Registry) Match(groupName string, key string) (*Service, error) {
	// Get the group associated with the group name.
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return nil, err
	}

	// Find the element in the group that matches the key.
	_, payload, err := group.Match(key)
	if err != nil {
		return nil, err
	}

	// Unmarshal the payload to create a Member object.
	m := &Member{}
	if err := m.Unmarshal(payload); err != nil {
		return nil, err
	}

	// Return the Service associated with the Member.
	return &m.Service, nil
}

// Members returns a list of services for a given group name.
func (s *Registry) Members(groupName string) []*Service {
	// Create an empty list of services.
	services := make([]*Service, 0)

	// Get the group associated with the group name.
	group, err := chash.GetGroup(groupName)
	if err != nil {
		return services
	}

	// Get the elements in the group and create a Service object for each.
	elements := group.GetElements()
	for _, element := range elements {
		// Unmarshal the payload to create a Member object.
		m := &Member{}
		if err := m.Unmarshal(element.Payload); err != nil {
			log.Printf("[ERROR] element to member err:%s\n", err.Error())
			continue
		}

		// Add the Service associated with the Member to the list.
		services = append(services, &m.Service)
	}

	// Return the list of services.
	return services
}
