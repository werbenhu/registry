package registry

import (
	"encoding/json"
	"sync"
)

// The service object
type Service struct {
	// The service id
	Id string `json:"id"`

	// The group name of this service
	Group string `json:"group"`

	// The service addr provided to the client
	Addr string `json:"addr"`
}

// NewService Create a new service object
// id:
//
//	The service id
//
// group:
//
//	The group name of this service
//
// addr:
//
//	The service addr provided to the client
func NewService(id string, group string, addr string) *Service {
	return &Service{
		Id:    id,
		Group: group,
		Addr:  addr,
	}
}

// Member is used for auto-discover, when a service is discoverd a Member Object be created.
type Member struct {
	sync.Mutex

	// The service id
	Id string `json:"id"`

	// The address used to register the service to registry server.
	Bind string `json:"bind"`

	// The address that the service will advertise to registry server.
	Advertise string `json:"advertise"`

	// The addresses of the registry servers, if there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370"
	Registries string `json:"-"`

	// How many replicated elements of a service need to be virtualized
	Replicas string `json:"replicas"`

	// Service info
	Service Service `json:"service"`

	// Tags for extra info
	tags map[string]string
}

// NewSimpleMember create a simple Member object, it does not contain the address of the service
// id:
//
//	The service id
//
// bind:
//
//	The address used to register the service to registry server.
//
// advertise:
//
//	The address that the service will advertise to registry server.
func NewSimpleMember(id string, bind string, advertise string) *Member {
	return &Member{
		Id:         id,
		Bind:       bind,
		Advertise:  advertise,
		Registries: "",
		Replicas:   DefaultReplicas,
		Service: Service{
			Id: id,
		},
	}
}

// NewMember create a Member object
// id: service id
// bind:
//
//	The address used to register the service to registry server. If there is a firewall, please remember that the port needs to open both tcp and udp.
//
// advertise:
//
//	The address that the service will advertise to registry server. Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
//
// registries:
//
//	The addresses of the registry servers, if there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370"
//
// group:
//
//	Group name of the current service belongs to.
//
// addr:
//
//	The address currently provided by this service to the client, for example, the current service is an http server, that is the address 172.16.3.3:80 that http listens to.
func NewMember(id string, bind string, advertise string, registries string, group string, addr string) *Member {
	return &Member{
		Id:         id,
		Bind:       bind,
		Advertise:  advertise,
		Registries: registries,
		Replicas:   DefaultReplicas,
		Service: Service{
			Id:    id,
			Group: group,
			Addr:  addr,
		},
	}
}

// IsSelf return true if the two member's id is the same
func (m *Member) IsSelf(b *Member) bool {
	return m.Id == b.Id
}

// SetTag set extra info of this service by tag
func (m *Member) SetTag(key string, val string) {
	m.Lock()
	defer m.Unlock()
	if m.tags == nil {
		m.tags = make(map[string]string)
	}
	m.tags[key] = val
}

// GetTag get extra info of this service by tag
func (m *Member) GetTag(key string) (string, bool) {
	m.Lock()
	defer m.Unlock()
	if m.tags == nil {
		return "", false
	}
	val, ok := m.tags[key]
	return val, ok
}

// SetTags set tags for this service
func (m *Member) SetTags(tags map[string]string) {
	if m.tags == nil {
		m.tags = make(map[string]string)
	}
	for k, v := range tags {
		m.SetTag(k, v)
	}

	m.Service.Group, _ = m.GetTag(TagGroup)
	m.Service.Addr, _ = m.GetTag(TagAddr)
	m.Replicas, _ = m.GetTag(TagReplicas)
}

// GetTag get all tags of this service
func (m *Member) GetTags() map[string]string {
	m.SetTag(TagAddr, m.Service.Addr)
	m.SetTag(TagGroup, m.Service.Group)
	m.SetTag(TagReplicas, m.Replicas)

	m.Lock()
	defer m.Unlock()
	clone := make(map[string]string)
	for k, v := range m.tags {
		clone[k] = v
	}
	return clone
}

// Marshal returns the JSON encoding of member.
func (m *Member) Marshal() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return json.Marshal(m)
}

// Unmarshal parses the JSON-encoded data and stores the result
// into a Member object.
func (m *Member) Unmarshal(paylaod []byte) error {
	m.Lock()
	defer m.Unlock()
	return json.Unmarshal(paylaod, m)
}
